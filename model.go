package main

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	step            step
	projectName     string
	djangoVersion   string
	features        []string // Currently only "vanilla"
	spinner         spinner.Model
	progress        progress.Model
	progressStatus  string
	error           error
	done            bool

	program *tea.Program

	// Single comprehensive form
	mainForm *huh.Form

	// Configuration options
	selectedOptions   []string // For multiselect
	appName           string
	createTemplates   bool // For global templates/static
	createAppTemplates bool // For app-specific templates
	runServer         bool
	initializeGit     bool
	setupTailwind     bool // For Tailwind CSS v4 setup

	stepMessages    []string
	splashCountdown int
	width           int
}

func NewModel() *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginRight(2)

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)

	m := &Model{
		spinner:            s,
		progress:           p,
		step:               stepSplashScreen,
		splashCountdown:    3,
		features:           []string{"vanilla"},
		createTemplates:    true, // Default to Yes
		createAppTemplates: true, // Default to Yes
		runServer:          true, // Default to Yes
		initializeGit:      true, // Default to Yes
		progressStatus:     "Initializing...",
		selectedOptions:    []string{"Global Templates", "Run Server", "Initialize Git"},
	}

	theme := huh.ThemeBase()
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	theme.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	theme.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("102"))
	theme.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)
	theme.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	theme.Blurred.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("102"))

	// Create comprehensive form with all options
	m.mainForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Description("Enter a memorable name for your Django project").
				Value(&m.projectName).
				Validate(validateProjectName),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Django Version").
				Description("Enter Django version (leave empty for latest stable)").
				Placeholder("5.2.0").
				Value(&m.djangoVersion).
				Validate(validateDjangoVersion),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("App Name (Optional)").
				Description("Enter initial Django app name (leave empty to skip)").
				Value(&m.appName),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Project Configuration").
				Description("Select the features you want to include in your Django project").
				Options(
					huh.NewOption("Global Templates & Static Directories", "Global Templates").Selected(true),
					huh.NewOption("App Templates (if creating an app)", "App Templates").Selected(true),
					huh.NewOption("Auto-start Development Server", "Run Server").Selected(true),
					huh.NewOption("Initialize Git Repository", "Initialize Git").Selected(true),
					huh.NewOption("Vanilla + Tailwind CSS v4", "Tailwind"),
				).
				Limit(5).
				Value(&m.selectedOptions),
		),
	).WithTheme(theme)

	return m
}

func (m *Model) SetProgram(p *tea.Program) {
	m.program = p
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
			return tickMsg{} // For splash screen countdown
		}),
		m.spinner.Tick,
	)
}