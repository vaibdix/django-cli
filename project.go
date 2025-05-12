package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *Model) createProject() {
	if m.projectName == "" {
		m.error = fmt.Errorf("project name cannot be empty")
		return
	}

	projectPath := m.projectName
	if !filepath.IsAbs(projectPath) {
		wd, err := os.Getwd()
		if err != nil {
			m.error = fmt.Errorf("failed to get working directory: %v", err)
			return
		}
		projectPath = filepath.Join(wd, m.projectName)
	}

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		m.error = fmt.Errorf("failed to create project directory: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "Project directory created.")

	m.progressStatus = "Creating virtual environment..."
	cmd := exec.Command("uv", "venv", ".venv")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to create virtual environment: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "Virtual environment created.")

	version := m.djangoVersion
	if version == "" {
		version = "5.2.0"
	}
	m.progressStatus = fmt.Sprintf("Installing Django %s...", version)
	cmd = exec.Command("uv", "pip", "install", "django=="+version)
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to install Django: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Django %s installed.", version))

	m.progressStatus = "Creating Django project..."
	pythonPath := filepath.Join(projectPath, ".venv", "bin", "python")
	cmd = exec.Command(pythonPath, "-m", "django", "startproject", m.projectName, ".")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to create Django project: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "Django project created.")
	m.stepMessages = append(m.stepMessages, "Using vanilla Django setup")

	if m.createTemplates {
		globalTemplatesPath := filepath.Join(projectPath, "templates")
		if err := os.MkdirAll(globalTemplatesPath, 0755); err != nil {
			m.error = fmt.Errorf("failed to create global templates directory: %v", err)
			return
		}

		staticPath := filepath.Join(projectPath, "static")
		if err := os.MkdirAll(filepath.Join(staticPath, "css"), 0755); err != nil {
			m.error = fmt.Errorf("failed to create static/css directory: %v", err)
			return
		}
		if err := os.MkdirAll(filepath.Join(staticPath, "js"), 0755); err != nil {
			m.error = fmt.Errorf("failed to create static/js directory: %v", err)
			return
		}

		styleContent := `/* Global styles */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: Arial, sans-serif;
    line-height: 1.6;
    padding: 20px;
}

h1 {
    color: red;
    margin-bottom: 20px;
}`
		stylePath := filepath.Join(staticPath, "css", "style.css")
		if err := os.WriteFile(stylePath, []byte(styleContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to create style.css: %v", err)
			return
		}

		jsContent := `// Main JavaScript file
document.addEventListener('DOMContentLoaded', function() {
    console.log('Django project initialized!');
});
`
		jsPath := filepath.Join(staticPath, "js", "main.js")
		if err := os.WriteFile(jsPath, []byte(jsContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to create main.js: %v", err)
			return
		}

		baseContent := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>{{ title|default:"Django Project" }}</title>
			{% load static %}
			<link rel="stylesheet" href="{% static 'css/style.css' %}">
		</head>
		<body>
			{% block content %}
			{% endblock %}
			<script src="{% static 'js/main.js' %}"></script>
		</body>
		</html>`
		basePath := filepath.Join(globalTemplatesPath, "base.html")
		if err := os.WriteFile(basePath, []byte(baseContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to create base.html: %v", err)
			return
		}

		indexContent := `{% extends 'base.html' %}
		{% block content %}
		<h1>Welcome to {{ project_name }}</h1>
		{% endblock %}`

		indexPath := filepath.Join(globalTemplatesPath, "index.html")
		if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to create index.html: %v", err)
			return
		}

		m.stepMessages = append(m.stepMessages, "✅ Created global templates directory with base.html and index.html")

		settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
		settingsContent, err := os.ReadFile(settingsPath)
		if err != nil {
			m.error = fmt.Errorf("failed to read settings.py: %v", err)
			return
		}

		settingsStr := string(settingsContent)
		templatesIndex := strings.Index(settingsStr, "'DIRS': []")
		if templatesIndex == -1 {
			m.error = fmt.Errorf("could not find TEMPLATES setting in settings.py")
			return
		}

		newSettingsContent := strings.Replace(settingsStr, "'DIRS': []", "'DIRS': [BASE_DIR / 'templates',]", 1)

		if !strings.Contains(newSettingsContent, "import os") {
			importIndex := strings.Index(newSettingsContent, "from pathlib import Path")
			if importIndex != -1 {
				newSettingsContent = newSettingsContent[:importIndex] + "import os\n" + newSettingsContent[importIndex:]
			}
		}

		if !strings.Contains(newSettingsContent, "STATICFILES_DIRS") {
			newSettingsContent += "\n# Static files (CSS, JavaScript, Images)\nSTATICFILES_DIRS = [\n    BASE_DIR / 'static',\n]\n"
		}

		if err := os.WriteFile(settingsPath, []byte(newSettingsContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to update settings.py: %v", err)
			return
		}
		m.stepMessages = append(m.stepMessages, "✅ Updated settings.py with templates configuration")
	}

	if m.appName != "" {
		wd, err := os.Getwd()
		if err != nil {
			m.error = fmt.Errorf("failed to get working directory: %v", err)
			return
		}
		projectPath := filepath.Join(wd, m.projectName)
		pythonPath := filepath.Join(projectPath, ".venv", "bin", "python")
		cmd := exec.Command(pythonPath, "manage.py", "startapp", m.appName)
		cmd.Dir = projectPath
		if err := cmd.Run(); err != nil {
			m.error = fmt.Errorf("failed to create app: %v", err)
			return
		}

		if m.createTemplates {
			appTemplatesPath := filepath.Join(projectPath, "templates", m.appName)
			if err := os.MkdirAll(appTemplatesPath, 0755); err != nil {
				m.error = fmt.Errorf("failed to create app templates directory: %v", err)
				return
			}
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created templates directory for %s app", m.appName))
		}

		settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
		settingsContent, err := os.ReadFile(settingsPath)
		if err != nil {
			m.error = fmt.Errorf("failed to read settings.py: %v", err)
			return
		}

		settingsStr := string(settingsContent)
		installedAppsIndex := strings.Index(settingsStr, "INSTALLED_APPS = [")
		if installedAppsIndex == -1 {
			m.error = fmt.Errorf("could not find INSTALLED_APPS in settings.py")
			return
		}

		closeBracketIndex := strings.Index(settingsStr[installedAppsIndex:], "]")
		if closeBracketIndex == -1 {
			m.error = fmt.Errorf("malformed INSTALLED_APPS in settings.py")
			return
		}

		newSettingsContent := settingsStr[:installedAppsIndex+closeBracketIndex] +
			"    '" + m.appName + "',\n" +
			settingsStr[installedAppsIndex+closeBracketIndex:]

		if err := os.WriteFile(settingsPath, []byte(newSettingsContent), 0644); err != nil {
			m.error = fmt.Errorf("failed to update settings.py: %v", err)
			return
		}

		m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))
	}

	m.stepMessages = append(m.stepMessages, "✅ Project setup finished!")
	m.doneChan <- true
}