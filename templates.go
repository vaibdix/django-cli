package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Template variations for different configurations
var (
	basicBaseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{% block title %}{{ project_name|default:"Django Site" }}{% endblock %}</title>
    <style>
        /* Using base with vanilla CSS */
        :root {
            --color-white: #ffffff;
            --color-black: #000000;
            --color-gray-300: #d1d5db;
            --color-gray-400: #9ca3af;
            --color-gray-700: #374151;
            --color-gray-800: #1f2937;
            --color-gray-900: #111827;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
            background-color: var(--color-black);
            color: var(--color-white);
            line-height: 1.5;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }

        header {
            position: sticky;
            top: 0;
            z-index: 50;
            background-color: rgba(0, 0, 0, 0.8);
            backdrop-filter: blur(12px);
            border-bottom: 1px solid var(--color-gray-800);
        }

        .container {
            max-width: 1280px;
            margin: 0 auto;
            padding: 0 1rem;
        }

        .header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem 0;
        }

        .logo {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            text-decoration: none;
            color: var(--color-white);
        }

        .logo svg {
            width: 32px;
            height: 32px;
        }

        .logo h1 {
            font-size: 1.25rem;
            font-weight: 600;
        }

        nav {
            display: none;
        }

        @media (min-width: 768px) {
            nav {
                display: flex;
                gap: 2rem;
            }
        }

        nav a {
            color: var(--color-gray-300);
            text-decoration: none;
            font-size: 0.875rem;
            transition: color 0.2s;
        }

        nav a:hover {
            color: var(--color-white);
        }

        .mobile-menu-button {
            display: block;
            padding: 0.5rem;
            background: none;
            border: none;
            color: var(--color-white);
            cursor: pointer;
        }

        @media (min-width: 768px) {
            .mobile-menu-button {
                display: none;
            }
        }

        main {
            flex: 1;
        }

        footer {
            border-top: 1px solid var(--color-gray-800);
            margin-top: 5rem;
            padding: 4rem 0;
        }

        .footer-content {
            display: grid;
            grid-template-columns: 1fr;
            gap: 2rem;
        }

        @media (min-width: 768px) {
            .footer-content {
                grid-template-columns: 2fr 1fr;
            }
        }

        .footer-brand {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            margin-bottom: 1rem;
        }

        .footer-brand svg {
            width: 24px;
            height: 24px;
        }

        .footer-brand span {
            font-size: 1.125rem;
            font-weight: 600;
        }

        .footer-description {
            color: var(--color-gray-400);
            max-width: 24rem;
            margin-bottom: 1.5rem;
        }

        .footer-links h3 {
            color: var(--color-white);
            font-size: 0.875rem;
            font-weight: 600;
            margin-bottom: 1rem;
        }

        .footer-links ul {
            list-style: none;
        }

        .footer-links li:not(:last-child) {
            margin-bottom: 0.75rem;
        }

        .footer-links a {
            color: var(--color-gray-400);
            text-decoration: none;
            font-size: 0.875rem;
            transition: color 0.2s;
        }

        .footer-links a:hover {
            color: var(--color-white);
        }

        .footer-bottom {
            border-top: 1px solid var(--color-gray-800);
            margin-top: 3rem;
            padding-top: 2rem;
            display: flex;
            flex-direction: column;
            align-items: center;
            gap: 1rem;
        }

        @media (min-width: 768px) {
            .footer-bottom {
                flex-direction: row;
                justify-content: space-between;
            }
        }

        .footer-bottom p {
            color: var(--color-gray-400);
            font-size: 0.875rem;
        }

        .footer-bottom-links {
            display: flex;
            gap: 1.5rem;
        }
    </style>
    {% block extra_head %}{% endblock %}
