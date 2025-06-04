package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (m *Model) setupGlobalTemplates(projectPath string) error {
	globalTemplatesPath := filepath.Join(projectPath, "templates")
	if err := os.MkdirAll(globalTemplatesPath, 0755); err != nil {
		return fmt.Errorf("failed to create global templates directory: %v", err)
	}
	staticPath := filepath.Join(projectPath, "static")
	if err := os.MkdirAll(filepath.Join(staticPath, "css"), 0755); err != nil {
		return fmt.Errorf("failed to create static/css directory: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(staticPath, "js"), 0755); err != nil {
		return fmt.Errorf("failed to create static/js directory: %v", err)
	}

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
{% block title %}Welcome - {{ project_name }}{% endblock %}

{% block content %}
<div class="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 flex items-center justify-center">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 w-full">
        <div class="text-center">
            <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl md:text-6xl">
                Welcome to {{ project_name|default:project_name }}
            </h1>
            <p class="mt-3 text-base text-gray-500 sm:text-lg md:mt-5 md:text-xl max-w-prose mx-auto">
                Your Django project is ready for development
            </p>
        </div>

        <div class="mt-16 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <!-- Documentation Card -->
            <a href="{% url 'api_docs' %}" class="group relative bg-white rounded-lg shadow-lg overflow-hidden transform transition duration-200 hover:scale-105">
                <div class="px-6 py-8">
                    <div class="text-center">
                        <div class="h-12 w-12 mx-auto bg-indigo-100 rounded-lg flex items-center justify-center">
                            <svg class="h-6 w-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                        </div>
                        <h3 class="mt-4 text-lg font-medium text-gray-900">Documentation</h3>
                        <p class="mt-2 text-sm text-gray-500">
                            View API endpoints, authentication, and admin documentation
                        </p>
                    </div>
                </div>
            </a>

            {% if app_name %}
            <!-- App Card -->
            <a href="/{{ app_name }}/" class="group relative bg-white rounded-lg shadow-lg overflow-hidden transform transition duration-200 hover:scale-105">
                <div class="px-6 py-8">
                    <div class="text-center">
                        <div class="h-12 w-12 mx-auto bg-green-100 rounded-lg flex items-center justify-center">
                            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                            </svg>
                        </div>
                        <h3 class="mt-4 text-lg font-medium text-gray-900">{{ app_name|title }} App</h3>
                        <p class="mt-2 text-sm text-gray-500">
                            Access your application's main page
                        </p>
                    </div>
                </div>
            </a>
            {% endif %}

            <!-- Quick Links Card -->
            <div class="group relative bg-white rounded-lg shadow-lg overflow-hidden">
                <div class="px-6 py-8">
                    <div class="text-center">
                        <div class="h-12 w-12 mx-auto bg-purple-100 rounded-lg flex items-center justify-center">
                            <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                        </div>
                        <h3 class="mt-4 text-lg font-medium text-gray-900">Quick Links</h3>
                        <div class="mt-4 space-y-2">
                            <a href="/admin/" class="block text-sm text-indigo-600 hover:text-indigo-500">
                                Admin Interface →
                            </a>
                            <a href="/api/v1/" class="block text-sm text-indigo-600 hover:text-indigo-500">
                                API Root →
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
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

	// Write the API docs template
	apiDocsContent := `{% extends 'base.html' %}
{% block title %}API Documentation - {{ project_name }}{% endblock %}

{% block content %}
<div class="min-h-screen bg-gray-50 py-12">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center">
            <h1 class="text-4xl font-bold text-gray-900">
                API Documentation
            </h1>
            <p class="mt-3 text-lg text-gray-500">
                Explore available endpoints and features
            </p>
        </div>

        <!-- API Endpoints Section -->
        <div class="mt-12 space-y-8">
            <div class="bg-white shadow overflow-hidden sm:rounded-lg">
                <div class="px-4 py-5 border-b border-gray-200">
                    <h2 class="text-xl font-semibold text-gray-900">
                        API Endpoints
                    </h2>
                </div>
                <div class="divide-y divide-gray-200">
                    <div class="px-4 py-5 sm:p-6">
                        <h3 class="text-lg font-medium text-gray-900">Books API</h3>
                        <div class="mt-4 space-y-4">
                            <div class="flex items-center justify-between">
                                <div class="flex items-center space-x-2">
                                    <span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-2 py-1 text-xs font-semibold text-blue-800 bg-blue-100 rounded-full">POST</span>
                                    <code class="ml-2 text-sm text-indigo-600">/api/v1/books/</code>
                                </div>
                                <span class="text-sm text-gray-500">List and create books</span>
                            </div>
                            <div class="flex items-center justify-between">
                                <div class="flex items-center space-x-2">
                                    <span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-2 py-1 text-xs font-semibold text-yellow-800 bg-yellow-100 rounded-full">PUT</span>
                                    <span class="px-2 py-1 text-xs font-semibold text-red-800 bg-red-100 rounded-full">DELETE</span>
                                    <code class="ml-2 text-sm text-indigo-600">/api/v1/books/{id}/</code>
                                </div>
                                <span class="text-sm text-gray-500">Manage individual books</span>
                            </div>
                            <div class="flex items-center justify-between">
                                <div class="flex items-center space-x-2">
                                    <span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <code class="ml-2 text-sm text-indigo-600">/api/v1/books/recent/</code>
                                </div>
                                <span class="text-sm text-gray-500">List recent books</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Authentication Section -->
            <div class="bg-white shadow overflow-hidden sm:rounded-lg">
                <div class="px-4 py-5 border-b border-gray-200">
                    <h2 class="text-xl font-semibold text-gray-900">
                        Authentication
                    </h2>
                </div>
                <div class="px-4 py-5 sm:p-6 space-y-4">
                    <div class="flex items-center justify-between">
                        <div class="flex items-center space-x-2">
                            <span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">GET</span>
                            <code class="ml-2 text-sm text-indigo-600">/api-auth/login/</code>
                        </div>
                        <span class="text-sm text-gray-500">API authentication login</span>
                    </div>
                    <div class="flex items-center justify-between">
                        <div class="flex items-center space-x-2">
                            <span class="px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded-full">GET</span>
                            <code class="ml-2 text-sm text-indigo-600">/api-auth/logout/</code>
                        </div>
                        <span class="text-sm text-gray-500">API authentication logout</span>
                    </div>
                </div>
            </div>
        </div>

        <div class="mt-8 text-center">
            <a href="/" class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700">
                Back to Home
            </a>
        </div>
    </div>
</div>
{% endblock %}`
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "api-docs.html"), []byte(apiDocsContent), 0644); err != nil {
		return fmt.Errorf("failed to create api-docs.html: %v", err)
	}

	m.stepMessages = append(m.stepMessages, "✅ Created global templates and static files.")
	m.updateProgress("Setting up global templates and static files...")
	return nil
}

