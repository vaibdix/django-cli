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
    <html lang="en" class="dark">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>{% block title %}{{ project_name|default:"Django Site" }}{% endblock %}</title>
        <link rel="stylesheet" href="{% static 'css/style.css' %}">
        <script>
            tailwind.config = {
                darkMode: 'class',
                theme: {
                    extend: {
                        colors: {
                            'gray-950': '#0a0a0a',
                            'gray-925': '#111111',
                            'gray-900': '#171717',
                            'gray-850': '#1f1f1f',
                        },
                        fontFamily: {
                            'geist': ['-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', 'sans-serif'],
                            'geist-mono': ['Menlo', 'Monaco', 'Lucida Console', 'Liberation Mono', 'DejaVu Sans Mono', 'Bitstream Vera Sans Mono', 'Courier New', 'monospace'],
                        }
                    }
                }
            }
        </script>
        {% block extra_head %}{% endblock %}
    </head>
    <body class="bg-black text-white font-geist antialiased">
        <!-- Header -->
        <header class="sticky top-0 z-50 backdrop-blur-xl bg-black/80 border-b border-gray-800">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div class="flex justify-between items-center py-4">
                    <div class="flex items-center space-x-3">
                            <a href="http://localhost:8000" class="flex items-center space-x-2">
                            <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <rect width="32" height="32" rx="8" fill="white"/>
                                <path d="M12 8L20 16L12 24" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                            <h1 class="text-xl font-semibold">{{ project_name|default:"Django" }}</h1>
                        </div>
                    </div>

                    <nav class="hidden md:flex items-center space-x-8">
                        <a href="/" class="text-gray-300 hover:text-white transition-colors duration-200 text-sm">Home</a>
                        <a href="{% url 'api_docs' %}" class="text-gray-300 hover:text-white transition-colors duration-200 text-sm">Docs</a>
                        <a href="/admin/" class="text-gray-300 hover:text-white transition-colors duration-200 text-sm">Admin</a>
                        <a href="/api/v1/" class="bg-white text-black px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-200 transition-colors duration-200">
                            API
                        </a>
                    </nav>

                    <!-- Mobile menu button -->
                    <button class="md:hidden p-2 rounded-md hover:bg-gray-800 transition-colors">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                        </svg>
                    </button>
                </div>
            </div>
        </header>

        <!-- Main Content -->
        <main class="flex-1">
            {% block content %}{% endblock %}
        </main>

        <!-- Footer -->
        <footer class="border-t border-gray-800 mt-20">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
                <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
                    <div class="col-span-1 md:col-span-2">
                        <div class="flex items-center space-x-2 mb-6">
                            <svg width="24" height="24" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <rect width="32" height="32" rx="8" fill="white"/>
                                <path d="M12 8L20 16L12 24" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                            </svg>
                            <span class="text-lg font-semibold">{{ project_name|default:"Django Site" }}</span>
                        </div>
                        <p class="text-gray-400 mb-6 max-w-md">The Django framework that gives you everything you need to build full-stack web applications.</p>
                    </div>

                    <div>
                        <h3 class="text-sm font-semibold text-white mb-4">Resources</h3>
                        <ul class="space-y-3">
                            <li><a href="{% url 'api_docs' %}" class="text-gray-400 hover:text-white transition-colors text-sm">Documentation</a></li>
                            <li><a href="/api/v1/" class="text-gray-400 hover:text-white transition-colors text-sm">API Reference</a></li>
                            <li><a href="/admin/" class="text-gray-400 hover:text-white transition-colors text-sm">Admin Panel</a></li>
                        </ul>
                    </div>

                    <div>
                        <h3 class="text-sm font-semibold text-white mb-4">Support</h3>
                        <ul class="space-y-3">
                            <li><a href="#" class="text-gray-400 hover:text-white transition-colors text-sm">Help Center</a></li>
                            <li><a href="#" class="text-gray-400 hover:text-white transition-colors text-sm">Contact</a></li>
                            <li><a href="#" class="text-gray-400 hover:text-white transition-colors text-sm">Status</a></li>
                        </ul>
                    </div>
                </div>

                <div class="border-t border-gray-800 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
                    <p class="text-gray-400 text-sm">© 2025 {{ project_name|default:"Django Site" }}. All rights reserved.</p>
                    <div class="flex space-x-6 mt-4 md:mt-0">
                        <a href="#" class="text-gray-400 hover:text-white text-sm transition-colors">Privacy</a>
                        <a href="#" class="text-gray-400 hover:text-white text-sm transition-colors">Terms</a>
                    </div>
                </div>
            </div>
        </footer>

        <script src="{% static 'js/main.js' %}"></script>
        {{ django_browser_reload_script }}
    </body>
    </html>`
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "base.html"), []byte(baseContent), 0644); err != nil {
		return fmt.Errorf("failed to create base.html: %v", err)
	}

	indexContent := `{% extends 'base.html' %}
{% block title %}{{ project_name|default:"Django" }} - The Django Framework{% endblock %}