</head>
<body>
    <header>
        <div class="container">
            <div class="header-content">
                <a href="/" class="logo">
                    <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect width="32" height="32" rx="8" fill="white"/>
                        <path d="M12 8L20 16L12 24" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    </svg>
                    <h1>{{ project_name|default:"Django" }}</h1>
                </a>

                <nav>
                    <a href="/">Home</a>
                    <a href="/admin/">Admin</a>
                </nav>

                <button class="mobile-menu-button">
                    <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                </button>
            </div>
        </div>
    </header>

    <main>
        {% block content %}{% endblock %}
    </main>

    <footer>
        <div class="container">
            <div class="footer-content">
                <div>
                    <div class="footer-brand">
                        <svg width="24" height="24" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <rect width="32" height="32" rx="8" fill="white"/>
                            <path d="M12 8L20 16L12 24" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                        <span>{{ project_name|default:"Django Site" }}</span>
                    </div>
                    <p class="footer-description">The Django framework that gives you everything you need to build full-stack web applications.</p>
                </div>

                <div class="footer-links">
                    <h3>Resources</h3>
                    <ul>
                        <li><a href="/admin/">Admin Panel</a></li>
                        <li><a href="https://docs.djangoproject.com/">Django Docs</a></li>
                    </ul>
                </div>
            </div>

            <div class="footer-bottom">
                <p>© {% now "Y" %} {{ project_name|default:"Django Site" }}. All rights reserved.</p>
                <div class="footer-bottom-links">
                    <a href="#" class="footer-links">Privacy</a>
                    <a href="#" class="footer-links">Terms</a>
                </div>
            </div>
        </div>
    </footer>

    {% block extra_body %}{% endblock %}
</body>
</html>`

	basicIndexTemplate = `{% extends 'base.html' %}

{% block title %}{{ project_name }} - Home{% endblock %}