func updateSettingsForTemplates(settingsContent string) string {
	templatesDirsSetting := "'DIRS': [BASE_DIR / 'templates']"
	if !strings.Contains(settingsContent, templatesDirsSetting) {
		settingsContent = strings.Replace(settingsContent, "'DIRS': []", templatesDirsSetting, 1)
	}
	staticfilesDirsSetting := "STATICFILES_DIRS = [\n    BASE_DIR / 'static',\n]"
	if !strings.Contains(settingsContent, "STATICFILES_DIRS") {
		staticUrlMarker := "STATIC_URL = "
		idx := strings.Index(settingsContent, staticUrlMarker)
		if idx != -1 {
			lineEndIdx := strings.Index(settingsContent[idx:], "\n")
			if lineEndIdx != -1 {
				insertPos := idx + lineEndIdx + 1
				settingsContent = settingsContent[:insertPos] + "\n" + staticfilesDirsSetting + "\n" + settingsContent[insertPos:]
			} else {
				settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
			}
		} else {
			settingsContent += "\n\n" + staticfilesDirsSetting + "\n"
		}
	}

	return settingsContent
}

func (m *Model) setupProjectUrls(projectPath string) error {
	// First create a context processor to make project_name available globally
	contextProcessorsPath := filepath.Join(projectPath, m.projectName, "context_processors.py")
	contextProcessorsContent := fmt.Sprintf(`def project_context(request):
    return {
        'project_name': '%s'
    }
`, m.projectName)

	if err := os.WriteFile(contextProcessorsPath, []byte(contextProcessorsContent), 0644); err != nil {
		return fmt.Errorf("failed to create context_processors.py: %v", err)
	}

	// Update settings.py to include the context processor
	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings.py: %v", err)
	}

	// Add context processor to templates settings
	updatedSettings := strings.Replace(
		string(settingsContent),
		"'django.contrib.messages.context_processors.messages',",
		"'django.contrib.messages.context_processors.messages',\n                '"+m.projectName+".context_processors.project_context',",
		1,
	)

	if err := os.WriteFile(settingsPath, []byte(updatedSettings), 0644); err != nil {
		return fmt.Errorf("failed to update settings.py: %v", err)
	}

	// Create views.py
	viewsContent := `from django.views.generic import TemplateView

class HomeView(TemplateView):
    template_name = 'index.html'`

	viewsPath := filepath.Join(projectPath, m.projectName, "views.py")
	if err := os.WriteFile(viewsPath, []byte(viewsContent), 0644); err != nil {
		return fmt.Errorf("failed to create views.py: %v", err)
	}

	// Update urls.py
	urlsContent := fmt.Sprintf(`from django.contrib import admin
from django.urls import path, include
from . import views

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),
    path('admin/', admin.site.urls),
    path('api-docs/', views.HomeView.as_view(template_name='api-docs.html'), name='api_docs'),
    path('__reload__/', include('django_browser_reload.urls')),
]`)

	urlsPath := filepath.Join(projectPath, m.projectName, "urls.py")
	if err := os.WriteFile(urlsPath, []byte(urlsContent), 0644); err != nil {
		return fmt.Errorf("failed to create urls.py: %v", err)
	}

	return nil
}