{% block content %}
<div class="relative">
    <!-- Hero Section -->
    <div class="relative overflow-hidden">
        <!-- Background gradient -->
        <div class="absolute inset-0 bg-gradient-to-b from-transparent via-black to-black pointer-events-none"></div>

        <!-- Grid background -->
        <div class="absolute inset-0 opacity-20">
            <div class="h-full w-full" style="background-image: radial-gradient(rgba(255,255,255,0.1) 1px, transparent 1px); background-size: 40px 40px;"></div>
        </div>

        <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-32">
            <div class="text-center">
                <!-- Badge -->
                <div class="inline-flex items-center rounded-full border border-gray-800 bg-gray-900/50 backdrop-blur-sm px-4 py-2 text-sm mb-8">
                    <span class="text-gray-300">🚀 Production ready Django application</span>
                </div>

                <!-- Main heading -->
                <h1 class="text-5xl md:text-7xl lg:text-8xl font-bold tracking-tight mb-8">
                    <span class="block">The Django</span>
                    <span class="block bg-gradient-to-r from-blue-400 via-purple-400 to-pink-400 bg-clip-text text-transparent">
                        Framework
                    </span>
                </h1>

                <!-- Subtitle -->
                <p class="text-xl md:text-2xl text-gray-400 max-w-3xl mx-auto mb-12 leading-relaxed">
                    Django provides everything you need to build fast, secure, and scalable web applications.
                    <span class="text-white">Used by thousands of developers worldwide.</span>
                </p>

                <!-- CTA Buttons -->
                <div class="flex flex-col sm:flex-row gap-4 justify-center mb-16">
                    <a href="{% url 'api_docs' %}" class="bg-white text-black px-8 py-4 rounded-md font-semibold hover:bg-gray-200 transition-colors duration-200 inline-flex items-center justify-center">
                        Get Started
                        <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                        </svg>
                    </a>
                    <a href="/api/v1/" class="border border-gray-700 text-white px-8 py-4 rounded-md font-semibold hover:border-gray-600 hover:bg-gray-900 transition-colors duration-200 inline-flex items-center justify-center">
                        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        Try API
                    </a>
                </div>

                <!-- Code example -->
                <div class="max-w-2xl mx-auto">
                    <div class="bg-gray-925 border border-gray-800 rounded-lg p-6 text-left">
                        <div class="flex items-center justify-between mb-4">
                            <div class="flex space-x-2">
                                <div class="w-3 h-3 rounded-full bg-red-500"></div>
                                <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
                                <div class="w-3 h-3 rounded-full bg-green-500"></div>
                            </div>
                            <span class="text-gray-400 text-sm">Django Project</span>
                        </div>
                        <pre class="text-sm text-gray-300 font-geist-mono"><code><span class="text-purple-400">from</span> <span class="text-blue-400">django.http</span> <span class="text-purple-400">import</span> <span class="text-yellow-400">JsonResponse</span>

