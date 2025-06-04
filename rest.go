package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (m *Model) setupDjangoRestFramework(projectPath string) error {
	if !m.setupRestFramework {
		return nil
	}

	m.updateProgress("Setting up Django REST Framework...")

	// Install Django REST Framework
	pythonCmd := getPythonPath(projectPath)
	cmd := exec.Command(pythonCmd, "-m", "pip", "install", "djangorestframework")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.stepMessages = append(m.stepMessages, fmt.Sprintf("⚠️  Warning: Failed to install Django REST Framework: %v\nOutput: %s", err, string(output)))
		return err
	}
	m.stepMessages = append(m.stepMessages, "✅ Django REST Framework installed.")

	// Update settings.py
	settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings.py: %v", err)
	}

	// Add 'rest_framework' to INSTALLED_APPS
	settingsStr := string(settingsContent)
	if !strings.Contains(settingsStr, "'rest_framework',") {
		settingsStr = strings.Replace(
			settingsStr,
			"INSTALLED_APPS = [",
			"INSTALLED_APPS = [\n    'rest_framework',",
			1,
		)

		// Add REST_FRAMEWORK settings
		restSettings := `
# Django REST Framework settings
REST_FRAMEWORK = {
    'DEFAULT_PERMISSION_CLASSES': [
        'rest_framework.permissions.AllowAny',
    ],
    'DEFAULT_RENDERER_CLASSES': [
        'rest_framework.renderers.JSONRenderer',
        'rest_framework.renderers.BrowsableAPIRenderer',
    ],
    'DEFAULT_PAGINATION_CLASS': 'rest_framework.pagination.PageNumberPagination',
    'PAGE_SIZE': 20
}
`
		settingsStr += restSettings

		if err := os.WriteFile(settingsPath, []byte(settingsStr), 0644); err != nil {
			return fmt.Errorf("failed to update settings.py: %v", err)
		}
		m.stepMessages = append(m.stepMessages, "✅ Added REST Framework to INSTALLED_APPS and configured settings.")
	}

	// Update urls.py to include REST Framework URLs
	urlsPath := filepath.Join(projectPath, m.projectName, "urls.py")
	urlsContent, err := os.ReadFile(urlsPath)
	if err != nil {
		return fmt.Errorf("failed to read urls.py: %v", err)
	}

	urlsStr := string(urlsContent)
	if !strings.Contains(urlsStr, "rest_framework.urls") {
		// First, ensure we have the include import
		urlsStr = strings.Replace(
			urlsStr,
			"from django.urls import path",
			"from django.urls import path, include",
			1,
		)

		// Create a new api.py file in the project's configuration directory
		apiUrlsContent := fmt.Sprintf(`from django.urls import path, include
from rest_framework.routers import DefaultRouter
from %s.views import BookViewSet

router = DefaultRouter()
router.register(r'books', BookViewSet, basename='book')

app_name = 'api'

urlpatterns = [
    path('', include(router.urls)),
]`, m.appName)

		// Write the api.py file
		projectConfigDir := filepath.Join(projectPath, m.projectName)
		if err := os.WriteFile(filepath.Join(projectConfigDir, "api.py"), []byte(apiUrlsContent), 0644); err != nil {
			return fmt.Errorf("failed to create api.py: %v", err)
		}

		// Update main urls.py to include the API docs view
		projectUrlsContent := fmt.Sprintf(`from django.contrib import admin
from django.urls import path, include
from django.views.generic import TemplateView

urlpatterns = [
    # Root URL - points to our main index.html
    path('', TemplateView.as_view(template_name='index.html'), name='home'),
    path('admin/', admin.site.urls),
    path('__reload__/', include('django_browser_reload.urls')),
    path('api-docs/', TemplateView.as_view(template_name='api-docs.html'), name='api_docs'),
    path('api/v1/', include('%s.api')),
    path('api-auth/', include('rest_framework.urls', namespace='rest_framework')),
`, m.projectName)

		if m.appName != "" {
			projectUrlsContent += fmt.Sprintf("    path('%s/', include('%s.urls')),\n", m.appName, m.appName)
		}

		projectUrlsContent += "]"

		// Write the complete urls.py file
		if err := os.WriteFile(urlsPath, []byte(projectUrlsContent), 0644); err != nil {
			return fmt.Errorf("failed to update urls.py: %v", err)
		}
		m.stepMessages = append(m.stepMessages, "✅ Added REST Framework URLs and API endpoints.")
	}

	// Create example API if an app is created
	if m.appName != "" {
		if err := m.createExampleAPI(projectPath); err != nil {
			return err
		}
	}

	// Update the app's urls.py to separate API endpoints
	appUrlsContent := fmt.Sprintf(`from django.urls import path, include
from . import views

app_name = '%s'

urlpatterns = [
    # Regular app views only
    path('', views.index, name='index'),
]
`, m.appName)

	// Create a new api.py file in the project's configuration directory
	apiUrlsContent := fmt.Sprintf(`from django.urls import path, include
from rest_framework.routers import DefaultRouter
from %s.views import BookViewSet

router = DefaultRouter()
router.register(r'books', BookViewSet, basename='book')

urlpatterns = [
    path('', include(router.urls)),
    path('auth/', include('rest_framework.urls')),
]
`, m.appName)

	// Write the api.py file
	projectConfigDir := filepath.Join(projectPath, m.projectName)
	if err := os.WriteFile(filepath.Join(projectConfigDir, "api.py"), []byte(apiUrlsContent), 0644); err != nil {
		return fmt.Errorf("failed to create api.py: %v", err)
	}

	// Update the main urls.py to include API URLs at root level
	urlsStr = strings.Replace(
		urlsStr,
		"urlpatterns = [",
		fmt.Sprintf(`urlpatterns = [
    # API endpoints
    path('api/v1/', include('%s.api')),`, m.projectName),
		1,
	)

	// Write the app's urls.py
	if err := os.WriteFile(filepath.Join(projectPath, m.appName, "urls.py"), []byte(appUrlsContent), 0644); err != nil {
		return fmt.Errorf("failed to update app urls.py: %v", err)
	}

	return nil
}

