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
	if m.done {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tickMsg:
		if m.step == stepSplashScreen {
			m.splashCountdown--
			if m.splashCountdown <= 0 {
				m.step = stepProjectName
				return m, m.inputForm.Init()
			}
			return m, tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
				return tickMsg{}
			})
		}
	case splashDoneMsg:
		m.step = stepProjectName
		return m, m.inputForm.Init()
	case progressMsg:
		if float64(msg) >= 1.0 {
			m.progress.SetPercent(1.0)
			m.step = stepCreateApp
			return m, m.appForm.Init()
		}
		m.progress.SetPercent(float64(msg))
		return m, m.updateProgress()
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if msg.Width > 0 {
			m.width = msg.Width
		}
	}

	switch m.step {
	case stepProjectName:
		form, cmd := m.inputForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.inputForm = f
			if f.State == huh.StateCompleted {
				m.step = stepDjangoVersion
				m.stepMessages = append(m.stepMessages, "Project name selected: " + m.projectName)
				return m, m.versionForm.Init()
			}
			return m, cmd
		}
	case stepDjangoVersion:
		form, cmd := m.versionForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.versionForm = f
			if f.State == huh.StateCompleted {
				m.step = stepFeatures
				m.stepMessages = append(m.stepMessages, "Django version selected: " + m.djangoVersion)
				return m, m.featureForm.Init()
			}
			return m, cmd
		}
	case stepFeatures:
		form, cmd := m.featureForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.featureForm = f
			if f.State == huh.StateCompleted {
				m.step = stepTemplates
				m.stepMessages = append(m.stepMessages, "Features selected: " + fmt.Sprint(m.features))
				return m, m.templateForm.Init()
			}
			return m, cmd
		}
	case stepTemplates:
		form, cmd := m.templateForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.templateForm = f
			if f.State == huh.StateCompleted {
				m.step = stepSetup
				m.stepMessages = append(m.stepMessages, fmt.Sprintf("Templates setup: %v", m.createTemplates))
				go m.CreateProject()
				return m, tea.Batch(m.spinner.Tick, m.updateProgress())
			}
			return m, cmd
		}

	case stepCreateApp:
		form, cmd := m.appForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.appForm = f
			if f.State == huh.StateCompleted {
				if m.appName != "" {
					wd, err := os.Getwd()
					if err != nil {
						m.error = fmt.Errorf("failed to get working directory: %v", err)
						return m, nil
					}
					projectPath := filepath.Join(wd, m.projectName)
					pythonPath := filepath.Join(projectPath, ".venv", "bin", "python")
					cmd := exec.Command(pythonPath, "manage.py", "startapp", m.appName)
					cmd.Dir = projectPath
					if err := cmd.Run(); err != nil {
						m.error = fmt.Errorf("failed to create app: %v", err)
						return m, nil
					}

					settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
					settingsContent, err := os.ReadFile(settingsPath)
					if err != nil {
						m.error = fmt.Errorf("failed to read settings.py: %v", err)
						return m, nil
					}

					settingsStr := string(settingsContent)
					installedAppsIndex := strings.Index(settingsStr, "INSTALLED_APPS = [")
					if installedAppsIndex == -1 {
						m.error = fmt.Errorf("could not find INSTALLED_APPS in settings.py")
						return m, nil
					}

					closeBracketIndex := strings.Index(settingsStr[installedAppsIndex:], "]")
					if closeBracketIndex == -1 {
						m.error = fmt.Errorf("malformed INSTALLED_APPS in settings.py")
						return m, nil
					}

					if m.appName != "" {
						newSettingsContent := settingsStr[:installedAppsIndex+closeBracketIndex] +
							"    '" + m.appName + "',\n" +
							settingsStr[installedAppsIndex+closeBracketIndex:]

						if err := os.WriteFile(settingsPath, []byte(newSettingsContent), 0644); err != nil {
							m.error = fmt.Errorf("failed to update settings.py: %v", err)
							return m, nil
						}
					}

					m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))
					m.step = stepAppTemplates
					m.createAppTemplates = false // Reset the value
					return m, m.appForm.Init()
				} else {
					m.step = stepServerOption
					return m, m.serverForm.Init()
				}
			}
			return m, cmd
		}
	case stepAppTemplates:
		form, cmd := m.appTemplateForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.appTemplateForm = f
			if f.State == huh.StateCompleted {
				wd, err := os.Getwd()
				if err != nil {
					m.error = fmt.Errorf("failed to get working directory: %v", err)
					return m, nil
				}
				projectPath := filepath.Join(wd, m.projectName)
				appTemplatesPath := filepath.Join(projectPath, m.appName, "templates", m.appName)
				if err := os.MkdirAll(appTemplatesPath, 0755); err != nil {
					m.error = fmt.Errorf("failed to create app templates directory: %v", err)
					return m, nil
				}

				// Check for global base.html
				globalBasePath := filepath.Join(projectPath, "templates", "base.html")
				baseExists := true
				if _, err := os.Stat(globalBasePath); os.IsNotExist(err) {
					baseExists = false
				}

				var indexContent string
				if baseExists {
					// Extend the base template if it exists
					indexContent = `{% extends 'base.html' %}

                    {% block content %}
                    <div class="container mt-5">
                        <h1 class="display-4">Welcome to {{ app_name }}</h1>
                        <p class="lead">Your Django app is ready to go!</p>
                    </div>
                    {% endblock content %}`
				} else {
					// Use a complete HTML structure if base.html doesn't exist
					indexContent = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
</head>
<body>
    <div class="container mt-5">
        <h1 class="display-4">Welcome to {{ app_name }}</h1>
        <p class="lead">Your Django app is ready to go!</p>
    </div>
</body>
</html>`
				}

				indexPath := filepath.Join(appTemplatesPath, "index.html")
				if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create index.html: %v", err)
					return m, nil
				}

				viewsPath := filepath.Join(projectPath, m.appName, "views.py")
				viewsContent := `from django.shortcuts import render

def index(request):
    app_name = getattr(request.resolver_match, 'app_name', '` + m.appName + `')
    context = {
        'app_name': app_name
    }
    return render(request, f'{app_name}/index.html', context)`
				if err := os.WriteFile(viewsPath, []byte(viewsContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create views.py: %v", err)
					return m, nil
				}

				urlsPath := filepath.Join(projectPath, m.appName, "urls.py")
				urlsContent := fmt.Sprintf(`from django.urls import path

from . import views

app_name = '%s'
urlpatterns = [
    path('', views.index, name='index'),
]
`, m.appName)
				if err := os.WriteFile(urlsPath, []byte(urlsContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to create urls.py: %v", err)
					return m, nil
				}

				projectUrlsPath := filepath.Join(projectPath, m.projectName, "urls.py")
				projectUrlsContent := `from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('admin/', admin.site.urls),
    path('__reload__/', include('django_browser_reload.urls')),
    path('', include('` + m.appName + `.urls', namespace='` + m.appName + `')),
]`
				if err := os.WriteFile(projectUrlsPath, []byte(projectUrlsContent), 0644); err != nil {
					m.error = fmt.Errorf("failed to update project urls.py: %v", err)
					return m, nil
				}

				m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created templates directory and index.html for %s app", m.appName))

				m.step = stepServerOption
				return m, m.serverForm.Init()
			}
			return m, cmd
		}

	case stepServerOption:
		form, cmd := m.serverForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.serverForm = f
			if f.State == huh.StateCompleted {
				m.step = stepGitInit
				return m, m.gitForm.Init()
			}
			return m, cmd
		}
	case stepGitInit:
		form, cmd := m.gitForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.gitForm = f
			if f.State == huh.StateCompleted {
				wd, err := os.Getwd()
				if err != nil {
					m.error = fmt.Errorf("failed to get working directory: %v", err)
					return m, nil
				}
				projectPath := filepath.Join(wd, m.projectName)

				if m.initializeGit {
					// Initialize Git repository
					gitInitCmd := exec.Command("git", "init")
					gitInitCmd.Dir = projectPath
					if err := gitInitCmd.Run(); err != nil {
						m.error = fmt.Errorf("failed to initialize Git repository: %v", err)
						return m, nil
					}
					m.stepMessages = append(m.stepMessages, "✅ Git repository initialized.")

					// Create .gitignore file
					gitignoreContent := `# Django
*.log
*.pot
*.pyc
__pycache__/
db.sqlite3
/media/
/static/

# Environments
.env
.venv
venv/
ENV/
env/
ENV.yml

# IDE / Editor
.idea/
.vscode/
*.swp

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db`
					gitignorePath := filepath.Join(projectPath, ".gitignore")
					if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
						m.error = fmt.Errorf("failed to create .gitignore file: %v", err)
						return m, nil
					}
					m.stepMessages = append(m.stepMessages, "✅ .gitignore file created.")
				}

				if m.runServer {
					pythonPath := filepath.Join(projectPath, ".venv", "bin", "python")
					serverCmd := exec.Command(pythonPath, "manage.py", "runserver")
					serverCmd.Dir = projectPath
					serverCmd.Stdout = os.Stdout
					serverCmd.Stderr = os.Stderr
					if err := serverCmd.Start(); err != nil {
						m.error = fmt.Errorf("failed to start development server: %v", err)
						return m, nil
					}
					m.stepMessages = append(m.stepMessages, "✨ Development server started at http://127.0.0.1:8000/")
				}
				m.done = true
				return m, nil
			}
			return m, cmd
		}
	}

	return m, nil
}

