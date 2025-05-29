package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Global exit, regardless of state (unless error/done message is being shown)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC || keyMsg.String() == "q" {
			if !m.done && m.error == nil {
				return m, tea.Quit
			}
		}
	}

	if m.error != nil || m.done {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyEnter || keyMsg.Type == tea.KeyCtrlC || keyMsg.String() == "q" || keyMsg.String() == "esc" {
				return m, tea.Quit
			}
		}
		return m, nil
	}

	// Process general messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		// Adjust progress bar width dynamically
		m.progress.Width = msg.Width / 2
		if m.progress.Width < 10 {
			m.progress.Width = 10
		}
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}

	case tickMsg: // For splash screen countdown
		if m.step == stepSplashScreen {
			m.splashCountdown--
			if m.splashCountdown <= 0 {
				m.step = stepProjectName
				cmds = append(cmds, m.inputForm.Init())
			} else {
				cmds = append(cmds, tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
					return tickMsg{}
				}))
			}
		}
		return m, tea.Batch(cmds...)

	case projectProgressMsg:
		if m.step == stepSetup {
			m.progress.SetPercent(msg.percent)
			m.progressStatus = msg.status
			m.stepMessages = append(m.stepMessages, "PROGRESS: "+msg.status)
		}
		return m, m.spinner.Tick

	case projectCreationDoneMsg:
		if m.step == stepSetup {
			if msg.err != nil {
				m.error = msg.err
				m.progressStatus = "Error during project setup!"
				return m, nil
			}
			m.progress.SetPercent(1.0)
			m.progressStatus = "Django project core setup complete!"
			m.stepMessages = append(m.stepMessages, "✅ Django project core setup complete!")
			m.step = stepCreateApp
			cmds = append(cmds, m.appForm.Init())
		}
		return m, tea.Batch(cmds...)
	}

	activeForm := m.getActiveForm()
	if activeForm != nil {
		if activeForm.State != huh.StateCompleted {
			formModel, formCmd := activeForm.Update(msg)
			if castedForm, ok := formModel.(*huh.Form); ok {
				m.setActiveForm(castedForm)
			}
			cmds = append(cmds, formCmd)
		}
	}

	switch m.step {
	case stepProjectName:
		if m.inputForm.State == huh.StateCompleted {
			m.step = stepDjangoVersion
			m.stepMessages = append(m.stepMessages, "Project name: "+m.projectName)
			cmds = append(cmds, m.versionForm.Init())
		}
	case stepDjangoVersion:
		if m.versionForm.State == huh.StateCompleted {
			if m.djangoVersion == "" {
				m.djangoVersion = "latest"
			}
			m.step = stepFeatures
			m.stepMessages = append(m.stepMessages, "Django version: "+m.djangoVersion)
			cmds = append(cmds, m.featureForm.Init())
		}
	case stepFeatures:
		if m.featureForm.State == huh.StateCompleted {
			m.step = stepTemplates
			m.stepMessages = append(m.stepMessages, "Features: "+fmt.Sprint(m.features))
			cmds = append(cmds, m.templateForm.Init())
		}
	case stepTemplates:
		if m.templateForm.State == huh.StateCompleted {
			m.step = stepSetup
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("Global templates/static setup: %v", m.createTemplates))
			m.progressStatus = "Starting project setup..."
			go m.CreateProject()
			cmds = append(cmds, m.spinner.Tick)
		}

	case stepCreateApp:
		if m.appForm.State == huh.StateCompleted {
			if m.appName != "" {
				projectAbsPath, err := filepath.Abs(m.projectName)
				if err != nil {
					m.error = fmt.Errorf("could not get project absolute path: %v", err)
					break // Break from switch, error will be displayed
				}
				pythonVenvPath := getPythonPath(projectAbsPath)
				cmd := exec.Command(pythonVenvPath, "manage.py", "startapp", m.appName)
				cmd.Dir = projectAbsPath
				if output, err := cmd.CombinedOutput(); err != nil {
					m.error = fmt.Errorf("failed to create app '%s': %v\nOutput: %s", m.appName, err, string(output))
					break
				}

				settingsPath := filepath.Join(projectAbsPath, m.projectName, "settings.py")
				settingsContentBytes, err := os.ReadFile(settingsPath)
				if err != nil {
					m.error = fmt.Errorf("failed to read settings.py to add app: %v", err)
					break
				}
				settingsContent := string(settingsContentBytes)
				updatedSettings, err := addToListInSettingsPy(settingsContent, "INSTALLED_APPS", m.appName)
				if err != nil {
					m.error = fmt.Errorf("failed to add app '%s' to INSTALLED_APPS: %v", m.appName, err)
					break
				}
				if err := os.WriteFile(settingsPath, []byte(updatedSettings), 0644); err != nil {
					m.error = fmt.Errorf("failed to write updated settings.py after adding app: %v", err)
					break
				}
				m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))
				m.step = stepAppTemplates
				// Update the title of the appTemplateSelect field before initializing its form
				m.appTemplateSelect.Title(fmt.Sprintf("App Templates for '%s'", m.appName))
				cmds = append(cmds, m.appTemplateForm.Init())
			} else {
				m.stepMessages = append(m.stepMessages, "⏩ Skipped app creation.")
				m.step = stepServerOption
				cmds = append(cmds, m.serverForm.Init())
			}
		}
	case stepAppTemplates:
		if m.appTemplateForm.State == huh.StateCompleted {
			if m.appName != "" && m.createAppTemplates {
				projectAbsPath, _ := filepath.Abs(m.projectName)
				appPath := filepath.Join(projectAbsPath, m.appName)
				appTemplatesDir := filepath.Join(appPath, "templates", m.appName)
				if err := os.MkdirAll(appTemplatesDir, 0755); err != nil {
					m.error = fmt.Errorf("failed to create app templates directory %s: %v", appTemplatesDir, err)
					break
				}

				appIndexContent := `{% extends 'base.html' %}
{% block title %}` + strings.Title(m.appName) + ` Home{% endblock %}
{% block content %}<h1>Welcome to the ` + m.appName + ` app!</h1>{% endblock %}`
				if _, err := os.Stat(filepath.Join(projectAbsPath, "templates", "base.html")); os.IsNotExist(err) {
					appIndexContent = `<!DOCTYPE html><html><head><title>` + strings.Title(m.appName) + `</title></head><body><h1>Welcome to the ` + m.appName + ` app!</h1></body></html>`
				}
				if err := os.WriteFile(filepath.Join(appTemplatesDir, "index.html"), []byte(appIndexContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create index.html for app %s: %v", m.appName, err)
					break
				}

				// Fixed viewsContent with correct formatting
				viewsContent := fmt.Sprintf(`from django.shortcuts import render

def index(request):
    return render(request, '%s/index.html')
`, m.appName)

				if err := os.WriteFile(filepath.Join(appPath, "views.py"), []byte(viewsContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create views.py for app %s: %v", m.appName, err)
					break
				}

				appUrlsContent := fmt.Sprintf(`from django.urls import path
from . import views

app_name = '%s'
urlpatterns = [
    path('', views.index, name='index'),
]
`, m.appName)
				if err := os.WriteFile(filepath.Join(appPath, "urls.py"), []byte(appUrlsContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create urls.py for app %s: %v", m.appName, err)
					break
				}

				projectUrlsPath := filepath.Join(projectAbsPath, m.projectName, "urls.py")
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
					m.error = fmt.Errorf("failed to update project urls.py: %v", err)
					break
				}
				m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Configured templates, views, and URLs for app: %s", m.appName))
			}
			m.step = stepServerOption
			cmds = append(cmds, m.serverForm.Init())
		}

	case stepServerOption:
		if m.serverForm.State == huh.StateCompleted {
			m.step = stepGitInit
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("Run server choice: %v", m.runServer))
			cmds = append(cmds, m.gitForm.Init())
		}
	case stepGitInit:
		if m.gitForm.State == huh.StateCompleted {
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("Initialize Git choice: %v", m.initializeGit))
			projectAbsPath, _ := filepath.Abs(m.projectName)

			if m.initializeGit {
				gitCmd := exec.Command("git", "init")
				gitCmd.Dir = projectAbsPath
				if output, err := gitCmd.CombinedOutput(); err != nil {
					m.error = fmt.Errorf("failed to initialize Git repository: %v\nOutput: %s", err, string(output))
					break
				}
				m.stepMessages = append(m.stepMessages, "✅ Git repository initialized.")

				gitignoreContent := `# Django
*.log
*.pot
*.pyc
__pycache__/
local_settings.py
db.sqlite3
db.sqlite3-journal
media

# Virtual environment
venv/
env/
ENV/

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
`
				if err := os.WriteFile(filepath.Join(projectAbsPath, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create .gitignore: %v", err)
					break
				}
				m.stepMessages = append(m.stepMessages, "✅ .gitignore file created.")
			}

			if m.runServer {
				pythonVenvPath := getPythonPath(projectAbsPath)
				m.stepMessages = append(m.stepMessages, "✨ To start the server: cd "+m.projectName+" && "+pythonVenvPath+" manage.py runserver")
			}
			m.done = true
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) getActiveForm() *huh.Form {
	switch m.step {
	case stepProjectName:
		return m.inputForm
	case stepDjangoVersion:
		return m.versionForm
	case stepFeatures:
		return m.featureForm
	case stepTemplates:
		return m.templateForm
	case stepCreateApp:
		return m.appForm
	case stepAppTemplates:
		return m.appTemplateForm
	case stepServerOption:
		return m.serverForm
	case stepGitInit:
		return m.gitForm
	}
	return nil
}

func (m *Model) setActiveForm(form *huh.Form) {
	if form == nil {
		return
	}
	switch m.step {
	case stepProjectName:
		m.inputForm = form
	case stepDjangoVersion:
		m.versionForm = form
	case stepFeatures:
		m.featureForm = form
	case stepTemplates:
		m.templateForm = form
	case stepCreateApp:
		m.appForm = form
	case stepAppTemplates:
		m.appTemplateForm = form
	case stepServerOption:
		m.serverForm = form
	case stepGitInit:
		m.gitForm = form
	}
}

func Ternary[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
