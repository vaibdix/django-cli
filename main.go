package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type CLIArgs struct {
	ProjectName     string
	DjangoVersion   string
	SkipInteractive bool
	Help            bool
	Install         bool
}

func parseArgs() CLIArgs {
	var args CLIArgs

	flag.StringVar(&args.ProjectName, "name", "", "Project name")
	flag.StringVar(&args.ProjectName, "n", "", "Project name (shorthand)")
	flag.StringVar(&args.DjangoVersion, "version", "", "Django version")
	flag.StringVar(&args.DjangoVersion, "v", "", "Django version (shorthand)")
	flag.BoolVar(&args.SkipInteractive, "auto", false, "Skip interactive mode with defaults")
	flag.BoolVar(&args.Help, "help", false, "Show help")
	flag.BoolVar(&args.Help, "h", false, "Show help (shorthand)")
	flag.BoolVar(&args.Install, "install", false, "Install CLI globally (Windows only)")

	flag.Parse()

	return args
}

func showHelp() {
	fmt.Println(`Django Forge CLI - Interactive Django Project Creator

Usage:
  django-forge [flags]

Flags:
  -n, --name string      Project name
  -v, --version string   Django version (default: latest)
  --auto                 Skip interactive mode with defaults
  --install             Install CLI globally (Windows only)
  -h, --help            Show this help message

Examples:
  django-forge                           # Interactive mode
  django-forge -n myproject              # Set project name
  django-forge -n myproject -v 4.2.7     # Set name and Django version
  django-forge --auto -n myproject       # Non-interactive with defaults
  django-forge --install                 # Install globally on Windows

Config file: ~/.django-forge.json (auto-created with your preferences)`)
}

func main() {
	args := parseArgs()

	if args.Help {
		showHelp()
		return
	}

	if args.Install {
		if err := installOnWindows(); err != nil {
			fmt.Printf("Error during installation: %v\n", err)
			os.Exit(1)
		}
		return
	}

	m := NewModel()

	if args.ProjectName != "" {
		m.projectName = args.ProjectName
	}
	if args.DjangoVersion != "" {
		m.djangoVersion = args.DjangoVersion
	}

	if args.SkipInteractive && args.ProjectName != "" {
		m.step = stepSetup
		go m.CreateProject()
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	m.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
