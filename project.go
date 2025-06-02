package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *Model) CreateProject() {
	var currentErr error
	defer func() {
		if m.program != nil {
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
		m.program.Send(projectProgressMsg{percent: 0.05, status: "Creating project directory..."})
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
		m.program.Send(projectProgressMsg{percent: 0.15, status: "Virtual environment created."})
	}

	// Install Django and django-browser-reload
	djangoInstallVersion := m.djangoVersion
	if djangoInstallVersion == "" || djangoInstallVersion == "latest" {
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
		m.program.Send(projectProgressMsg{percent: 0.35, status: "Installing Django and dependencies..."})
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
		m.program.Send(projectProgressMsg{percent: 0.55, status: "Creating Django project structure..."})
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
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.65, status: "Configuring settings.py..."})
	}

	// Add django_browser_reload.middleware.BrowserReloadMiddleware to MIDDLEWARE
	settingsContent, err = addToListInSettingsPy(settingsContent, "MIDDLEWARE", "django_browser_reload.middleware.BrowserReloadMiddleware")
	if err != nil {
		currentErr = fmt.Errorf("failed to add BrowserReloadMiddleware to MIDDLEWARE: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "✅ Added django-browser-reload middleware.")
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.75, status: "Configuring middleware..."})
	}

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
		if m.program != nil {
			m.program.Send(projectProgressMsg{percent: 0.85, status: "Setting up global templates and static files..."})
		}

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
		m.stepMessages = append(m.stepMessages, "✅ Configured settings for global templates and static files.")
	}

	// Write updated settings.py
	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		currentErr = fmt.Errorf("failed to write updated settings.py: %v", err)
		return
	}
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: 0.95, status: "Finalizing settings configuration..."})
	}

	// Handle app creation if specified
	if m.appName != "" {
		pythonVenvPath := getPythonPath(projectPath)
		cmd = exec.Command(pythonVenvPath, "manage.py", "startapp", m.appName)
		cmd.Dir = projectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			currentErr = fmt.Errorf("failed to create app '%s': %v\nOutput: %s", m.appName, err, string(output))
			return
		}

		// Add app to INSTALLED_APPS
		settingsContentBytes, err := os.ReadFile(settingsPath)
		if err != nil {
			currentErr = fmt.Errorf("failed to read settings.py to add app: %v", err)
			return
		}
		settingsContent = string(settingsContentBytes)
		updatedSettings, err := addToListInSettingsPy(settingsContent, "INSTALLED_APPS", m.appName)
		if err != nil {
			currentErr = fmt.Errorf("failed to add app '%s' to INSTALLED_APPS: %v", m.appName, err)
			return
		}
		if err := os.WriteFile(settingsPath, []byte(updatedSettings), 0644); err != nil {
			currentErr = fmt.Errorf("failed to write updated settings.py after adding app: %v", err)
			return
		}
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))

		// Handle app templates if selected
		if m.createAppTemplates {
			appPath := filepath.Join(projectPath, m.appName)
			appTemplatesDir := filepath.Join(appPath, "templates", m.appName)
			if err := os.MkdirAll(appTemplatesDir, 0755); err != nil {
				currentErr = fmt.Errorf("failed to create app templates directory %s: %v", appTemplatesDir, err)
				return
			}

			appIndexContent := `{% extends 'base.html' %}
{% block title %}` + strings.Title(m.appName) + ` Home{% endblock %}
{% block content %}<h1>Welcome to the ` + m.appName + `</h1>{% endblock %}`
			if _, err := os.Stat(filepath.Join(projectPath, "templates", "base.html")); os.IsNotExist(err) {
				appIndexContent = `<!DOCTYPE html><html><head><title>` + strings.Title(m.appName) + `</title></head><body><h1>Welcome to the ` + m.appName + ` app!</h1></body></html>`
			}
			if err := os.WriteFile(filepath.Join(appTemplatesDir, "index.html"), []byte(appIndexContent), 0644); err != nil {
				currentErr = fmt.Errorf("failed to create index.html for app %s: %v", m.appName, err)
				return
			}

			// Create views.py
			viewsContent := fmt.Sprintf(`from django.shortcuts import render

def index(request):
    return render(request, '%s/index.html')
`, m.appName)
			if err := os.WriteFile(filepath.Join(appPath, "views.py"), []byte(viewsContent), 0644); err != nil {
				currentErr = fmt.Errorf("failed to create views.py for app %s: %v", m.appName, err)
				return
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
				currentErr = fmt.Errorf("failed to create urls.py for app %s: %v", m.appName, err)
				return
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
				currentErr = fmt.Errorf("failed to update project urls.py: %v", err)
				return
			}
			m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Configured templates, views, and URLs for app: %s", m.appName))
		}
	}

	// Handle Git initialization if selected
	if m.initializeGit {
		gitCmd := exec.Command("git", "init")
		gitCmd.Dir = projectPath
		if output, err := gitCmd.CombinedOutput(); err != nil {
			currentErr = fmt.Errorf("failed to initialize Git repository: %v\nOutput: %s", err, string(output))
			return
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
.venv/
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
		if err := os.WriteFile(filepath.Join(projectPath, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
			currentErr = fmt.Errorf("failed to create .gitignore: %v", err)
			return
		}
		m.stepMessages = append(m.stepMessages, "✅ .gitignore file created.")
	}

	// Setup Tailwind CSS v4 if selected
	if m.setupTailwind {
		if m.program != nil {
			m.program.Send(projectProgressMsg{percent: 0.90, status: "Setting up Tailwind CSS v4..."})
		}

		// Check if Node.js is available
		if !isCommandAvailable("npm") {
			m.stepMessages = append(m.stepMessages, "⚠️  Warning: npm not found. Please install Node.js to use Tailwind CSS.")
		} else {
			// Initialize npm
			cmd := exec.Command("npm", "init", "-y")
			cmd.Dir = projectPath
			if output, err := cmd.CombinedOutput(); err != nil {
				m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to initialize npm: %v\nOutput: %s", err, string(output)))
			} else {
				m.stepMessages = append(m.stepMessages, "✅ npm initialized.")

				// Install Tailwind CSS v4
				cmd = exec.Command("npm", "install", "tailwindcss", "@tailwindcss/cli")
				cmd.Dir = projectPath
				if output, err := cmd.CombinedOutput(); err != nil {
					m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to install Tailwind CSS: %v\nOutput: %s", err, string(output)))
				} else {
					m.stepMessages = append(m.stepMessages, "✅ Tailwind CSS v4 installed.")

					// Create static/src directory structure
					staticSrcPath := filepath.Join(projectPath, "static", "src")
					staticDistPath := filepath.Join(projectPath, "static", "dist")
					if err := os.MkdirAll(staticSrcPath, 0755); err != nil {
						m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create static/src directory: %v", err))
					} else if err := os.MkdirAll(staticDistPath, 0755); err != nil {
						m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create static/dist directory: %v", err))
					} else {
						m.stepMessages = append(m.stepMessages, "✅ Tailwind directory structure created.")

						// Create source CSS file
						tailwindCSS := `@import "tailwindcss";`
						if err := os.WriteFile(filepath.Join(staticSrcPath, "styles.css"), []byte(tailwindCSS), 0644); err != nil {
							m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to create styles.css: %v", err))
						} else {
							m.stepMessages = append(m.stepMessages, "✅ Tailwind source CSS created.")

							// Update package.json with build scripts
							packageJSONPath := filepath.Join(projectPath, "package.json")
							if packageData, err := os.ReadFile(packageJSONPath); err != nil {
								m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to read package.json: %v", err))
							} else {
								// Parse and update package.json
								var packageJSON map[string]interface{}
								if err := json.Unmarshal(packageData, &packageJSON); err != nil {
									m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to parse package.json: %v", err))
								} else {
									// Add scripts
									scripts := map[string]interface{}{
										"build:css": fmt.Sprintf("npx tailwindcss -i ./static/src/styles.css -o ./static/dist/styles.css"),
										"watch:css": fmt.Sprintf("npx tailwindcss -i ./static/src/styles.css -o ./static/dist/styles.css --watch"),
									}
									packageJSON["scripts"] = scripts

									// Write updated package.json
									if updatedData, err := json.MarshalIndent(packageJSON, "", "  "); err != nil {
										m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to marshal package.json: %v", err))
									} else if err := os.WriteFile(packageJSONPath, updatedData, 0644); err != nil {
										m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to write package.json: %v", err))
									} else {
										m.stepMessages = append(m.stepMessages, "✅ package.json updated with Tailwind scripts.")

										// Update base.html template to include Tailwind CSS
										if m.createTemplates {
											baseTemplatePath := filepath.Join(projectPath, "templates", "base.html")
											if baseContent, err := os.ReadFile(baseTemplatePath); err != nil {
												m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to read base.html: %v", err))
											} else {
												// Replace the CSS link with Tailwind CSS
												updatedBaseContent := strings.Replace(string(baseContent), 
													`<link rel="stylesheet" href="{% static 'css/style.css' %}">`,
													`<link rel="stylesheet" href="{% static 'dist/styles.css' %}">`, 1)
												if err := os.WriteFile(baseTemplatePath, []byte(updatedBaseContent), 0644); err != nil {
													m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to update base.html: %v", err))
												} else {
													m.stepMessages = append(m.stepMessages, "✅ base.html updated to use Tailwind CSS.")
												}
											}
										}

										// Build initial CSS
										cmd = exec.Command("npm", "run", "build:css")
										cmd.Dir = projectPath
										if output, err := cmd.CombinedOutput(); err != nil {
											m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to build Tailwind CSS: %v\nOutput: %s", err, string(output)))
										} else {
											m.stepMessages = append(m.stepMessages, "✅ Tailwind CSS compiled successfully.")
											m.stepMessages = append(m.stepMessages, "💡 Run 'npm run watch:css' for development or 'npm run build:css' for production.")
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Handle server startup if selected
	if m.runServer {
		pythonVenvPath := getPythonPath(projectPath)
		m.stepMessages = append(m.stepMessages, "✨ To start the server: cd "+m.projectName+" && "+pythonVenvPath+" manage.py runserver")
		if m.setupTailwind {
			m.stepMessages = append(m.stepMessages, "✨ To watch Tailwind CSS: cd "+m.projectName+" && npm run watch:css")
		}
	}

	m.stepMessages = append(m.stepMessages, "✅ Django project setup complete!")
	// Final progress update will be handled by projectCreationDoneMsg logic in Update()
}