<span class="text-purple-400">def</span> <span class="text-blue-400">api_view</span>(<span class="text-orange-400">request</span>):
    <span class="text-purple-400">return</span> <span class="text-yellow-400">JsonResponse</span>({
        <span class="text-green-400">'message'</span>: <span class="text-green-400">'Hello, Django!'</span>,
        <span class="text-green-400">'status'</span>: <span class="text-green-400">'success'</span>
    })</code></pre>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Features Section -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div class="text-center mb-16">
            <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">Why Django?</h2>
            <p class="text-xl text-gray-400 max-w-2xl mx-auto">
                Built for speed, security, and scalability. Trusted by startups and enterprises.
            </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <!-- Feature 1 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-blue-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">Fast Development</h3>
                <p class="text-gray-400">Django's batteries-included approach means you can build full-featured applications quickly without reinventing the wheel.</p>
            </div>

            <!-- Feature 2 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-green-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">Security First</h3>
                <p class="text-gray-400">Built-in protection against common security threats like SQL injection, CSRF, and XSS attacks.</p>
            </div>

            <!-- Feature 3 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-purple-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">Scalable</h3>
                <p class="text-gray-400">From small projects to high-traffic applications, Django scales with your needs and handles millions of users.</p>
            </div>

            <!-- Feature 4 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-orange-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">Rich Ecosystem</h3>
                <p class="text-gray-400">Thousands of packages and a vibrant community provide solutions for almost any use case.</p>
            </div>

            <!-- Feature 5 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-pink-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-pink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">Admin Interface</h3>
                <p class="text-gray-400">Automatic admin interface for content management, user authentication, and database operations.</p>
            </div>

            <!-- Feature 6 -->
            <div class="border border-gray-800 rounded-lg p-8 bg-gray-925 hover:border-gray-700 transition-colors duration-200">
                <div class="w-12 h-12 bg-cyan-500/10 rounded-lg flex items-center justify-center mb-6">
                    <svg class="w-6 h-6 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                </div>
                <h3 class="text-xl font-semibold text-white mb-3">REST API</h3>
                <p class="text-gray-400">Built-in support for creating powerful REST APIs with authentication, serialization, and documentation.</p>
            </div>
        </div>
    </div>

    <!-- Stats Section -->
    <div class="border-t border-gray-800">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
            <div class="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
                <div>
                    <div class="text-4xl font-bold text-white mb-2">15+</div>
                    <div class="text-gray-400 text-sm">Years of Development</div>
                </div>
                <div>
                    <div class="text-4xl font-bold text-white mb-2">1M+</div>
                    <div class="text-gray-400 text-sm">Websites Built</div>
                </div>
                <div>
                    <div class="text-4xl font-bold text-white mb-2">99.9%</div>
                    <div class="text-gray-400 text-sm">Uptime</div>
                </div>
                <div>
                    <div class="text-4xl font-bold text-white mb-2">24/7</div>
                    <div class="text-gray-400 text-sm">Community Support</div>
                </div>
            </div>
        </div>
    </div>

    <!-- CTA Section -->
    <div class="border-t border-gray-800">
        <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
            <h2 class="text-3xl md:text-4xl font-bold text-white mb-6">
                Start building today
            </h2>
            <p class="text-xl text-gray-400 mb-12 max-w-2xl mx-auto">
                Join thousands of developers who trust Django to build their next big project.
            </p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
                <a href="{% url 'api_docs' %}" class="bg-white text-black px-8 py-4 rounded-md font-semibold hover:bg-gray-200 transition-colors duration-200 inline-flex items-center justify-center">
                    Get Started
                    <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                </a>
                <a href="/admin/" class="border border-gray-700 text-white px-8 py-4 rounded-md font-semibold hover:border-gray-600 hover:bg-gray-900 transition-colors duration-200">
                    Admin Panel
                </a>
            </div>
        </div>
    </div>
</div>
{% endblock %}`
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "index.html"), []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.html: %v", err)
	}

	styleContent := `
    /* Custom styles for enhanced animations and effects */
    @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700;800;900&display=swap');

    body {
        font-family: 'Inter', sans-serif;
        scroll-behavior: smooth;
    }

    `
	if err := os.WriteFile(filepath.Join(staticPath, "css", "style.css"), []byte(styleContent), 0644); err != nil {
		return fmt.Errorf("failed to create style.css: %v", err)
	}

	jsContent := `
    // Enhanced JavaScript for better interactivity
    console.log('Django project initialized with modern design!');

    `
	if err := os.WriteFile(filepath.Join(staticPath, "js", "main.js"), []byte(jsContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.js: %v", err)
	}

	// Write the API docs template
	apiDocsContent := `{% extends 'base.html' %}
{% block title %}API Documentation - {{ project_name }}{% endblock %}

