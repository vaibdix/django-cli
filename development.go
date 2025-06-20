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

	if isUvAvailable() {
		m.stepMessages = append(m.stepMessages, "⚡ Using uv for faster development server startup.")
	}
}

func createVSCodeTasks(projectPath string, setupTailwind bool) {
	// Create .vscode directory if it doesn't exist
	vscodeDir := filepath.Join(projectPath, ".vscode")
	os.MkdirAll(vscodeDir, 0755)

	var pythonCmd string
	var activateCmd string

	// Check for virtual environment
	venvDirs := []string{"venv", "env", ".venv"}
	venvFound := false
	
	for _, dir := range venvDirs {
		if _, err := os.Stat(filepath.Join(projectPath, dir)); err == nil {
			venvFound = true
			if isUvAvailable() {
				pythonCmd = "uv run python"
				activateCmd = ""
			} else {
				if runtime.GOOS == "windows" {
					activateCmd = ".\\" + dir + "\\Scripts\\activate && "
					pythonCmd = "python"
				} else {
					activateCmd = "source ./" + dir + "/bin/activate && "
					pythonCmd = "python"
				}
			}
			break
		}
	}

	if !venvFound {
		if isUvAvailable() {
			pythonCmd = "uv run python"
			activateCmd = ""
		} else {
			pythonCmd = "python"
			activateCmd = ""
		}
	}

	// Create tasks configuration
	tasks := map[string]interface{}{
		"version": "2.0.0",
		"tasks":   []map[string]interface{}{},
	}

	// Task for running Django server
	djangoTaskLabel := "Django: Run server"
	if isUvAvailable() {
		djangoTaskLabel += " (uv optimized)"
	}

	djangoTask := map[string]interface{}{
		"label":   djangoTaskLabel,
		"type":    "shell",
		"command": activateCmd + pythonCmd + " manage.py runserver",
		"presentation": map[string]interface{}{
			"reveal":          "always",
			"panel":           "new",
			"group":           "development",
			"showReuseMessage": false,
		},
		"runOptions": map[string]interface{}{
			"runOn": "folderOpen",
		},
		"group": map[string]interface{}{
			"kind":      "build",
			"isDefault": true,
		},
	}

	if runtime.GOOS == "windows" {
		djangoTask["options"] = map[string]interface{}{
			"shell": map[string]interface{}{
				"executable": "cmd.exe",
				"args":       []string{"/c"},
			},
		}
	}

	// Add Django task
	tasks["tasks"] = append(tasks["tasks"].([]map[string]interface{}), djangoTask)

	// If Tailwind is set up, add task for watching CSS
	if setupTailwind {
		tailwindTask := map[string]interface{}{
			"label":   "Tailwind: Watch CSS",
			"type":    "shell",
			"command": "npm run watch:css",
			"presentation": map[string]interface{}{
				"reveal":          "always",
				"panel":           "new",
				"group":           "development",
				"showReuseMessage": false,
			},
			"runOptions": map[string]interface{}{
				"runOn": "folderOpen",
			},
			"group": "build",
		}

		// Add Tailwind task
		tasks["tasks"] = append(tasks["tasks"].([]map[string]interface{}), tailwindTask)
	}


	if setupTailwind {
		combinedTask := map[string]interface{}{
			"label": "Start Development Environment",
			"dependsOrder": "parallel",
			"dependsOn": []string{djangoTaskLabel, "Tailwind: Watch CSS"},
			"group": map[string]interface{}{
				"kind":      "build",
				"isDefault": false,
			},
		}
		tasks["tasks"] = append(tasks["tasks"].([]map[string]interface{}), combinedTask)
	}

	// Write tasks.json
	tasksJSON, _ := json.MarshalIndent(tasks, "", "  ")
	tasksFile := filepath.Join(vscodeDir, "tasks.json")
	os.WriteFile(tasksFile, tasksJSON, 0644)

	// Create a welcome file with instructions including performance tips
	welcomeContent := "# Welcome to Your Django Project!\n\n"
	
	if isUvAvailable() {
		welcomeContent += "🚀 **Performance Optimized**: This project uses `uv` for faster package management and execution.\n\n"
	}
	
	welcomeContent += "Two terminal windows should automatically open with:\n\n"

	if setupTailwind {
		welcomeContent += "1. Django development server (`" + activateCmd + pythonCmd + " manage.py runserver`)\n"
		welcomeContent += "2. Tailwind CSS watcher (`npm run watch:css`)\n\n"
	} else {
		welcomeContent += "1. Django development server (`" + activateCmd + pythonCmd + " manage.py runserver`)\n\n"
	}

	welcomeContent += "If the terminals didn't open automatically, you can run these tasks manually from the menu:\n"
	welcomeContent += "**Terminal > Run Task**\n\n"

	welcomeContent += "## Available VS Code Tasks:\n"
	welcomeContent += "- `" + djangoTaskLabel + "`: Start the Django development server\n"
	if setupTailwind {
		welcomeContent += "- `Tailwind: Watch CSS`: Watch and compile Tailwind CSS\n"
		welcomeContent += "- `Start Development Environment`: Run both Django and Tailwind concurrently\n"
	}
	
	if !isUvAvailable() {
		welcomeContent += "\n## Performance Tip:\n"
		welcomeContent += "💡 Install `uv` for faster package management and development server startup:\n"
		welcomeContent += "```bash\n"
		welcomeContent += "pip install uv\n"
		welcomeContent += "```\n"
	}

	welcomeFile := filepath.Join(projectPath, "WELCOME.md")
	os.WriteFile(welcomeFile, []byte(welcomeContent), 0644)
}

func (m *Model) setupServerInstructions(projectPath string) {
	if m.runServer {
		var startCommand string
		
		if isUvAvailable() {
			startCommand = "cd " + m.projectName + " && uv run python manage.py runserver"
		} else {
			pythonVenvPath := getPythonPath(projectPath)
			startCommand = "cd " + m.projectName + " && " + pythonVenvPath + " manage.py runserver"
		}
		
		m.stepMessages = append(m.stepMessages, "✨ To start the server: "+startCommand)
		
		if m.setupTailwind {
			m.stepMessages = append(m.stepMessages, "✨ To watch Tailwind CSS: cd "+m.projectName+" && npm run watch:css")
		}
	}
}