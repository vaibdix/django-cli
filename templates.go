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
    <body class="bg-gradient-to-br from-slate-50 via-blue-50 to-indigo-100 min-h-screen">
        <!-- Header -->
        <header class="bg-white/80 backdrop-blur-lg border-b border-white/20 sticky top-0 z-50">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div class="flex justify-between items-center py-4">
                    <div class="flex items-center space-x-3">
                        <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center">
                            <span class="text-white font-bold text-lg">D</span>
                        </div>
                        <div>
                            <h1 class="text-xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                                {{ project_name|default:"Django Site" }}
                            </h1>
                            <p class="text-sm text-gray-500">Professional Development</p>
                        </div>
                    </div>
                    
                    <nav class="hidden md:flex items-center space-x-8">
                        <a href="/" class="text-gray-700 hover:text-indigo-600 font-medium transition-colors duration-200">Home</a>
                        <a href="{% url 'api_docs' %}" class="text-gray-700 hover:text-indigo-600 font-medium transition-colors duration-200">Docs</a>
                        <a href="/admin/" class="text-gray-700 hover:text-indigo-600 font-medium transition-colors duration-200">Admin</a>
                        <a href="/api/v1/" class="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-4 py-2 rounded-lg font-medium hover:shadow-lg transform hover:scale-105 transition-all duration-200">
                            API
                        </a>
                    </nav>
    
                    <!-- Mobile menu button -->
                    <button class="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors">
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
        <footer class="bg-white/60 backdrop-blur-lg border-t border-white/20 mt-20">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
                <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
                    <div class="col-span-1 md:col-span-2">
                        <div class="flex items-center space-x-3 mb-4">
                            <div class="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                                <span class="text-white font-bold">D</span>
                            </div>
                            <span class="text-xl font-bold text-gray-900">{{ project_name|default:"Django Site" }}</span>
                        </div>
                        <p class="text-gray-600 mb-4">Building powerful web applications with Django. Professional, modern, and scalable solutions.</p>
                        <div class="flex space-x-4">
                            <a href="#" class="w-10 h-10 bg-gray-200 rounded-lg flex items-center justify-center hover:bg-indigo-500 hover:text-white transition-all duration-200">
                                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"/></svg>
                            </a>
                            <a href="#" class="w-10 h-10 bg-gray-200 rounded-lg flex items-center justify-center hover:bg-indigo-500 hover:text-white transition-all duration-200">
                                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M22.46 6c-.77.35-1.6.58-2.46.69.88-.53 1.56-1.37 1.88-2.38-.83.5-1.75.85-2.72 1.05C18.37 4.5 17.26 4 16 4c-2.35 0-4.27 1.92-4.27 4.29 0 .34.04.67.11.98C8.28 9.09 5.11 7.38 3 4.79c-.37.63-.58 1.37-.58 2.15 0 1.49.75 2.81 1.91 3.56-.71 0-1.37-.2-1.95-.5v.03c0 2.08 1.48 3.82 3.44 4.21a4.22 4.22 0 0 1-1.93.07 4.28 4.28 0 0 0 4 2.98 8.521 8.521 0 0 1-5.33 1.84c-.34 0-.68-.02-1.02-.06C3.44 20.29 5.7 21 8.12 21 16 21 20.33 14.46 20.33 8.79c0-.19 0-.37-.01-.56.84-.6 1.56-1.36 2.14-2.23z"/></svg>
                            </a>
                            <a href="#" class="w-10 h-10 bg-gray-200 rounded-lg flex items-center justify-center hover:bg-indigo-500 hover:text-white transition-all duration-200">
                                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
                            </a>
                        </div>
                    </div>
                    
                    <div>
                        <h3 class="text-sm font-semibold text-gray-900 tracking-wider uppercase mb-4">Resources</h3>
                        <ul class="space-y-3">
                            <li><a href="{% url 'api_docs' %}" class="text-gray-600 hover:text-indigo-600 transition-colors">Documentation</a></li>
                            <li><a href="/api/v1/" class="text-gray-600 hover:text-indigo-600 transition-colors">API Reference</a></li>
                            <li><a href="/admin/" class="text-gray-600 hover:text-indigo-600 transition-colors">Admin Panel</a></li>
                        </ul>
                    </div>
                    
                    <div>
                        <h3 class="text-sm font-semibold text-gray-900 tracking-wider uppercase mb-4">Support</h3>
                        <ul class="space-y-3">
                            <li><a href="#" class="text-gray-600 hover:text-indigo-600 transition-colors">Help Center</a></li>
                            <li><a href="#" class="text-gray-600 hover:text-indigo-600 transition-colors">Contact Us</a></li>
                            <li><a href="#" class="text-gray-600 hover:text-indigo-600 transition-colors">Status</a></li>
                        </ul>
                    </div>
                </div>
                
                <div class="border-t border-gray-200 mt-8 pt-8 flex flex-col md:flex-row justify-between items-center">
                    <p class="text-gray-500 text-sm">© 2025 {{ project_name|default:"Django Site" }}. All rights reserved.</p>
                    <div class="flex space-x-6 mt-4 md:mt-0">
                        <a href="#" class="text-gray-500 hover:text-indigo-600 text-sm transition-colors">Privacy Policy</a>
                        <a href="#" class="text-gray-500 hover:text-indigo-600 text-sm transition-colors">Terms of Service</a>
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
    {% block title %}Welcome - {{ project_name }}{% endblock %}
    
    {% block content %}
    <div class="min-h-screen">
        <!-- Hero Section -->
        <div class="relative overflow-hidden">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
                <div class="text-center animate-fade-in">
                    <div class="inline-flex items-center bg-gradient-to-r from-indigo-500/10 to-purple-500/10 rounded-full px-6 py-2 mb-8">
                        <span class="text-sm font-medium bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent">
                            ✨ Some make promises... We have proof
                        </span>
                    </div>
                    
                    <h1 class="text-5xl md:text-7xl font-bold mb-6">
                        <span class="bg-gradient-to-r from-gray-900 via-indigo-900 to-purple-900 bg-clip-text text-transparent">
                            Welcome to
                        </span>
                        <br>
                        <span class="bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 bg-clip-text text-transparent animate-bounce-gentle">
                            {{ project_name|default:"Django Site" }}
                        </span>
                    </h1>
                    
                    <p class="text-xl text-gray-600 max-w-3xl mx-auto mb-12 leading-relaxed">
                        Your Django project is ready for development. Experience award-winning design with powerful functionality.
                    </p>
                    
                    <div class="flex flex-col sm:flex-row gap-4 justify-center">
                        <a href="{% url 'api_docs' %}" class="bg-gradient-to-r from-indigo-500 to-purple-600 text-white px-8 py-4 rounded-2xl font-semibold hover:shadow-2xl hover:shadow-indigo-500/25 transform hover:scale-105 transition-all duration-300 inline-flex items-center justify-center">
                            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                            Explore Documentation
                        </a>
                        <a href="/api/v1/" class="bg-white text-gray-900 border-2 border-gray-200 px-8 py-4 rounded-2xl font-semibold hover:border-indigo-300 hover:shadow-xl transform hover:scale-105 transition-all duration-300 inline-flex items-center justify-center">
                            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                            Try API
                        </a>
                    </div>
                </div>
            </div>
        </div>
    
        <!-- Feature Cards Section -->
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
            <div class="text-center mb-16 animate-slide-up">
                <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                    Award-Winning Features
                </h2>
                <p class="text-xl text-gray-600">Discover what makes our platform exceptional</p>
            </div>
    
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                <!-- Documentation Card -->
                <div class="group relative bg-gradient-to-br from-purple-400 to-indigo-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-purple-500/25 animate-slide-up" style="animation-delay: 0.1s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-sm font-bold">W.</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">Documentation</h3>
                        <p class="text-purple-100 text-sm mb-4">Site Of The Day</p>
                        <p class="text-white/90">
                            Comprehensive API docs, authentication guides, and admin documentation
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <a href="{% url 'api_docs' %}" class="bg-white/20 hover:bg-white/30 rounded-full p-3 transition-all duration-200 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                            </svg>
                        </a>
                    </div>
                </div>
    
                {% if app_name %}
                <!-- App Card -->
                <div class="group relative bg-gradient-to-br from-green-400 to-emerald-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-green-500/25 animate-slide-up" style="animation-delay: 0.2s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-lg">⚡</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">{{ app_name|title }} App</h3>
                        <p class="text-green-100 text-sm mb-4">Website Of The Day</p>
                        <p class="text-white/90">
                            Access your application's main functionality and features
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <a href="/{{ app_name }}/" class="bg-white/20 hover:bg-white/30 rounded-full p-3 transition-all duration-200 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                            </svg>
                        </a>
                    </div>
                </div>
                {% endif %}
    
                <!-- API Card -->
                <div class="group relative bg-gradient-to-br from-cyan-400 to-blue-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-cyan-500/25 animate-slide-up" style="animation-delay: 0.3s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-lg">🔥</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">RESTful API</h3>
                        <p class="text-cyan-100 text-sm mb-4">Developer Award</p>
                        <p class="text-white/90">
                            Powerful REST API endpoints with full CRUD operations
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <a href="/api/v1/" class="bg-white/20 hover:bg-white/30 rounded-full p-3 transition-all duration-200 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                            </svg>
                        </a>
                    </div>
                </div>
    
                <!-- Admin Panel Card -->
                <div class="group relative bg-gradient-to-br from-pink-400 to-rose-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-pink-500/25 animate-slide-up" style="animation-delay: 0.4s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-sm font-bold">W.</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">Admin Panel</h3>
                        <p class="text-pink-100 text-sm mb-4">Honorable Mention</p>
                        <p class="text-white/90">
                            Powerful Django admin interface for content management
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <a href="/admin/" class="bg-white/20 hover:bg-white/30 rounded-full p-3 transition-all duration-200 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                            </svg>
                        </a>
                    </div>
                </div>
    
                <!-- Authentication Card -->
                <div class="group relative bg-gradient-to-br from-amber-400 to-orange-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-amber-500/25 animate-slide-up" style="animation-delay: 0.5s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-lg">🛡️</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">Authentication</h3>
                        <p class="text-amber-100 text-sm mb-4">Security Award</p>
                        <p class="text-white/90">
                            Secure user authentication and session management
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <a href="/api-auth/login/" class="bg-white/20 hover:bg-white/30 rounded-full p-3 transition-all duration-200 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                            </svg>
                        </a>
                    </div>
                </div>
    
                <!-- Performance Card -->
                <div class="group relative bg-gradient-to-br from-violet-400 to-purple-600 rounded-3xl p-8 text-white transform hover:scale-105 transition-all duration-500 hover:shadow-2xl hover:shadow-violet-500/25 animate-slide-up" style="animation-delay: 0.6s">
                    <div class="absolute top-4 right-4">
                        <div class="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                            <span class="text-lg">⚡</span>
                        </div>
                    </div>
                    
                    <div class="mb-6">
                        <h3 class="text-2xl font-bold mb-2">Performance</h3>
                        <p class="text-violet-100 text-sm mb-4">Speed Champion</p>
                        <p class="text-white/90">
                            Optimized for speed with modern web technologies
                        </p>
                    </div>
                    
                    <div class="flex justify-between items-end">
                        <span class="text-3xl font-bold">2025</span>
                        <div class="bg-white/20 rounded-full p-3 backdrop-blur-sm">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    
        <!-- Stats Section -->
        <div class="bg-white/60 backdrop-blur-lg">
            <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
                <div class="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
                    <div class="animate-fade-in">
                        <div class="text-4xl font-bold bg-gradient-to-r from-indigo-600 to-purple-600 bg-clip-text text-transparent mb-2">99.9%</div>
                        <div class="text-gray-600 font-medium">Uptime</div>
                    </div>
                    <div class="animate-fade-in" style="animation-delay: 0.1s">
                        <div class="text-4xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent mb-2">< 100ms</div>
                        <div class="text-gray-600 font-medium">Response Time</div>
                    </div>
                    <div class="animate-fade-in" style="animation-delay: 0.2s">
                        <div class="text-4xl font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent mb-2">10k+</div>
                        <div class="text-gray-600 font-medium">API Calls/day</div>
                    </div>
                    <div class="animate-fade-in" style="animation-delay: 0.3s">
                        <div class="text-4xl font-bold bg-gradient-to-r from-pink-600 to-rose-600 bg-clip-text text-transparent mb-2">24/7</div>
                        <div class="text-gray-600 font-medium">Support</div>
                    </div>
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
    <div class="min-h-screen py-12">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <!-- Header Section -->
            <div class="text-center mb-16 animate-fade-in">
                <div class="inline-flex items-center bg-gradient-to-r from-blue-500/10 to-purple-500/10 rounded-full px-6 py-2 mb-6">
                    <span class="text-sm font-medium bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                        📚 Comprehensive Documentation
                    </span>
                </div>
                
                <h1 class="text-4xl md:text-5xl font-bold mb-6">
                    <span class="bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 bg-clip-text text-transparent">
                        API Documentation
                    </span>
                </h1>
                <p class="text-xl text-gray-600 max-w-3xl mx-auto">
                    Explore our powerful REST API endpoints, authentication methods, and integration guides
                </p>
            </div>
    
            <!-- Quick Stats -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-6 mb-16">
                <div class="bg-white/60 backdrop-blur-lg rounded-2xl p-6 text-center animate-slide-up">
                    <div class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-cyan-600 bg-clip-text text-transparent mb-2">12+</div>
                    <div class="text-gray-600 text-sm font-medium">Endpoints</div>
                </div>
                <div class="bg-white/60 backdrop-blur-lg rounded-2xl p-6 text-center animate-slide-up" style="animation-delay: 0.1s">
                    <div class="text-3xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent mb-2">REST</div>
                    <div class="text-gray-600 text-sm font-medium">Architecture</div>
                </div>
                <div class="bg-white/60 backdrop-blur-lg rounded-2xl p-6 text-center animate-slide-up" style="animation-delay: 0.2s">
                    <div class="text-3xl font-bold bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent mb-2">JSON</div>
                    <div class="text-gray-600 text-sm font-medium">Response</div>
                </div>
                <div class="bg-white/60 backdrop-blur-lg rounded-2xl p-6 text-center animate-slide-up" style="animation-delay: 0.3s">
                    <div class="text-3xl font-bold bg-gradient-to-r from-orange-600 to-red-600 bg-clip-text text-transparent mb-2">Auth</div>
                    <div class="text-gray-600 text-sm font-medium">Secured</div>
                </div>
            </div>
    
            <!-- API Endpoints Section -->
            <div class="space-y-8">
                <div class="bg-white/80 backdrop-blur-lg shadow-2xl rounded-3xl overflow-hidden animate-slide-up">
                    <div class="bg-gradient-to-r from-blue-500 to-purple-600 px-8 py-6">
                        <h2 class="text-2xl font-bold text-white flex items-center">
                            <svg class="w-8 h-8 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                            Books API Endpoints
                        </h2>
                        <p class="text-blue-100 mt-2">Manage your book collection with full CRUD operations</p>
                    </div>
                    
                    <div class="p-8 space-y-6">
                        <!-- GET & POST /api/v1/books/ -->
                        <div class="border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-all duration-300 hover:border-blue-300">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-3 py-1 text-xs font-bold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-3 py-1 text-xs font-bold text-blue-800 bg-blue-100 rounded-full">POST</span>
                                    <code class="px-4 py-2 bg-gray-100 rounded-lg text-sm font-mono text-indigo-600">/api/v1/books/</code>
                                </div>
                                <div class="flex items-center text-sm text-gray-500">
                                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C20.832 18.477 19.246 18 17.5 18c-1.746 0-3.332.477-4.5 1.253z" />
                                    </svg>
                                    List and create books
                                </div>
                            </div>
                            <p class="text-gray-600 text-sm mb-3">Retrieve all books or create a new book entry with title, author, and publication details.</p>
                            <div class="bg-gray-50 rounded-lg p-3">
                                <code class="text-xs text-gray-700">
                                    GET: Returns paginated list of books<br>
                                    POST: Creates new book (requires: title, author, isbn)
                                </code>
                            </div>
                        </div>
    
                        <!-- GET, PUT, DELETE /api/v1/books/{id}/ -->
                        <div class="border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-all duration-300 hover:border-green-300">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-3 py-1 text-xs font-bold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-3 py-1 text-xs font-bold text-yellow-800 bg-yellow-100 rounded-full">PUT</span>
                                    <span class="px-3 py-1 text-xs font-bold text-red-800 bg-red-100 rounded-full">DELETE</span>
                                    <code class="px-4 py-2 bg-gray-100 rounded-lg text-sm font-mono text-indigo-600">/api/v1/books/{id}/</code>
                                </div>
                                <div class="flex items-center text-sm text-gray-500">
                                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                    </svg>
                                    Manage individual books
                                </div>
                            </div>
                            <p class="text-gray-600 text-sm mb-3">Retrieve, update, or delete a specific book by its unique identifier.</p>
                            <div class="bg-gray-50 rounded-lg p-3">
                                <code class="text-xs text-gray-700">
                                    GET: Returns book details<br>
                                    PUT: Updates book (partial updates supported)<br>
                                    DELETE: Removes book from collection
                                </code>
                            </div>
                        </div>
    
                        <!-- GET /api/v1/books/recent/ -->
                        <div class="border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-all duration-300 hover:border-purple-300">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-3 py-1 text-xs font-bold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <code class="px-4 py-2 bg-gray-100 rounded-lg text-sm font-mono text-indigo-600">/api/v1/books/recent/</code>
                                </div>
                                <div class="flex items-center text-sm text-gray-500">
                                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    List recent books
                                </div>
                            </div>
                            <p class="text-gray-600 text-sm mb-3">Get the most recently added books, sorted by creation date.</p>
                            <div class="bg-gray-50 rounded-lg p-3">
                                <code class="text-xs text-gray-700">
                                    Returns: Latest 10 books by default (configurable with ?limit parameter)
                                </code>
                            </div>
                        </div>
                    </div>
                </div>
    
                <!-- Authentication Section -->
                <div class="bg-white/80 backdrop-blur-lg shadow-2xl rounded-3xl overflow-hidden animate-slide-up" style="animation-delay: 0.2s">
                    <div class="bg-gradient-to-r from-orange-500 to-red-600 px-8 py-6">
                        <h2 class="text-2xl font-bold text-white flex items-center">
                            <svg class="w-8 h-8 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                            </svg>
                            Authentication
                        </h2>
                        <p class="text-orange-100 mt-2">Secure access to protected endpoints</p>
                    </div>
                    
                    <div class="p-8 space-y-6">
                        <!-- Login Endpoint -->
                        <div class="border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-all duration-300 hover:border-orange-300">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-3 py-1 text-xs font-bold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-3 py-1 text-xs font-bold text-blue-800 bg-blue-100 rounded-full">POST</span>
                                    <code class="px-4 py-2 bg-gray-100 rounded-lg text-sm font-mono text-indigo-600">/api-auth/login/</code>
                                </div>
                                <div class="flex items-center text-sm text-gray-500">
                                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                                    </svg>
                                    API authentication login
                                </div>
                            </div>
                            <p class="text-gray-600 text-sm">Authenticate users and obtain session credentials for API access.</p>
                        </div>
    
                        <!-- Logout Endpoint -->
                        <div class="border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-all duration-300 hover:border-red-300">
                            <div class="flex flex-col md:flex-row md:items-center justify-between mb-4">
                                <div class="flex items-center space-x-3 mb-3 md:mb-0">
                                    <span class="px-3 py-1 text-xs font-bold text-green-800 bg-green-100 rounded-full">GET</span>
                                    <span class="px-3 py-1 text-xs font-bold text-blue-800 bg-blue-100 rounded-full">POST</span>
                                    <code class="px-4 py-2 bg-gray-100 rounded-lg text-sm font-mono text-indigo-600">/api-auth/logout/</code>
                                </div>
                                <div class="flex items-center text-sm text-gray-500">
                                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                                    </svg>
                                    API authentication logout
                                </div>
                            </div>
                            <p class="text-gray-600 text-sm">Safely terminate user sessions and invalidate authentication credentials.</p>
                        </div>
                    </div>
                </div>
    
                <!-- Getting Started Section -->
                <div class="bg-gradient-to-br from-indigo-500 to-purple-600 rounded-3xl p-8 text-white animate-slide-up" style="animation-delay: 0.4s">
                    <h2 class="text-2xl font-bold mb-6 flex items-center">
                        <svg class="w-8 h-8 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        Quick Start Guide
                    </h2>
                    
                    <div class="grid md:grid-cols-2 gap-6">
                        <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-6">
                            <h3 class="text-lg font-semibold mb-3">1. Authentication</h3>
                            <code class="text-sm bg-black/20 rounded-lg p-3 block">
                                curl -X POST http://localhost:8000/api-auth/login/ \<br>
                                -d "username=your_username&password=your_password"
                            </code>
                        </div>
                        
                        <div class="bg-white/10 backdrop-blur-sm rounded-2xl p-6">
                            <h3 class="text-lg font-semibold mb-3">2. Fetch Books</h3>
                            <code class="text-sm bg-black/20 rounded-lg p-3 block">
                                curl -X GET http://localhost:8000/api/v1/books/ \<br>
                                -H "Authorization: Bearer your_token"
                            </code>
                        </div>
                    </div>
                </div>
            </div>
    
            <!-- Action Buttons -->
            <div class="text-center mt-16 space-y-4">
                <div class="flex flex-col sm:flex-row gap-4 justify-center">
                    <a href="/api/v1/" class="bg-gradient-to-r from-blue-500 to-purple-600 text-white px-8 py-4 rounded-2xl font-semibold hover:shadow-2xl hover:shadow-blue-500/25 transform hover:scale-105 transition-all duration-300 inline-flex items-center justify-center">
                        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        Try API Now
                    </a>
                    <a href="/" class="bg-white text-gray-900 border-2 border-gray-200 px-8 py-4 rounded-2xl font-semibold hover:border-indigo-300 hover:shadow-xl transform hover:scale-105 transition-all duration-300 inline-flex items-center justify-center">
                        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
