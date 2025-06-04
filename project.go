package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Model) CreateProject() {
	var currentErr error
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
	if currentErr = m.configureDjangoSettings(settingsPath); currentErr != nil {
		return
	}

	if currentErr = m.setupProjectUrls(projectPath); currentErr != nil {
		return
	}

	if m.createTemplates {
		if currentErr = m.setupGlobalTemplates(projectPath); currentErr != nil {
			return
		}
		settingsContentBytes, err := os.ReadFile(settingsPath)
		if err != nil {
			currentErr = fmt.Errorf("failed to read settings.py for templates: %v", err)
			return
		}
		settingsContent := updateSettingsForTemplates(string(settingsContentBytes))
		if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to write updated settings.py: %v", err)
			return
		}
		m.stepMessages = append(m.stepMessages, "✅ Configured settings for global templates and static files.")
	}

	m.updateProgress("Finalizing settings configuration...")

	if m.appName != "" {
		if currentErr = m.createDjangoApp(projectPath, settingsPath); currentErr != nil {
			return
		}
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created Django app '%s' with templates and URLs.", m.appName))
	}

	if m.initializeGit {
		if currentErr = m.initializeGitRepository(projectPath); currentErr != nil {
			return
		}
	}

	if m.setupTailwind {
		m.updateProgress("Setting up Tailwind CSS v4...")
		if currentErr = m.setupTailwindCSS(projectPath); currentErr != nil {
			return
		}
	}

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
