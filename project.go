package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func (m *Model) CreateProject() {
	var currentErr error
	errChan := make(chan error, 5) // Buffer for multiple potential errors
	var wg sync.WaitGroup

	defer func() {
		if m.program != nil {
			m.program.Send(projectCreationDoneMsg{err: currentErr})
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
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := m.createVirtualEnvironment(projectPath); err != nil {
			errChan <- err
			return
		}
		if err := m.installDjango(projectPath); err != nil {
			errChan <- err
			return
		}
	}()

	// Wait for virtual environment and Django installation before proceeding with project creation
	wg.Wait()

	// Check for any errors from the goroutine
	select {
	case err := <-errChan:
		currentErr = err
		return
	default:
	}

	if currentErr = m.createDjangoProject(projectPath); currentErr != nil {
		return
	}

	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")

	// Start concurrent tasks that don't depend on each other
	var templateWg sync.WaitGroup

	// Configure Django settings
	templateWg.Add(1)
	go func() {
		defer templateWg.Done()
		if err := m.configureDjangoSettings(settingsPath); err != nil {
			errChan <- err
		}
	}()

	// Setup project URLs
	templateWg.Add(1)
	go func() {
		defer templateWg.Done()
		if err := m.setupProjectUrls(projectPath); err != nil {
			errChan <- err
		}
	}()

	// Setup templates if enabled
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

	// Initialize Git repository in parallel if enabled
	if m.initializeGit {
		templateWg.Add(1)
		go func() {
			defer templateWg.Done()
			if err := m.initializeGitRepository(projectPath); err != nil {
				errChan <- err
			}
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

	// Wait for all concurrent tasks to complete
	templateWg.Wait()

	// Check for any errors from goroutines
	select {
	case err := <-errChan:
		currentErr = err
		return
	default:
	}

	m.updateProgress("Finalizing settings configuration...")

	// Create Django app if specified (must be done after settings are configured)
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

	// Run migrations (must be done last)
	if currentErr = m.runDjangoMigrations(projectPath); currentErr != nil {
		return
	}

	if m.runServer {
		m.setupServerInstructions(projectPath)
	}

	m.stepMessages = append(m.stepMessages, "✅ Django project setup complete!")
}
