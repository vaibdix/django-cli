package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// startDevelopmentEnvironment opens VS Code and starts development servers in integrated terminal
func (m *Model) startDevelopmentEnvironment() {
	projectPath := m.projectName
	if !filepath.IsAbs(projectPath) {
		wd, _ := os.Getwd()
		projectPath = filepath.Join(wd, m.projectName)
	}

	// Open VS Code first
	exec.Command("code", projectPath).Start()

	// Get Python path for Django server
	pythonVenvPath := getPythonPath(projectPath)

	if m.setupTailwind {
		// Use AppleScript to open VS Code integrated terminal and split it
		appleScript := fmt.Sprintf(`tell application "Visual Studio Code"
			activate
			delay 2
		end tell
		tell application "System Events"
			tell process "Visual Studio Code"
				key code 50 using {control down}
				delay 1
				keystroke "npm run watch:css"
				key code 36
				delay 1
				key code 42 using {command down}
				delay 1
				keystroke "%s manage.py runserver"
				key code 36
			end tell
		end tell`, pythonVenvPath)
		exec.Command("osascript", "-e", appleScript).Start()
	} else {
		// Use AppleScript to open VS Code integrated terminal for Django server only
		appleScript := fmt.Sprintf(`tell application "Visual Studio Code"
			activate
			delay 2
		end tell
		tell application "System Events"
			tell process "Visual Studio Code"
				key code 50 using {control down}
				delay 1
				keystroke "%s manage.py runserver"
				key code 36
			end tell
		end tell`, pythonVenvPath)
		exec.Command("osascript", "-e", appleScript).Start()
	}
}

// setupServerInstructions provides instructions for starting the development server
func (m *Model) setupServerInstructions(projectPath string) {
	if m.runServer {
		pythonVenvPath := getPythonPath(projectPath)
		m.stepMessages = append(m.stepMessages, "✨ To start the server: cd "+m.projectName+" && "+pythonVenvPath+" manage.py runserver")
		if m.setupTailwind {
			m.stepMessages = append(m.stepMessages, "✨ To watch Tailwind CSS: cd "+m.projectName+" && npm run watch:css")
		}
	}
}