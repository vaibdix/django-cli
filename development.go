package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func (m *Model) startDevelopmentEnvironment() {
	projectPath := m.projectName
	if !filepath.IsAbs(projectPath) {
		wd, _ := os.Getwd()
		projectPath = filepath.Join(wd, m.projectName)
	}

	// Create VS Code tasks.json to automate terminal setup
	createVSCodeTasks(projectPath, m.setupTailwind)

	// Open VS Code with the project
	cmd := exec.Command("code", projectPath)
	cmd.Start()

	// Add instructions to the step messages
	m.stepMessages = append(m.stepMessages, "✨ VS Code will open with the project.")
	m.stepMessages = append(m.stepMessages, "✨ Two terminals will automatically open with your development servers.")
	m.stepMessages = append(m.stepMessages, "✨ You can run the tasks manually from the Terminal menu > Run Task.")
}

func createVSCodeTasks(projectPath string, setupTailwind bool) {
	// Create .vscode directory if it doesn't exist
	vscodeDir := filepath.Join(projectPath, ".vscode")
	os.MkdirAll(vscodeDir, 0755)

	// Determine the Python command for different platforms
	pythonCmd := "python"
	activateCmd := ""

	// Check for virtual environment
	venvDirs := []string{"venv", "env", ".venv"}
	for _, dir := range venvDirs {
		if _, err := os.Stat(filepath.Join(projectPath, dir)); err == nil {
			if runtime.GOOS == "windows" {
				activateCmd = ".\\" + dir + "\\Scripts\\activate && "
				pythonCmd = "python"
			} else {
				activateCmd = "source ./" + dir + "/bin/activate && "
				pythonCmd = "python"
			}
			break
		}
	}

	// Create tasks configuration
	tasks := map[string]interface{}{
		"version": "2.0.0",
		"tasks":   []map[string]interface{}{},
	}

	// Task for running Django server
	djangoTask := map[string]interface{}{
		"label":       "Django: Run server",
		"type":        "shell",
		"command":     activateCmd + pythonCmd + " manage.py runserver",
		"presentation": map[string]interface{}{
			"reveal":          "always",
			"panel":           "new",
			"group":           "development",
			"showReuseMessage": false,
		},
		"runOptions": map[string]interface{}{
			"runOn": "folderOpen",
		},
	}

	// Add Django task
	tasks["tasks"] = append(tasks["tasks"].([]map[string]interface{}), djangoTask)

	// If Tailwind is set up, add task for watching CSS
	if setupTailwind {
		tailwindTask := map[string]interface{}{
			"label":       "Tailwind: Watch CSS",
			"type":        "shell",
			"command":     "npm run watch:css",
			"presentation": map[string]interface{}{
				"reveal":          "always",
				"panel":           "new",
				"group":           "development",
				"showReuseMessage": false,
			},
			"runOptions": map[string]interface{}{
				"runOn": "folderOpen",
			},
		}

		// Add Tailwind task
		tasks["tasks"] = append(tasks["tasks"].([]map[string]interface{}), tailwindTask)
	}

	// Write tasks.json
	tasksJSON, _ := json.MarshalIndent(tasks, "", "  ")
	tasksFile := filepath.Join(vscodeDir, "tasks.json")
	os.WriteFile(tasksFile, tasksJSON, 0644)

	// Create a welcome file with instructions
	welcomeContent := "# Welcome to Your Django Project!\n\n"
	welcomeContent += "Two terminal windows should automatically open with:\n\n"

	if setupTailwind {
		welcomeContent += "1. Django development server (`python manage.py runserver`)\n"
		welcomeContent += "2. Tailwind CSS watcher (`npm run watch:css`)\n\n"
	} else {
		welcomeContent += "1. Django development server (`python manage.py runserver`)\n\n"
	}

	welcomeContent += "If the terminals didn't open automatically, you can run these tasks manually from the menu:\n"
	welcomeContent += "Terminal > Run Task\n\n"

	welcomeFile := filepath.Join(projectPath, "WELCOME.md")
	os.WriteFile(welcomeFile, []byte(welcomeContent), 0644)
}

func (m *Model) setupServerInstructions(projectPath string) {
	if m.runServer {
		pythonVenvPath := getPythonPath(projectPath)
		m.stepMessages = append(m.stepMessages, "✨ To start the server: cd "+m.projectName+" && "+pythonVenvPath+" manage.py runserver")
		if m.setupTailwind {
			m.stepMessages = append(m.stepMessages, "✨ To watch Tailwind CSS: cd "+m.projectName+" && npm run watch:css")
		}
	}
}