func (m *Model) createExampleAPI(projectPath string) error {
	// Move appUrlsContent declaration to the beginning of createExampleAPI
	appUrlsContent := fmt.Sprintf(`from django.urls import path, include
from . import views

app_name = '%s'

urlpatterns = [
    # Regular app views only
    path('', views.index, name='index'),
]
`, m.appName)

	// Create serializers.py
	serializersPath := filepath.Join(projectPath, m.appName, "serializers.py")
	serializersContent := fmt.Sprintf(`from rest_framework import serializers
from .models import Book

class BookSerializer(serializers.ModelSerializer):
    class Meta:
        model = Book
        fields = '__all__'
        read_only_fields = ('created_at', 'updated_at')
`)

	if err := os.WriteFile(serializersPath, []byte(serializersContent), 0644); err != nil {
		return fmt.Errorf("failed to create serializers.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Created serializers.py with example BookSerializer.")

	// Update models.py with example model
	modelsPath := filepath.Join(projectPath, m.appName, "models.py")
	modelsContent := `from django.db import models

class Book(models.Model):
    title = models.CharField(max_length=200)
    author = models.CharField(max_length=100)
    isbn = models.CharField(max_length=13, unique=True)
    publication_date = models.DateField()
    price = models.DecimalField(max_digits=10, decimal_places=2)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.title
`

	if err := os.WriteFile(modelsPath, []byte(modelsContent), 0644); err != nil {
		return fmt.Errorf("failed to update models.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Created example Book model.")

	// Update views.py with both regular view and ViewSet
	viewsPath := filepath.Join(projectPath, m.appName, "views.py")
	viewsContent := fmt.Sprintf(`from django.shortcuts import render
from rest_framework import viewsets
from rest_framework.decorators import action
from rest_framework.response import Response
from django.utils import timezone
from datetime import timedelta
from .models import Book
from .serializers import BookSerializer

def index(request):
    return render(request, '%s/index.html')

class BookViewSet(viewsets.ModelViewSet):
    queryset = Book.objects.all()
    serializer_class = BookSerializer

    @action(detail=False, methods=['get'])
    def recent(self, request):
        recent_books = Book.objects.filter(
            created_at__gte=timezone.now() - timedelta(days=30)
        )
        serializer = self.get_serializer(recent_books, many=True)
        return Response(serializer.data)
`, m.appName)

	if err := os.WriteFile(viewsPath, []byte(viewsContent), 0644); err != nil {
		return fmt.Errorf("failed to update views.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Created BookViewSet with custom action.")

	// Create or update app's urls.py
	appUrlsPath := filepath.Join(projectPath, m.appName, "urls.py")
	if err := os.WriteFile(appUrlsPath, []byte(appUrlsContent), 0644); err != nil {
		return fmt.Errorf("failed to create/update app urls.py: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Configured API URLs with DefaultRouter.")

	// Create management command for sample data
	managementDir := filepath.Join(projectPath, m.appName, "management", "commands")
	if err := os.MkdirAll(managementDir, 0755); err != nil {
		return fmt.Errorf("failed to create management directory: %v", err)
	}

	sampleDataPath := filepath.Join(managementDir, "create_sample_data.py")
	sampleDataContent := fmt.Sprintf(`from django.core.management.base import BaseCommand
from %s.models import Book
from datetime import date

class Command(BaseCommand):
    help = 'Creates sample book data'

    def handle(self, *args, **options):
        books = [
            {
                'title': 'Django for Beginners',
                'author': 'William Vincent',
                'isbn': '9781735467200',
                'publication_date': date(2022, 1, 1),
                'price': 39.99
            },
            {
                'title': 'Two Scoops of Django',
                'author': 'Daniel Roy Greenfeld',
                'isbn': '9780692915738',
                'publication_date': date(2021, 5, 15),
                'price': 49.99
            }
        ]
        
        for book_data in books:
            Book.objects.get_or_create(**book_data)
        
        self.stdout.write(self.style.SUCCESS('Sample data created successfully!'))
`, m.appName)

	if err := os.WriteFile(sampleDataPath, []byte(sampleDataContent), 0644); err != nil {
		return fmt.Errorf("failed to create sample data command: %v", err)
	}
	m.stepMessages = append(m.stepMessages, "✅ Created management command for sample data.")

	// Create __init__.py files for management commands
	initPath := filepath.Join(projectPath, m.appName, "management", "__init__.py")
	if err := os.WriteFile(initPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create management __init__.py: %v", err)
	}

	initCommandsPath := filepath.Join(projectPath, m.appName, "management", "commands", "__init__.py")
	if err := os.WriteFile(initCommandsPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create commands __init__.py: %v", err)
	}

	return nil
}
