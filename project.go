package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func (m *Model) CreateProject() {
	var currentErr error
	errChan := make(chan error, 10) // Increased buffer for more concurrent operations

	defer func() {
		if m.program != nil {
			m.program.Send(projectCreationDoneMsg{err: currentErr})
		}
		if currentErr == nil {
			m.showPerformanceTip()
		}
	}()

	if m.projectName == "" {
		currentErr = fmt.Errorf("project name cannot be empty")
		return
	}

	m.totalSteps = m.calculateTotalSteps()
	m.completedSteps = 0

	projectPath := m.projectName
	if !filepath.IsAbs(projectPath) {
		wd, err := os.Getwd()
		if err != nil {
			currentErr = fmt.Errorf("failed to get working directory: %v", err)
			return
		}
		projectPath = filepath.Join(wd, m.projectName)
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		currentErr = fmt.Errorf("failed to create project directory: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Project directory created: %s", projectPath))
	m.updateProgress("Creating project directory...")

	// Create virtual environment and install Django concurrently
	if currentErr = m.createVirtualEnvironment(projectPath); currentErr != nil {
		return
	}
	if currentErr = m.installDjango(projectPath); currentErr != nil {
		return
	}
	if currentErr = m.createDjangoProject(projectPath); currentErr != nil {
		return
	}

	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")

	// Phase 2: Parallel operations that can run concurrently
	var parallelWg sync.WaitGroup

	// Configure Django settings (required for other operations)
	parallelWg.Add(1)
	go func() {
		defer parallelWg.Done()
		if err := m.configureDjangoSettings(settingsPath); err != nil {
			errChan <- err
		}
	}()

	// Setup project URLs (independent operation)
	parallelWg.Add(1)
	go func() {
		defer parallelWg.Done()
		if err := m.setupProjectUrls(projectPath); err != nil {
			errChan <- err
		}
	}()

	// Initialize Git repository in parallel if enabled (independent operation)
	if m.initializeGit {
		parallelWg.Add(1)
		go func() {
			defer parallelWg.Done()
			if err := m.initializeGitRepository(projectPath); err != nil {
				errChan <- err
			}
		}()
	}
	parallelWg.Wait()
	select {
	case err := <-errChan:
		currentErr = err
		return
	default:
	}
	var templateWg sync.WaitGroup


	if m.createTemplates {
		templateWg.Add(1)
		go func() {
			defer templateWg.Done()
			if err := m.setupGlobalTemplates(projectPath); err != nil {
				errChan <- err
				return
			}
			settingsContentBytes, err := os.ReadFile(settingsPath)
			if err != nil {
				errChan <- fmt.Errorf("failed to read settings.py for templates: %v", err)
				return
			}
			settingsContent := updateSettingsForTemplates(string(settingsContentBytes))
			if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
				errChan <- fmt.Errorf("failed to write updated settings.py: %v", err)
				return
			}
			m.stepMessages = append(m.stepMessages, "✅ Configured settings for global templates and static files.")
		}()
	}

	// Setup Tailwind CSS in parallel if enabled
	if m.setupTailwind {
		templateWg.Add(1)
		go func() {
			defer templateWg.Done()
			m.updateProgress("Setting up Tailwind CSS v4...")
			if err := m.setupTailwindCSS(projectPath); err != nil {
				errChan <- err
			}
		}()
	}

	templateWg.Wait()
	select {
	case err := <-errChan:
		currentErr = err
		return
	default:
	}

	m.updateProgress("Finalizing settings configuration...")

	if m.appName != "" {
		if currentErr = m.createDjangoApp(projectPath, settingsPath); currentErr != nil {
			return
		}
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created Django app '%s' with templates and URLs.", m.appName))
	}

	// Setup REST framework if enabled (must be done after app creation)
	if m.setupRestFramework {
		if currentErr = m.setupDjangoRestFramework(projectPath); currentErr != nil {
			return
		}
	}

	if currentErr = m.runDjangoMigrations(projectPath); currentErr != nil {
		return
	}

	if m.runServer {
		m.setupServerInstructions(projectPath)
	}

	m.stepMessages = append(m.stepMessages, "✅ Django project setup complete!")
}