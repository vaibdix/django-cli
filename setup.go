package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (m *Model) createVirtualEnvironment(projectPath string) error {
	// Try different Python commands based on the OS
	pythonCommands := []string{"python", "python3"}
	if runtime.GOOS == "windows" {
		pythonCommands = []string{"python", "py", "python3"}
	}

	var pythonCmd string
	for _, cmd := range pythonCommands {
		if isCommandAvailable(cmd) {
			pythonCmd = cmd
			break
		}
	}

	if pythonCmd == "" {
		if runtime.GOOS == "windows" {
			return fmt.Errorf("python not found. Please:\n" +
				"1. Download Python from https://www.python.org/downloads/\n" +
				"2. During installation, CHECK 'Add Python to PATH'\n" +
				"3. Restart your terminal/command prompt\n" +
				"4. Try running this command again")
		}
		return fmt.Errorf("python not found. Please install Python 3.x and ensure it's in your PATH")
	}

	var cmd *exec.Cmd
	if isUvAvailable() {
		m.updateProgress("Creating virtual environment with uv...")
		cmd = exec.Command("uv", "venv", ".venv")
	} else {
		m.updateProgress("Creating virtual environment...")
		cmd = exec.Command(pythonCmd, "-m", "venv", ".venv")
	}

	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %v\nOutput: %s", err, string(output))
	}
	m.stepMessages = append(m.stepMessages, "✅ Virtual environment created.")
	return nil
}

func (m *Model) installDjango(projectPath string) error {
	packageManager, baseArgs := getPackageManager(projectPath)

	if isUvAvailable() {
		m.updateProgress("Installing Django with uv (faster)...")
	} else {
		m.updateProgress("Installing Django...")
	}

	args := append(baseArgs, "django", "django-browser-reload")

	if isUvAvailable() {
		args = append(args, "--quiet")
	} else {
		if runtime.GOOS == "windows" {
			args = append(args, "--progress-bar", "pretty")
		}
	}

	cmd := exec.Command(packageManager, args...)
	cmd.Dir = projectPath

	if runtime.GOOS == "windows" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Django packages: %v", err)
		}
	} else {
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to install Django packages: %v\nOutput: %s", err, string(output))
		}
	}

	m.stepMessages = append(m.stepMessages, "✅ Django installed.")
	m.stepMessages = append(m.stepMessages, "✅ django-browser-reload installed.")
	m.updateProgress("Installing development dependencies...")

	return nil
}

func (m *Model) createDjangoProject(projectPath string) error {
	var pythonVenvPath string
	if isUvAvailable() {
		pythonVenvPath = "uv"
		cmd := exec.Command(pythonVenvPath, "run", "python", "-m", "django", "startproject", m.projectName, ".")
		cmd.Dir = projectPath
		m.updateProgress("Creating Django project...")

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to create Django project: %v\nOutput: %s", err, string(output))
		}
	} else {
		pythonVenvPath = getPythonPath(projectPath)
		cmd := exec.Command(pythonVenvPath, "-m", "django", "startproject", m.projectName, ".")
		cmd.Dir = projectPath
		m.updateProgress("Creating Django project...")

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to create Django project: %v\nOutput: %s", err, string(output))
		}
	}

	m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Django project '%s' created.", m.projectName))
	return nil
}

func (m *Model) configureDjangoSettings(settingsPath string) error {
	m.updateProgress("Configuring Django settings...")

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
	return nil
}
func (m *Model) showPerformanceTip() {
	if !isUvAvailable() && runtime.GOOS == "windows" {
		m.stepMessages = append(m.stepMessages,
			"💡 Tip: Install 'uv' for faster package management: pip install uv")
	}
}
