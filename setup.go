package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (m *Model) createVirtualEnvironment(projectPath string) error {
	cmd := exec.Command("python3", "-m", "venv", ".venv")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Virtual environment created.")
	m.updateProgress("Creating virtual environment...")
	return nil
}

func (m *Model) installDjango(projectPath string) error {
	pipPath := getPipPath(projectPath)

	cmd := exec.Command(pipPath, "install", "django")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Django: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Django installed.")
	m.updateProgress("Installing Django...")

	cmd = exec.Command(pipPath, "install", "django-browser-reload")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install django-browser-reload: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ django-browser-reload installed.")
	m.updateProgress("Installing development dependencies...")

	return nil
}

func (m *Model) createDjangoProject(projectPath string) error {
	pythonVenvPath := getPythonPath(projectPath)
	cmd := exec.Command(pythonVenvPath, "-m", "django", "startproject", m.projectName, ".")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create Django project: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Django project '%s' created.", m.projectName))
	m.updateProgress("Creating Django project...")
	return nil
}

func (m *Model) configureDjangoSettings(settingsPath string) error {
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings.py: %v", err)
	}
	settingsStr := string(settingsContent)
	if !strings.Contains(settingsStr, "django_browser_reload") {
		updatedSettings, err := addToListInSettingsPy(settingsStr, "INSTALLED_APPS", "django_browser_reload")
		if err != nil {
			return fmt.Errorf("failed to add django_browser_reload to INSTALLED_APPS: %v", err)
		}
		settingsStr = updatedSettings
	}
	if !strings.Contains(settingsStr, "django_browser_reload.middleware.BrowserReloadMiddleware") {
		updatedSettings, err := addToListInSettingsPy(settingsStr, "MIDDLEWARE", "django_browser_reload.middleware.BrowserReloadMiddleware")
		if err != nil {
			return fmt.Errorf("failed to add BrowserReloadMiddleware to MIDDLEWARE: %v", err)
		}
		settingsStr = updatedSettings
	}
	if err := os.WriteFile(settingsPath, []byte(settingsStr), 0644); err != nil {
		return fmt.Errorf("failed to write updated settings.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Django settings configured.")
	m.updateProgress("Configuring Django settings...")
	return nil
}
