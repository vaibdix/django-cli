package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	// "time" // Not needed if progress is discrete
)

func (m *Model) CreateProject() {
	// This function will now send messages via m.program.Send()
	// instead of returning an error or writing to m.error directly in this goroutine.
	var currentErr error
	defer func() {
		if m.program != nil {
			// Send finalization message
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
		m.program.Send(projectProgressMsg{percent: 0.1, status: "Project directory created."})
	}

	// Create virtual environment
	pythonCmd := getPythonCommand()
	if pythonCmd == "" {
		currentErr = fmt.Errorf("no Python command (python3 or python) found in PATH")
		return
	}

	var cmd *exec.Cmd
	venvTool := ""
	if isCommandAvailable("uv") {
		cmd = exec.Command("uv", "venv", ".venv", "--python", pythonCmd)
		venvTool = "uv"
	} else {
		cmd = exec.Command(pythonCmd, "-m", "venv", ".venv")
		venvTool = "python -m venv"
	}
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		currentErr = fmt.Errorf("failed to create virtual environment using %s: %v\nOutput: %s", venvTool, err, string(output))
		return
	}
	m.stepMessages = append(m.stepMessages, "Virtual environment created.")
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.25, status: "Virtual environment created."})
	}

	// Install Django and django-browser-reload
	djangoInstallVersion := m.djangoVersion
	if djangoInstallVersion == "" || djangoInstallVersion == "latest" {
		// Fetch latest stable Django version here if desired, or use a fixed recent one.
		// For now, let Django/pip decide "latest" or use a default.
		// If you specify no version, pip installs the latest.
		// For explicit default:
		djangoInstallVersion = "Django" // Pip will get latest
		if m.djangoVersion != "" && m.djangoVersion != "latest" { // User specified a version
		    djangoInstallVersion = "Django==" + m.djangoVersion
        }
	} else {
		djangoInstallVersion = "Django==" + djangoInstallVersion
	}

	installCmdArgs := []string{"install", djangoInstallVersion, "django-browser-reload"}
	pipTool := ""
	if isCommandAvailable("uv") {
		cmd = exec.Command("uv", append([]string{"pip"}, installCmdArgs...)...)
		pipTool = "uv pip"
	} else {
		pipPath := getPipPath(projectPath)
		cmd = exec.Command(pipPath, installCmdArgs...)
		pipTool = "pip"
	}
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		currentErr = fmt.Errorf("failed to install Django and django-browser-reload using %s: %v\nOutput: %s", pipTool, err, string(output))
		return
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Django and django-browser-reload installed using %s.", pipTool))
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.5, status: "Django and dependencies installed."})
	}

	// Create Django project
	pythonVenvPath := getPythonPath(projectPath)
	cmd = exec.Command(pythonVenvPath, "-m", "django", "startproject", m.projectName, ".")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		currentErr = fmt.Errorf("failed to create Django project: %v\nOutput: %s", err, string(output))
		return
	}
	m.stepMessages = append(m.stepMessages, "Django project created.")
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.7, status: "Django project structure created."})
	}

	// Modify settings.py
	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
	settingsContentBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		currentErr = fmt.Errorf("failed to read settings.py: %v", err)
		return
	}
	settingsContent := string(settingsContentBytes)

	// Add django_browser_reload to INSTALLED_APPS
	settingsContent, err = addToListInSettingsPy(settingsContent, "INSTALLED_APPS", "django_browser_reload")
	if err != nil {
		currentErr = fmt.Errorf("failed to add django_browser_reload to INSTALLED_APPS: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "✅ Added django-browser-reload to INSTALLED_APPS.")

	// Add django_browser_reload.middleware.BrowserReloadMiddleware to MIDDLEWARE
	settingsContent, err = addToListInSettingsPy(settingsContent, "MIDDLEWARE", "django_browser_reload.middleware.BrowserReloadMiddleware")
	if err != nil {
		currentErr = fmt.Errorf("failed to add BrowserReloadMiddleware to MIDDLEWARE: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "✅ Added django-browser-reload middleware.")

	// Configure templates and static files if chosen
	if m.createTemplates {
		// Global templates directory
		globalTemplatesPath := filepath.Join(projectPath, "templates")
		if err := os.MkdirAll(globalTemplatesPath, 0755); err != nil {
			currentErr = fmt.Errorf("failed to create global templates directory: %v", err)
			return
		}

		// Global static directory and subdirectories
		staticPath := filepath.Join(projectPath, "static")
		if err := os.MkdirAll(filepath.Join(staticPath, "css"), 0755); err != nil {
			currentErr = fmt.Errorf("failed to create static/css directory: %v", err)
			return
		}
		if err := os.MkdirAll(filepath.Join(staticPath, "js"), 0755); err != nil {
			currentErr = fmt.Errorf("failed to create static/js directory: %v", err)
			return
		}
		// Create base.html, index.html, style.css, main.js
		baseContent := `{% load static %}
{% load django_browser_reload %}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{% block title %}My Django Site{% endblock %}</title>
    <link rel="stylesheet" href="{% static 'css/style.css' %}">
    {% block extra_head %}{% endblock %}
</head>
<body>
    <div class="container">
        {% block content %}{% endblock %}
    </div>
    <script src="{% static 'js/main.js' %}"></script>
    {{ django_browser_reload_script }}
</body>
</html>`
		if err := os.WriteFile(filepath.Join(globalTemplatesPath, "base.html"), []byte(baseContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to create base.html: %v", err)
			return
		}
		indexContent := `{% extends 'base.html' %}
{% block title %}Home{% endblock %}
{% block content %}
    <h1>Welcome to {{ project_name }}!</h1>
    <p>Your Django project is ready.</p>
{% endblock %}`
		if err := os.WriteFile(filepath.Join(globalTemplatesPath, "index.html"), []byte(indexContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to create index.html: %v", err)
			return
		}
		styleContent := `body { font-family: sans-serif; margin: 20px; background-color: #f4f4f4; color: #333; } h1 { color: #2c3e50; }`
		if err := os.WriteFile(filepath.Join(staticPath, "css", "style.css"), []byte(styleContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to create style.css: %v", err)
			return
		}
		jsContent := `console.log('Django project initialized!');`
		if err := os.WriteFile(filepath.Join(staticPath, "js", "main.js"), []byte(jsContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to create main.js: %v", err)
			return
		}
		m.stepMessages = append(m.stepMessages, "✅ Created global templates and static files.")

		// Update DIRS in TEMPLATES setting
		// This is a simplified replacement. A more robust method would parse the TEMPLATES structure.
		templatesDirsSetting := "'DIRS': [BASE_DIR / 'templates']"
		if !strings.Contains(settingsContent, templatesDirsSetting) {
			settingsContent = strings.Replace(settingsContent, "'DIRS': []", templatesDirsSetting, 1)
		}

		// Add STATICFILES_DIRS
		staticfilesDirsSetting := "STATICFILES_DIRS = [\n    BASE_DIR / 'static',\n]"
		if !strings.Contains(settingsContent, "STATICFILES_DIRS") {
			// Add near STATIC_URL or at the end of the file
			staticUrlMarker := "STATIC_URL = "
			idx := strings.Index(settingsContent, staticUrlMarker)
			if idx != -1 {
				// Find end of that line
				lineEndIdx := strings.Index(settingsContent[idx:], "\n")
				if lineEndIdx != -1 {
					insertPos := idx + lineEndIdx +1
					settingsContent = settingsContent[:insertPos] + "\n" + staticfilesDirsSetting + "\n" + settingsContent[insertPos:]
				} else { // STATIC_URL is the last line
					settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
				}
			} else { // Fallback: add to the end
				settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
			}
		}
		m.stepMessages = append(m.stepMessages, "✅ Configured settings for global templates and static files.")
	}

	// Write updated settings.py
	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		currentErr = fmt.Errorf("failed to write updated settings.py: %v", err)
		return
	}
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.9, status: "Settings configured."})
	}

	m.stepMessages = append(m.stepMessages, "✅ Django project core setup finished!")
	// Final progress update will be handled by projectCreationDoneMsg logic in Update()
}