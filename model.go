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
	step               step
	projectName        string
	djangoVersion      string
	features           []string
	spinner            spinner.Model
	progress           progress.Model
	progressStatus     string
	error              error
	done               bool
	program            *tea.Program
	mainForm           *huh.Form
	devServerForm      *huh.Form
	selectedOptions    []string
	appName            string
	createTemplates    bool
	createAppTemplates bool
	runServer          bool
	initializeGit      bool
	setupTailwind      bool
	startDevServer     bool
	stepMessages       []string
	splashCountdown    int
	width              int
	totalSteps         int
	completedSteps     int
}

func (m *Model) calculateTotalSteps() int {
	// Base steps: project dir, venv, django, settings
	totalSteps := 4

	// Optional features
	if m.createTemplates {
		totalSteps++
	}
	if m.appName != "" {
		totalSteps++
	}
	if m.initializeGit {
		totalSteps++
	}
	if m.setupTailwind {
		totalSteps++
	}

	return totalSteps
}

func (m *Model) updateProgress(status string) {
	m.completedSteps++
	progress := float64(m.completedSteps) / float64(m.totalSteps)
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: progress, status: status})
	}
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
		createTemplates:    true,
		createAppTemplates: true,
		runServer:          false,
		initializeGit:      true,
		progressStatus:     "Initializing...",
		selectedOptions:    []string{"Global Templates", "Initialize Git"},
		completedSteps:     0,
	}

	theme := huh.ThemeBase()
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	theme.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	theme.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("102"))
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	theme.Blurred.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	theme.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)
	theme.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	theme.Blurred.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("102"))

	m.mainForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
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
					huh.NewOption("Initialize Git Repository", "Initialize Git").Selected(true),
					huh.NewOption("Vanilla + Tailwind CSS v4", "Tailwind"),
				).
				Limit(4).
				Value(&m.selectedOptions),
		),
	).WithTheme(theme)

	m.devServerForm = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Open and run in VS Code?").
				Affirmative("Yes").
				Negative("No").
				Value(&m.startDevServer),
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
			return tickMsg{}
		}),
		m.spinner.Tick,
	)
}
