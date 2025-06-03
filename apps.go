package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// createDjangoApp creates a Django app and configures it
func (m *Model) createDjangoApp(projectPath, settingsPath string) error {
	if m.appName == "" {
		return nil
	}

	pythonVenvPath := getPythonPath(projectPath)
	cmd := exec.Command(pythonVenvPath, "manage.py", "startapp", m.appName)
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create app '%s': %v\nOutput: %s", m.appName, err, string(output))
	}

	// Add app to INSTALLED_APPS
	settingsContentBytes, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings.py to add app: %v", err)
	}
	settingsContent := string(settingsContentBytes)
	updatedSettings, err := addToListInSettingsPy(settingsContent, "INSTALLED_APPS", m.appName)
	if err != nil {
		return fmt.Errorf("failed to add app '%s' to INSTALLED_APPS: %v", m.appName, err)
	}
	if err := os.WriteFile(settingsPath, []byte(updatedSettings), 0644); err != nil {
		return fmt.Errorf("failed to write updated settings.py after adding app: %v", err)
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))

	// Handle app templates if selected
	if m.createAppTemplates {
		if err := m.setupAppTemplates(projectPath); err != nil {
			return err
		}
	}

	return nil
}

// setupAppTemplates creates app-specific templates and views
func (m *Model) setupAppTemplates(projectPath string) error {
	appPath := filepath.Join(projectPath, m.appName)
	appTemplatesDir := filepath.Join(appPath, "templates", m.appName)
	if err := os.MkdirAll(appTemplatesDir, 0755); err != nil {
		return fmt.Errorf("failed to create app templates directory %s: %v", appTemplatesDir, err)
	}

	appIndexContent := `{% extends 'base.html' %}
{% block title %}` + strings.Title(m.appName) + ` Home{% endblock %}
{% block content %}<h1>Welcome to the ` + m.appName + `</h1>{% endblock %}`
	if _, err := os.Stat(filepath.Join(projectPath, "templates", "base.html")); os.IsNotExist(err) {
		appIndexContent = `<!DOCTYPE html><html><head><title>` + strings.Title(m.appName) + `</title></head><body><h1>Welcome to the ` + m.appName + ` app!</h1></body></html>`
	}
	if err := os.WriteFile(filepath.Join(appTemplatesDir, "index.html"), []byte(appIndexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.html for app %s: %v", m.appName, err)
	}

	// Create views.py
	viewsContent := fmt.Sprintf(`from django.shortcuts import render

def index(request):
    return render(request, '%s/index.html')
`, m.appName)
	if err := os.WriteFile(filepath.Join(appPath, "views.py"), []byte(viewsContent), 0644); err != nil {
		return fmt.Errorf("failed to create views.py for app %s: %v", m.appName, err)
	}

	// Create urls.py for app
	appUrlsContent := fmt.Sprintf(`from django.urls import path
from . import views

app_name = '%s'
urlpatterns = [
    path('', views.index, name='index'),
]
`, m.appName)
	if err := os.WriteFile(filepath.Join(appPath, "urls.py"), []byte(appUrlsContent), 0644); err != nil {
		return fmt.Errorf("failed to create urls.py for app %s: %v", m.appName, err)
	}

	// Update project urls.py
	projectUrlsPath := filepath.Join(projectPath, m.projectName, "urls.py")
	rootPathForProjectUrls := ""
	if m.createTemplates {
		rootPathForProjectUrls = "    path('', TemplateView.as_view(template_name='index.html'), name='home'),\n"
	}
	appIncludePath := fmt.Sprintf("    path('%s/', include('%s.urls', namespace='%s')),\n", m.appName, m.appName, m.appName)
	projectUrlsContent := fmt.Sprintf(`from django.contrib import admin
from django.urls import path, include
%s
urlpatterns = [
    path('admin/', admin.site.urls),
    path('__reload__/', include('django_browser_reload.urls')),
%s%s]
`, Ternary(m.createTemplates, "from django.views.generic import TemplateView", ""), appIncludePath, rootPathForProjectUrls)
	if err := os.WriteFile(projectUrlsPath, []byte(projectUrlsContent), 0644); err != nil {
		return fmt.Errorf("failed to update project urls.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Configured templates, views, and URLs for app: %s", m.appName))

	return nil
}