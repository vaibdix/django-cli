# Django Forge CLI 🚀

An interactive command-line tool for creating Django projects with modern development practices and automated setup.

## Features

-   **Dynamic Progress Bar**: Provides real-time feedback on the project setup progress, adapting to the selected features and steps.

### 🎯 Interactive Project Setup

-   **Project Name Validation**: Ensures valid project names and prevents conflicts
-   **Django Version Selection**: Choose specific Django versions or use latest stable
-   **App Creation**: Optionally create an initial Django app during setup
-   **Multi-select Configuration**: Choose features you want in your project

### 📁 Project Structure & Templates

-   **Global Templates**: Creates `templates/` directory with base.html and index.html
-   **Static Files**: Sets up `static/css/` and `static/js/` directories with starter files
-   **App Templates**: Creates app-specific template directories when creating apps
-   **Django Settings**: Automatically configures `settings.py` for templates and static files

### 🔧 Development Environment

-   **Virtual Environment**: Automatically creates `.venv` using `uv` (preferred) or `python -m venv`
-   **Django Installation**: Installs specified Django version and `django-browser-reload`
-   **Hot Reload**: Configures `django-browser-reload` for automatic browser refresh during development
-   **Development Server**: Optionally starts the Django development server after setup

### 🗂️ Version Control

-   **Git Initialization**: Optionally initializes a Git repository
-   **Gitignore**: Creates a comprehensive `.gitignore` file for Django projects

### 🎨 User Experience

-   **Dynamic Progress Bar**: Accurately tracks and displays project creation progress based on selected features and completed steps.
-   **Animated Progress Bar**: Visual feedback during project creation
-   **Splash Screen**: Welcome screen with countdown
-   **Error Handling**: Comprehensive validation and error messages
-   **Cross-platform**: Works on macOS, Linux, and Windows

## Installation

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd django-cli

# Build for your platform
# Linux/macOS:
go build -o django-cli

# Windows:
go build -o django-cli.exe

# Cross-compile for Windows from macOS/Linux:
GOOS=windows GOARCH=amd64 go build -o django-cli-windows.exe
```

## Usage

### Interactive Mode (Recommended)

```bash
# Run the interactive CLI
./django-cli
```

The interactive mode will guide you through:

1. **Project Name**: Enter a unique name for your Django project
2. **Django Version**: Specify version (e.g., "5.2.0") or leave empty for latest
3. **App Name**: Optionally create an initial Django app
4. **Project Configuration**: Select features using multi-select:
    - Global Templates & Static Directories
    - App Templates (if creating an app)
    - Auto-start Development Server
    - Initialize Git Repository

### Command Line Arguments

```bash
# Set project name
./django-cli -name myproject
./django-cli -n myproject

# Set Django version
./django-cli -version 4.2.7
./django-cli -v 4.2.7

# Combine flags
./django-cli -n myproject -v 5.2.0

# Non-interactive mode with defaults
./django-cli --auto -n myproject

# Show help
./django-cli -h
./django-cli --help
```

### Available Flags

| Flag        | Short | Description                         |
| ----------- | ----- | ----------------------------------- |
| `--name`    | `-n`  | Project name                        |
| `--version` | `-v`  | Django version (default: latest)    |
| `--auto`    |       | Skip interactive mode with defaults |
| `--help`    | `-h`  | Show help message                   |

## Project Structure Created

```
myproject/
├── .venv/                    # Virtual environment
├── .git/                     # Git repository (optional)
├── .gitignore               # Django-specific gitignore
├── manage.py                # Django management script
├── myproject/               # Django project directory
│   ├── __init__.py
│   ├── settings.py          # Configured with templates and middleware
│   ├── urls.py              # Updated to include app URLs (if app created)
│   ├── wsgi.py
│   └── asgi.py
├── templates/               # Global templates (optional)
│   ├── base.html            # Base template with django-browser-reload
│   └── index.html           # Homepage template
├── static/                  # Static files (optional)
│   ├── css/
│   │   └── style.css        # Basic styling
│   └── js/
│       └── main.js          # JavaScript starter
└── myapp/                   # Django app (optional)
    ├── __init__.py
    ├── admin.py
    ├── apps.py
    ├── models.py
    ├── tests.py
    ├── views.py             # With homepage view
    ├── urls.py              # App URL configuration
    └── templates/           # App-specific templates (optional)
        └── myapp/
            └── index.html
```

## Dependencies

### System Requirements

-   **Go**: 1.23.0 or later
-   **Python**: 3.8 or later (`python3` or `python` in PATH)
-   **uv** (optional but recommended): For faster virtual environment and package management

### Go Dependencies

-   `github.com/charmbracelet/bubbletea`: TUI framework
-   `github.com/charmbracelet/bubbles`: TUI components
-   `github.com/charmbracelet/huh`: Form components
-   `github.com/charmbracelet/lipgloss`: Styling

## Configuration

The tool supports a configuration file at `~/.django-forge.json` for setting defaults:

```json
{
    "default_django_version": "latest",
    "default_features": ["vanilla"],
    "create_templates": true,
    "create_app_templates": true,
    "run_server": true,
    "initialize_git": true,
    "prefer_uv": true
}
```

## Features in Detail

### Django Browser Reload

Automatically configured for development:

-   Added to `INSTALLED_APPS`
-   Middleware configured
-   Template tags included in base template
-   Enables automatic browser refresh when files change

### Template Configuration

-   Global templates directory added to `DIRS` in `settings.py`
-   Static files configuration updated
-   Base template includes:
    -   Responsive meta tags
    -   CSS and JavaScript loading
    -   Django browser reload integration
    -   Block structure for inheritance

### Git Integration

-   Initializes Git repository
-   Creates comprehensive `.gitignore` for Django projects
-   Excludes virtual environment, cache files, database, and IDE files

## Troubleshooting

### Common Issues

1. **Python not found**: Ensure `python3` or `python` is in your PATH
2. **Permission denied**: Make the binary executable with `chmod +x django-cli`
3. **Virtual environment creation fails**: Install `python3-venv` on Ubuntu/Debian
4. **uv not found**: Install uv with `pip install uv` or use system Python

### Error Messages

-   **Project name validation**: Checks for invalid characters and Python reserved words
-   **Directory exists**: Prevents overwriting existing projects
-   **Django version format**: Validates version format (e.g., "4.2.0")

## Development

### Building

```bash
# Development build
go run .

# Production build
go build -o django-cli

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o django-cli-linux
GOOS=windows GOARCH=amd64 go build -o django-cli-windows.exe
GOOS=darwin GOARCH=amd64 go build -o django-cli-macos
```

### Testing

```bash
# Test with flags
./django-cli -n testproject -v 5.0

# Test interactive mode
./django-cli
```

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

---

**Django Forge CLI** - Crafting your Django project, one step at a time. 🚀
