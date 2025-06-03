package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// setupGlobalTemplates creates global templates and static files
func (m *Model) setupGlobalTemplates(projectPath string) error {
	// Global templates directory
	globalTemplatesPath := filepath.Join(projectPath, "templates")
	if err := os.MkdirAll(globalTemplatesPath, 0755); err != nil {
		return fmt.Errorf("failed to create global templates directory: %v", err)
	}

	// Global static directory and subdirectories
	staticPath := filepath.Join(projectPath, "static")
	if err := os.MkdirAll(filepath.Join(staticPath, "css"), 0755); err != nil {
		return fmt.Errorf("failed to create static/css directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(staticPath, "js"), 0755); err != nil {
		return fmt.Errorf("failed to create static/js directory: %v", err)
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
		return fmt.Errorf("failed to create base.html: %v", err)
	}

	indexContent := `{% extends 'base.html' %}
{% block title %}Home{% endblock %}
{% block content %}
    <h1>Welcome to {{ project_name }}!</h1>
    <p>Your Django project is ready.</p>
{% endblock %}`
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "index.html"), []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.html: %v", err)
	}

	styleContent := `body { font-family: sans-serif; margin: 20px; background-color: #f4f4f4; color: #333; } h1 { color: #2c3e50; }`
	if err := os.WriteFile(filepath.Join(staticPath, "css", "style.css"), []byte(styleContent), 0644); err != nil {
		return fmt.Errorf("failed to create style.css: %v", err)
	}

	jsContent := `console.log('Django project initialized!');`
	if err := os.WriteFile(filepath.Join(staticPath, "js", "main.js"), []byte(jsContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.js: %v", err)
	}

	m.stepMessages = append(m.stepMessages, "✅ Created global templates and static files.")
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.85, status: "Setting up global templates and static files..."})
	}

	return nil
}

// updateSettingsForTemplates updates Django settings for templates and static files
func updateSettingsForTemplates(settingsContent string) string {
	// Update DIRS in TEMPLATES setting
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
				insertPos := idx + lineEndIdx + 1
				settingsContent = settingsContent[:insertPos] + "\n" + staticfilesDirsSetting + "\n" + settingsContent[insertPos:]
			} else { // STATIC_URL is the last line
				settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
			}
		} else { // Fallback: add to the end
			settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
		}
	}

	return settingsContent
}