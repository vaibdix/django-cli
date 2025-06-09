# Django Forge CLI

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/django-cli)](https://goreportcard.com/report/github.com/yourusername/django-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Django Forge CLI is an interactive command-line tool that streamlines Django project creation and setup. It provides a modern, intuitive interface for creating Django projects with best practices baked in.

![Django Forge CLI Demo](demo.gif)

## Features

- 🚀 Interactive project setup with a beautiful TUI
- ⚡️ Lightning-fast project creation with uv for dependency management
- 🎨 Built-in Tailwind CSS integration
- 🔄 Django REST Framework setup option
- 📝 Automatic VS Code configuration
- 🛠️ Multiple development server support
- 🎯 Git repository initialization
- 🔧 Customizable project templates
- 💻 Cross-platform support (Windows, macOS, Linux)

## Installation

### Prerequisites

- Go 1.21 or higher
- Python 3.8 or higher

### Using Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/yourusername/django-cli/releases).

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
git clone https://github.com/yourusername/django-cli.git
cd django-cli

# Build the binary
go build -o django-cli

# Optional: Install globally
mv django-cli /usr/local/bin/
```

## Usage

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
# Interactive mode
django-cli

# Set project name
django-cli -n myproject

# Set name and Django version
django-cli -n myproject -v 4.2.7

# Non-interactive with defaults
django-cli --auto -n myproject
```

## Configuration

Django CLI creates a configuration file at `~/.django-forge.json` to store your preferences. You can modify this file to set default values for:

- Django version
- Project structure
- Template settings
- Development server configuration
- VS Code settings

Example configuration:

```json
{
  "defaultDjangoVersion": "4.2.7",
  "defaultProjectStructure": "recommended",
  "useTailwind": true,
  "useRestFramework": true,
  "autoOpenVSCode": true
}
```

## Project Structure

The generated project follows modern Django best practices:

```
myproject/
├── .vscode/                # VS Code configuration
├── manage.py
├── myproject/             # Project configuration
│   ├── __init__.py
│   ├── settings.py
│   ├── urls.py
│   └── wsgi.py
├── apps/                  # Django applications
├── templates/            # Global templates
│   ├── base.html
│   └── index.html
├── static/              # Static files
│   ├── css/
│   └── js/
└── requirements.txt     # Project dependencies
```


## Development

### Requirements

- Go 1.21+
- Make (optional, for using Makefile)

### Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/django-cli.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build for development
go build
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes and version history.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the terminal UI
- [Django](https://www.djangoproject.com/) community
- All contributors who have helped shape this project

## Support

- 📫 Report bugs through [GitHub issues](https://github.com/yourusername/django-cli/issues)
- 💬 Get help in the [Discussions](https://github.com/yourusername/django-cli/discussions)
- 📖 Read the [Wiki](https://github.com/yourusername/django-cli/wiki) for detailed documentation