{% block extra_head %}
<style>
    .hero {
        position: relative;
        overflow: hidden;
    }

    .hero-gradient {
        position: absolute;
        inset: 0;
        background: linear-gradient(to bottom, transparent, var(--color-black));
        pointer-events: none;
    }

    .hero-grid {
        position: absolute;
        inset: 0;
        opacity: 0.2;
        background-image: radial-gradient(rgba(255,255,255,0.1) 1px, transparent 1px);
        background-size: 40px 40px;
    }

    .hero-content {
        position: relative;
        max-width: 1280px;
        margin: 0 auto;
        padding: 5rem 1rem 8rem;
        text-align: center;
    }

    .badge {
        display: inline-flex;
        align-items: center;
        border: 1px solid var(--color-gray-800);
        background-color: rgba(17, 24, 39, 0.5);
        backdrop-filter: blur(4px);
        padding: 0.5rem 1rem;
        border-radius: 9999px;
        font-size: 0.875rem;
        color: var(--color-gray-300);
        margin-bottom: 2rem;
    }

    .hero-title {
        font-size: 3rem;
        font-weight: 700;
        line-height: 1.2;
        margin-bottom: 2rem;
    }

    @media (min-width: 768px) {
        .hero-title {
            font-size: 4rem;
        }
    }

    @media (min-width: 1024px) {
        .hero-title {
            font-size: 6rem;
        }
    }

    .gradient-text {
        background: linear-gradient(to right, #60a5fa, #a78bfa, #f472b6);
        -webkit-background-clip: text;
        background-clip: text;
        color: transparent;
    }

    .hero-subtitle {
        font-size: 1.25rem;
        color: var(--color-gray-400);
        max-width: 48rem;
        margin: 0 auto 3rem;
        line-height: 1.625;
    }

    @media (min-width: 768px) {
        .hero-subtitle {
            font-size: 1.5rem;
        }
    }

    .hero-subtitle strong {
        color: var(--color-white);
    }

    .button-group {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        justify-content: center;
        margin-bottom: 4rem;
    }

    @media (min-width: 640px) {
        .button-group {
            flex-direction: row;
        }
    }

    .button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 1rem 2rem;
        border-radius: 0.375rem;
        font-weight: 600;
        transition: all 0.2s;
        text-decoration: none;
    }

    .button-primary {
        background-color: var(--color-white);
        color: var(--color-black);
    }

    .button-primary:hover {
        background-color: var(--color-gray-300);
    }

    .button-secondary {
        border: 1px solid var(--color-gray-700);
        color: var(--color-white);
    }

    .button-secondary:hover {
        border-color: var(--color-gray-600);
        background-color: rgba(55, 65, 81, 0.1);
    }

    .button svg {
        width: 1rem;
        height: 1rem;
        margin-left: 0.5rem;
    }

    .code-example {
        max-width: 42rem;
        margin: 0 auto;
    }

    .code-window {
        background-color: var(--color-gray-900);
        border: 1px solid var(--color-gray-800);
        border-radius: 0.5rem;
        padding: 1.5rem;
        text-align: left;
    }

    .code-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 1rem;
    }

    .window-dots {
        display: flex;
        gap: 0.5rem;
    }

    .window-dot {
        width: 0.75rem;
        height: 0.75rem;
        border-radius: 9999px;
    }

    .dot-red { background-color: #ef4444; }
    .dot-yellow { background-color: #f59e0b; }
    .dot-green { background-color: #10b981; }

    .window-title {
        color: var(--color-gray-400);
        font-size: 0.875rem;
    }

    .code-content {
        font-family: 'Menlo', 'Monaco', 'Lucida Console', monospace;
        font-size: 0.875rem;
        color: var(--color-gray-300);
    }

    .code-keyword { color: #c084fc; }
    .code-builtin { color: #60a5fa; }
    .code-string { color: #34d399; }
    .code-param { color: #fb923c; }

    .features {
        max-width: 1280px;
        margin: 0 auto;
        padding: 6rem 1rem;
    }

    .features-header {
        text-align: center;
        margin-bottom: 4rem;
    }

    .features-title {
        font-size: 1.875rem;
        font-weight: 700;
        color: var(--color-white);
        margin-bottom: 1rem;
    }

    @media (min-width: 768px) {
        .features-title {
            font-size: 2.25rem;
        }
    }

    .features-subtitle {
        font-size: 1.25rem;
        color: var(--color-gray-400);
        max-width: 42rem;
        margin: 0 auto;
    }

    .features-grid {
        display: grid;
        grid-template-columns: 1fr;
        gap: 2rem;
    }

    @media (min-width: 768px) {
        .features-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }

    .feature-card {
        border: 1px solid var(--color-gray-800);
        background-color: var(--color-gray-900);
        border-radius: 0.5rem;
        padding: 2rem;
        transition: border-color 0.2s;
    }

    .feature-card:hover {
        border-color: var(--color-gray-700);
    }

    .feature-icon {
        width: 3rem;
        height: 3rem;
        background-color: rgba(59, 130, 246, 0.1);
        border-radius: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 1.5rem;
    }

    .feature-icon svg {
        width: 1.5rem;
        height: 1.5rem;
        color: #60a5fa;
    }

    .feature-icon.security {
        background-color: rgba(16, 185, 129, 0.1);
    }

    .feature-icon.security svg {
        color: #34d399;
    }

    .feature-title {
        font-size: 1.25rem;
        font-weight: 600;
        color: var(--color-white);
        margin-bottom: 0.75rem;
    }

    .feature-description {
        color: var(--color-gray-400);
    }
</style>
{% endblock %}

{% block content %}
<div class="hero">
    <div class="hero-gradient"></div>
    <div class="hero-grid"></div>

    <div class="hero-content">
        <div class="badge">
            <span>🚀 Production ready Django application</span>
        </div>

        <h1 class="hero-title">
            <span>The Django</span>
            <span class="gradient-text">Framework</span>
        </h1>

        <p class="hero-subtitle">
            Django provides everything you need to build fast, secure, and scalable web applications.
            <strong>Used by thousands of developers worldwide.</strong>
        </p>

        <div class="button-group">
            <a href="https://docs.djangoproject.com/" class="button button-primary">
                Get Started
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
            </a>
            <a href="/admin/" class="button button-secondary">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Admin Panel
            </a>
        </div>

        <div class="code-example">
            <div class="code-window">
                <div class="code-header">
                    <div class="window-dots">
                        <div class="window-dot dot-red"></div>
                        <div class="window-dot dot-yellow"></div>
                        <div class="window-dot dot-green"></div>
                    </div>
                    <span class="window-title">Django Project</span>
                </div>
                <pre class="code-content"><code><span class="code-keyword">from</span> <span class="code-builtin">django.shortcuts</span> <span class="code-keyword">import</span> <span class="code-builtin">render</span>

<span class="code-keyword">def</span> <span class="code-builtin">home</span>(<span class="code-param">request</span>):
    <span class="code-keyword">return</span> <span class="code-builtin">render</span>(request, <span class="code-string">'index.html'</span>)</code></pre>
            </div>
        </div>
    </div>
</div>

<div class="features">
    <div class="features-header">
        <h2 class="features-title">Why Django?</h2>
        <p class="features-subtitle">
            Built for speed, security, and scalability. Trusted by startups and enterprises.
        </p>
    </div>

    <div class="features-grid">
        <div class="feature-card">
            <div class="feature-icon">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
            </div>
            <h3 class="feature-title">Fast Development</h3>
            <p class="feature-description">Django's batteries-included approach means you can build full-featured applications quickly without reinventing the wheel.</p>
        </div>

        <div class="feature-card">
            <div class="feature-icon security">
                <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
            </div>
            <h3 class="feature-title">Security First</h3>
            <p class="feature-description">Built-in protection against common security threats like SQL injection, CSRF, and XSS attacks.</p>
        </div>
    </div>
</div>
{% endblock %}`

	basicApiDocsTemplate = `{% extends 'base.html' %}

{% block title %}{{ project_name }} - API Documentation{% endblock %}

{% block content %}
<div>
    <h1>API Documentation</h1>
    
    <section>
        <h2>Authentication</h2>
        <p>Learn how to authenticate with our API.</p>
        
        <h3>Endpoints</h3>
        <ul>
            <li>
                <strong>Login:</strong> 
                <code>POST /api-auth/login/</code>
            </li>
            <li>
                <strong>Logout:</strong> 
                <code>POST /api-auth/logout/</code>
            </li>
        </ul>
    </section>

    <section>
        <h2>API Endpoints</h2>
        <p>Explore our available API endpoints.</p>
        
        <div>
            <h3>GET /api/v1/</h3>
            <p>Root endpoint that lists all available resources.</p>
        </div>
    </section>
</div>
{% endblock %}`

	tailwindBaseTemplateNoApi = `{% load static %}
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
                        <a href="/admin/" class="text-gray-300 hover:text-white transition-colors duration-200 text-sm">Admin</a>
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
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
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
                            <li><a href="/admin/" class="text-gray-400 hover:text-white transition-colors text-sm">Admin Panel</a></li>
                            <li><a href="https://docs.djangoproject.com/" class="text-gray-400 hover:text-white transition-colors text-sm">Django Docs</a></li>
                        </ul>
                    </div>
                </div>

                <div class="border-t border-gray-800 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
                    <p class="text-gray-400 text-sm">© {% now "Y" %} {{ project_name|default:"Django Site" }}. All rights reserved.</p>
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

	tailwindIndexTemplateNoApi = `{% extends 'base.html' %}
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
                    <a href="https://docs.djangoproject.com/" class="bg-white text-black px-8 py-4 rounded-md font-semibold hover:bg-gray-200 transition-colors duration-200 inline-flex items-center justify-center">
                        Get Started
                        <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                        </svg>
                    </a>
                    <a href="/admin/" class="border border-gray-700 text-white px-8 py-4 rounded-md font-semibold hover:border-gray-600 hover:bg-gray-900 transition-colors duration-200 inline-flex items-center justify-center">
                        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        Admin Panel
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
                        <pre class="text-sm text-gray-300 font-geist-mono"><code><span class="text-purple-400">from</span> <span class="text-blue-400">django.shortcuts</span> <span class="text-purple-400">import</span> <span class="text-yellow-400">render</span>

<span class="text-purple-400">def</span> <span class="text-blue-400">home</span>(<span class="text-orange-400">request</span>):
    <span class="text-purple-400">return</span> <span class="text-yellow-400">render</span>(request, <span class="text-green-400">'index.html'</span>)</code></pre>
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

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
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
        </div>
    </div>
</div>
{% endblock %}`

	tailwindBaseTemplate = `{% load static %}
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
                        <a href="/admin/" class="text-gray-300 hover:text-white transition-colors duration-200 text-sm">Admin</a>
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
                <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
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
                            <li><a href="/admin/" class="text-gray-400 hover:text-white transition-colors text-sm">Admin Panel</a></li>
                            <li><a href="https://docs.djangoproject.com/" class="text-gray-400 hover:text-white transition-colors text-sm">Django Docs</a></li>
                        </ul>
                    </div>
                </div>

                <div class="border-t border-gray-800 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
                    <p class="text-gray-400 text-sm">© {% now "Y" %} {{ project_name|default:"Django Site" }}. All rights reserved.</p>
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

	tailwindIndexTemplate = `{% extends 'base.html' %}
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

	tailwindApiDocsTemplate = `{% extends 'base.html' %}
{% block title %}{{ project_name }} - API Documentation{% endblock %}

{% block content %}
<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="text-center mb-16">
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">API Documentation</h1>
        <p class="text-xl text-gray-400 max-w-2xl mx-auto">
            Explore our comprehensive API documentation to integrate with our services.
        </p>
    </div>

    <!-- Authentication Section -->
    <section class="mb-16">
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
                <h3 class="text-lg font-semibold text-white mb-3">2. Make API Requests</h3>
                <div class="bg-black/50 border border-gray-700 rounded-lg p-4">
                    <pre class="text-sm text-gray-300 font-mono overflow-x-auto">
<span class="text-purple-400">curl</span> <span class="text-blue-400">-X GET</span> http://localhost:8000/api/v1/ \
  <span class="text-blue-400">-H</span> <span class="text-green-400">"Authorization: Bearer your_token"</span></pre>
                </div>
            </div>
        </div>
    </section>
</div>
{% endblock %}`
)

// Function to get the appropriate template content based on configuration
func getTemplateContent(templateType string, config map[string]bool) string {
	switch templateType {
	case "base.html":
		if config["tailwind"] {
			if config["api"] {
				return tailwindBaseTemplate
			}
			return tailwindBaseTemplateNoApi
		}
		return basicBaseTemplate
	case "index.html":
		if config["tailwind"] {
			if config["api"] {
				return tailwindIndexTemplate
			}
			return tailwindIndexTemplateNoApi
		}
		return basicIndexTemplate
	case "api-docs.html":
		if config["tailwind"] {
			return tailwindApiDocsTemplate
		}
		return basicApiDocsTemplate
	default:
		return ""
	}
}

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

	// Create configuration map for template selection
	config := map[string]bool{
		"tailwind": m.setupTailwind,
		"api":      m.setupRestFramework,
	}

	// Write base.html
	baseContent := getTemplateContent("base.html", config)
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "base.html"), []byte(baseContent), 0644); err != nil {
		return fmt.Errorf("failed to create base.html: %v", err)
	}

	// Write index.html
	indexContent := getTemplateContent("index.html", config)
	if err := os.WriteFile(filepath.Join(globalTemplatesPath, "index.html"), []byte(indexContent), 0644); err != nil {
		return fmt.Errorf("failed to create index.html: %v", err)
	}

	// Write api-docs.html if REST framework is enabled
	if m.setupRestFramework {
		apiDocsContent := getTemplateContent("api-docs.html", config)
		if err := os.WriteFile(filepath.Join(globalTemplatesPath, "api-docs.html"), []byte(apiDocsContent), 0644); err != nil {
			return fmt.Errorf("failed to create api-docs.html: %v", err)
		}
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
    template_name = 'index.html'

class ApiDocsView(TemplateView):
    template_name = 'api-docs.html'`

	viewsPath := filepath.Join(projectPath, m.projectName, "views.py")
	if err := os.WriteFile(viewsPath, []byte(viewsContent), 0644); err != nil {
		return fmt.Errorf("failed to create views.py: %v", err)
	}

	// Update urls.py
	urlsContent := fmt.Sprintf(`from django.contrib import admin
from django.urls import path, include
from . import views

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),  # Root route is always available
    path('admin/', admin.site.urls),
    path('api-docs/', views.ApiDocsView.as_view(), name='api_docs'),
    path('__reload__/', include('django_browser_reload.urls')),
]`)

	urlsPath := filepath.Join(projectPath, m.projectName, "urls.py")
	if err := os.WriteFile(urlsPath, []byte(urlsContent), 0644); err != nil {
		return fmt.Errorf("failed to create urls.py: %v", err)
	}

	return nil
}
