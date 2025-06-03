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

	projectPath := m.projectName // Relative path for now
	if !filepath.IsAbs(projectPath) {
		wd, err := os.Getwd()
		if err != nil {
			currentErr = fmt.Errorf("failed to get working directory: %v", err)
			return
		}
		projectPath = filepath.Join(wd, m.projectName)
	}

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		currentErr = fmt.Errorf("failed to create project directory: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Project directory created: %s", projectPath))
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.05, status: "Creating project directory..."})
	}

	// Create virtual environment
	if currentErr = m.createVirtualEnvironment(projectPath); currentErr != nil {
		return
	}

	// Install Django and dependencies
	if currentErr = m.installDjango(projectPath); currentErr != nil {
		return
	}

	// Create Django project
	if currentErr = m.createDjangoProject(projectPath); currentErr != nil {
		return
	}

	// Configure Django settings
	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
	if currentErr = m.configureDjangoSettings(settingsPath); currentErr != nil {
		return
	}

	// Setup templates and static files if chosen
	if m.createTemplates {
		if currentErr = m.setupGlobalTemplates(projectPath); currentErr != nil {
			return
		}
		// Update settings for templates
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

	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.95, status: "Finalizing settings configuration..."})
	}

	// Create Django app if chosen
	if m.appName != "" {
		if currentErr = m.createDjangoApp(projectPath, settingsPath); currentErr != nil {
			return
		}
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created Django app '%s' with templates and URLs.", m.appName))
	}



	// Handle Git initialization if selected
	if m.initializeGit {
		if currentErr = m.initializeGitRepository(projectPath); currentErr != nil {
			return
		}
	}

	// Setup Tailwind CSS v4 if selected
	if m.setupTailwind {
		if m.program != nil {
			m.program.Send(projectProgressMsg{percent: 0.90, status: "Setting up Tailwind CSS v4..."})
		}
		if currentErr = m.setupTailwindCSS(projectPath); currentErr != nil {
			return
		}
	}

	// Handle server startup if selected
	if m.runServer {
		m.setupServerInstructions(projectPath)
	}

	m.stepMessages = append(m.stepMessages, "✅ Django project setup complete!")
	// Final progress update will be handled by projectCreationDoneMsg logic in Update()
}