{% block content %}
<div class="bg-black text-white min-h-screen">
    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <!-- Hero Section -->
        <div class="text-center mb-20">
            <div class="inline-flex items-center bg-gray-900 border border-gray-800 rounded-full px-4 py-2 mb-8">
                <span class="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                <span class="text-sm text-gray-300">API Documentation</span>
            </div>

            <h1 class="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white via-gray-300 to-gray-500 bg-clip-text text-transparent">
                Books API
            </h1>
            <p class="text-xl text-gray-400 max-w-2xl mx-auto leading-relaxed">
                A powerful REST API for managing your book collection with full CRUD operations, authentication, and more.
            </p>
        </div>

        <!-- Quick Stats -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-20">
            <div class="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                <div class="text-2xl font-bold text-white mb-1">12+</div>
                <div class="text-sm text-gray-400">Endpoints</div>
            </div>
            <div class="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                <div class="text-2xl font-bold text-white mb-1">REST</div>
                <div class="text-sm text-gray-400">Architecture</div>
            </div>
            <div class="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                <div class="text-2xl font-bold text-white mb-1">JSON</div>
                <div class="text-sm text-gray-400">Response</div>
            </div>
            <div class="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                <div class="text-2xl font-bold text-white mb-1">Auth</div>
                <div class="text-sm text-gray-400">Secured</div>
            </div>
        </div>

        <!-- API Endpoints Section -->
        <div class="space-y-12">
            <!-- Books API -->
            <section>
                <div class="mb-8">
                    <h2 class="text-3xl font-bold text-white mb-3">Books API</h2>
                    <p class="text-gray-400">Manage your book collection with full CRUD operations</p>
                </div>

                <div class="space-y-6">
                    <!-- GET & POST /api/v1/books/ -->
                    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors">
                        <div class="p-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-2 py-1 text-xs font-mono bg-green-500/20 text-green-400 border border-green-500/30 rounded">GET</span>
                                    <span class="px-2 py-1 text-xs font-mono bg-blue-500/20 text-blue-400 border border-blue-500/30 rounded">POST</span>
                                    <code class="text-sm font-mono text-gray-300 bg-gray-800 px-3 py-1 rounded">/api/v1/books/</code>
                                </div>
                            </div>
                            <p class="text-gray-400 mb-4">Retrieve all books or create a new book entry with title, author, and publication details.</p>
                            <div class="bg-black/50 border border-gray-800 rounded-lg p-4">
                                <pre class="text-sm text-gray-300 font-mono">
<span class="text-green-400">GET</span>: Returns paginated list of books
<span class="text-blue-400">POST</span>: Creates new book (requires: title, author, isbn)</pre>
                            </div>
                        </div>
                    </div>

                    <!-- GET, PUT, DELETE /api/v1/books/{id}/ -->
                    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors">
                        <div class="p-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-2 py-1 text-xs font-mono bg-green-500/20 text-green-400 border border-green-500/30 rounded">GET</span>
                                    <span class="px-2 py-1 text-xs font-mono bg-yellow-500/20 text-yellow-400 border border-yellow-500/30 rounded">PUT</span>
                                    <span class="px-2 py-1 text-xs font-mono bg-red-500/20 text-red-400 border border-red-500/30 rounded">DELETE</span>
                                    <code class="text-sm font-mono text-gray-300 bg-gray-800 px-3 py-1 rounded">/api/v1/books/{id}/</code>
                                </div>
                            </div>
                            <p class="text-gray-400 mb-4">Retrieve, update, or delete a specific book by its unique identifier.</p>
                            <div class="bg-black/50 border border-gray-800 rounded-lg p-4">
                                <pre class="text-sm text-gray-300 font-mono">
