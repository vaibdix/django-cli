# Django Forge CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/vaibdix/django-cli)](https://goreportcard.com/report/github.com/vaibdix/django-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Django Forge CLI is an interactive command-line tool that streamlines Django project creation and setup. It provides a modern, intuitive interface for creating Django projects with best practices baked in.

## ✨ Features

- 🚀 Interactive project setup with a beautiful TUI
- ⚡️ Lightning-fast project creation with automatic `uv` detection for dependency management
- 🎨 Built-in Tailwind CSS v4 integration
- 🔄 Django REST Framework setup option with sample API
- 📝 Automatic VS Code configuration with optimized tasks
- 🛠️ Multiple development server support
- 🎯 Git repository initialization
- 🔧 Customizable project templates
- 💻 Cross-platform support (Windows, macOS, Linux)
- 🏃‍♂️ Windows performance optimizations for faster setup
- 🔄 Automatic fallback from `uv` to `pip` when needed

## 🚀 Performance

Django Forge CLI automatically detects and uses [`uv`](https://github.com/astral-sh/uv) for faster package management:

| Operation | With `pip` | With `uv` | Improvement |
|-----------|------------|-----------|-------------|
| Virtual environment creation | 5-10s | 0.5-1s | **10x faster** |
| Django installation | 15-30s | 2-5s | **10x faster** |
| Total setup time | 30-60s | 3-10s | **Up to 20x faster** |

*Performance improvements are especially noticeable on Windows systems.*

## 📦 Installation

### Prerequisites

- **Go 1.21** or higher
- **Python 3.8** or higher
- **uv** (optional, but recommended for better performance)

### Install uv for Better Performance

```bash
# Install uv for significantly faster package management
pip install uv
```

### Using Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/vaibdix/django-cli/releases).

#### Windows
```bash
# Download and install globally
django-cli-windows.exe --install
```

#### macOS/Linux
```bash
# Make the binary executable
chmod +x django-cli-linux

# Move to a directory in your PATH
sudo mv django-cli-linux /usr/local/bin/django-cli
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/vaibdix/django-cli.git
cd django-cli

# Build the binary
go build -o django-cli

# Optional: Install globally
mv django-cli /usr/local/bin/
```

## 🎯 Usage

### Interactive Mode

Simply run the CLI without any arguments for interactive mode:

```bash
django-cli
```

### Command-line Arguments

```bash
django-cli [flags]

Flags:
  -n, --name string      Project name
  -v, --version string   Django version (default: latest)
  --auto                 Skip interactive mode with defaults
  --install             Install CLI globally (Windows only)
  -h, --help            Show this help message
```

### Examples

```bash
# Interactive mode (recommended)
django-cli

# Set project name
django-cli -n myproject

# Set name and Django version
django-cli -n myproject -v 5.2

# Non-interactive with defaults
django-cli --auto -n myproject
```

## ⚙️ Configuration

Django CLI creates a configuration file at `~/.django-forge.json` to store your preferences. You can modify this file to set default values for:

- Django version
- Project structure
- Template settings
- Development server configuration
- VS Code settings
- Package manager preferences

Example configuration:

```json
{
  "defaultDjangoVersion": "5.2",
  "defaultProjectStructure": "recommended",
  "useTailwind": true,
  "useRestFramework": true,
  "autoOpenVSCode": true,
  "preferUV": true
}
```

## 📁 Project Structure

The generated project follows modern Django best practices:

```
myproject/
├── .vscode/                # VS Code configuration with optimized tasks
│   ├── tasks.json         # Auto-configured development tasks
│   └── settings.json      # Project-specific settings
├── .venv/                 # Virtual environment
├── manage.py
├── myproject/             # Project configuration
│   ├── __init__.py
│   ├── settings.py        # Optimized settings with auto-reload
│   ├── urls.py
│   ├── wsgi.py
│   └── api.py            # REST API configuration (if enabled)
├── apps/                  # Django applications directory
├── templates/            # Global templates
│   ├── base.html
│   ├── index.html
│   └── api-docs.html     # API documentation (if REST enabled)
├── static/              # Static files
│   ├── css/
│   │   └── styles.css   # Tailwind CSS (if enabled)
│   └── js/
├── requirements.txt     # Project dependencies
├── WELCOME.md          # Getting started guide
└── .gitignore          # Git ignore file (if Git enabled)
```

## 🛠️ Development Environment

### VS Code Integration

Django Forge CLI automatically configures VS Code with:

- **Optimized development tasks** using `uv` when available
- **Auto-starting development server** on project open
- **Tailwind CSS watcher** (if enabled)
- **Python environment detection**
- **Debugging configuration**

### Available VS Code Tasks

- `Django: Run server` - Start the Django development server
- `Tailwind: Watch CSS` - Watch and compile Tailwind CSS (if enabled)
- `Start Development Environment` - Run both Django and Tailwind concurrently

## 🔧 REST API Features

When Django REST Framework is enabled, Django Forge CLI creates:

- **Complete API setup** with ViewSets and Serializers
- **Sample Book model** with CRUD operations
- **API documentation** accessible at `/api-docs/`
- **Browsable API** at `/api/v1/`
- **Management command** for creating sample data
- **Authentication endpoints** at `/api-auth/`

### API Endpoints

```
GET    /api/v1/books/          # List all books
POST   /api/v1/books/          # Create a new book
GET    /api/v1/books/{id}/     # Retrieve a specific book
PUT    /api/v1/books/{id}/     # Update a specific book
DELETE /api/v1/books/{id}/     # Delete a specific book
GET    /api/v1/books/recent/   # Get recent books (custom action)
```

### Setup

```bash
# Clone the repository
git clone https://github.com/vaibdix/django-cli.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build
```

### Performance Testing

To test performance improvements:

```bash
# Without uv
time django-cli --auto -n test-project

# With uv installed
pip install uv
time django-cli --auto -n test-project-uv
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📋 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.

## 🙏 Acknowledgments

- [uv](https://github.com/astral-sh/uv) for blazing-fast Python package management
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the terminal UI
- [Django](https://www.djangoproject.com/) community
- All contributors who have helped shape this project

## 💡 Support

- 📫 Report bugs through [GitHub issues](https://github.com/vaibdix/django-cli/issues)