<span class="text-green-400">GET</span>: Returns book details
<span class="text-yellow-400">PUT</span>: Updates book (partial updates supported)
<span class="text-red-400">DELETE</span>: Removes book from collection</pre>
                            </div>
                        </div>
                    </div>

                    <!-- GET /api/v1/books/recent/ -->
                    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors">
                        <div class="p-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-2 py-1 text-xs font-mono bg-green-500/20 text-green-400 border border-green-500/30 rounded">GET</span>
                                    <code class="text-sm font-mono text-gray-300 bg-gray-800 px-3 py-1 rounded">/api/v1/books/recent/</code>
                                </div>
                            </div>
                            <p class="text-gray-400 mb-4">Get the most recently added books, sorted by creation date.</p>
                            <div class="bg-black/50 border border-gray-800 rounded-lg p-4">
                                <pre class="text-sm text-gray-300 font-mono">Returns: Latest 10 books by default (configurable with ?limit parameter)</pre>
                            </div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Authentication Section -->
            <section>
                <div class="mb-8">
                    <h2 class="text-3xl font-bold text-white mb-3">Authentication</h2>
                    <p class="text-gray-400">Secure access to protected endpoints</p>
                </div>

                <div class="space-y-6">
                    <!-- Login Endpoint -->
                    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors">
                        <div class="p-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-2 py-1 text-xs font-mono bg-green-500/20 text-green-400 border border-green-500/30 rounded">GET</span>
                                    <span class="px-2 py-1 text-xs font-mono bg-blue-500/20 text-blue-400 border border-blue-500/30 rounded">POST</span>
                                    <code class="text-sm font-mono text-gray-300 bg-gray-800 px-3 py-1 rounded">/api-auth/login/</code>
                                </div>
                            </div>
                            <p class="text-gray-400">Authenticate users and obtain session credentials for API access.</p>
                        </div>
                    </div>

                    <!-- Logout Endpoint -->
                    <div class="bg-gray-900/30 border border-gray-800 rounded-xl overflow-hidden hover:border-gray-700 transition-colors">
                        <div class="p-6">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-2 py-1 text-xs font-mono bg-green-500/20 text-green-400 border border-green-500/30 rounded">GET</span>
                                    <span class="px-2 py-1 text-xs font-mono bg-blue-500/20 text-blue-400 border border-blue-500/30 rounded">POST</span>
                                    <code class="text-sm font-mono text-gray-300 bg-gray-800 px-3 py-1 rounded">/api-auth/logout/</code>
                                </div>
                            </div>
                            <p class="text-gray-400">Safely terminate user sessions and invalidate authentication credentials.</p>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Quick Start Guide -->
            <section class="bg-gradient-to-r from-gray-900 to-gray-800 border border-gray-700 rounded-xl p-8">
                <h2 class="text-2xl font-bold text-white mb-6">Quick Start Guide</h2>

                <div class="grid md:grid-cols-2 gap-6">
                    <div>
                        <h3 class="text-lg font-semibold text-white mb-3">1. Authentication</h3>
                        <div class="bg-black/50 border border-gray-700 rounded-lg p-4">
                            <pre class="text-sm text-gray-300 font-mono overflow-x-auto">
<span class="text-purple-400">curl</span> <span class="text-blue-400">-X POST</span> http://localhost:8000/api-auth/login/ \
  <span class="text-blue-400">-d</span> <span class="text-green-400">"username=your_username&password=your_password"</span></pre>
                        </div>
                    </div>

                    <div>
                        <h3 class="text-lg font-semibold text-white mb-3">2. Fetch Books</h3>
                        <div class="bg-black/50 border border-gray-700 rounded-lg p-4">
                            <pre class="text-sm text-gray-300 font-mono overflow-x-auto">
<span class="text-purple-400">curl</span> <span class="text-blue-400">-X GET</span> http://localhost:8000/api/v1/books/ \
  <span class="text-blue-400">-H</span> <span class="text-green-400">"Authorization: Bearer your_token"</span></pre>
                        </div>
                    </div>
                </div>
            </section>
        </div>

        <!-- Action Buttons -->
        <div class="text-center mt-20">
            <div class="flex flex-col sm:flex-row gap-4 justify-center">
                <a href="/api/v1/" class="bg-white text-black px-8 py-3 rounded-lg font-semibold hover:bg-gray-200 transition-colors inline-flex items-center justify-center">
                    Try API Now
                    <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                    </svg>
                </a>
                <a href="/" class="bg-transparent border border-gray-700 text-white px-8 py-3 rounded-lg font-semibold hover:border-gray-600 transition-colors inline-flex items-center justify-center">
                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                    </svg>
                    Back to Home
                </a>
            </div>
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
