package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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
	projectNameForm    *huh.Form
	djangoVersionForm  *huh.Form
	projectConfigForm  *huh.Form
	appNameForm        *huh.Form
	devServerForm      *huh.Form
	selectedOptions    []string
	appName            string
	createTemplates    bool
	createAppTemplates bool
	runServer          bool
	initializeGit      bool
	setupTailwind      bool
	setupRestFramework bool
	useGlobalTemplates bool
	startDevServer     bool
	stepMessages       []string
	splashCountdown    int
	width              int
	totalSteps         int
	completedSteps     int
}

func (m *Model) calculateTotalSteps() int {
	steps := 3 // Base steps (create directory, venv, Django)

	if m.createTemplates {
		steps++
	}
	if m.appName != "" {
		steps++
	}
	if m.initializeGit {
		steps++
	}
	if m.setupTailwind {
		steps++
	}
	if m.setupRestFramework {
		steps++
	}
	steps++ // For migrations (makemigrations and migrate)

	return steps
}

func (m *Model) updateProgress(status string) {
	m.completedSteps++
	progress := float64(m.completedSteps) / float64(m.totalSteps)
	if m.program != nil {
		m.program.Send(projectProgressMsg{percent: progress, status: status})
	}
}

func (m *Model) getActiveForm() *huh.Form {
	switch m.step {
	case stepProjectName:
		return m.projectNameForm
	case stepDjangoVersion:
		return m.djangoVersionForm
	case stepProjectConfig:
		return m.projectConfigForm
	case stepAppName:
		return m.appNameForm
	case stepDevServerPrompt:
		return m.devServerForm
	default:
		return nil
	}
}

// func NewModel() *Model {
// 	s := GetSpinner()

// 	p := progress.New(
// 		progress.WithGradient("#7a6483", "#baa4ed"),
// 		// progress.WithDefaultGradient(),
// 		progress.WithWidth(50),
// 	)
// 	p.Empty = '░'
func NewModel() *Model {
	s := GetSpinner()
	p := progress.New(
		progress.WithGradient("#baa4ed", "#ffbc9e"),
		// progress.WithDefaultGradient(),
		progress.WithWidth(50),

		progress.WithFillCharacters('●', '○'), // moved this line up and added comma
	)
	p.ShowPercentage = true
	// p.Empty = '░'
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
		selectedOptions:    []string{"Standard Django Project", "Initialize Git"},
		completedSteps:     0,
	}

	theme := GetTheme()

	m.projectNameForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Value(&m.projectName).
				Validate(validateProjectName),
		),
	).WithTheme(theme)

	m.djangoVersionForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Django Version").
				Placeholder("5.2.1").
				Value(&m.djangoVersion).
				Validate(validateDjangoVersion),
		),
	).WithTheme(theme)

	m.projectConfigForm = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Project Configuration" + "\n").
				Options(
					huh.NewOption("Standard Django Project (includes templates & static files)", "Standard Django Project").Selected(true),
					huh.NewOption("Initialize Git Repository", "Initialize Git").Selected(true),
					huh.NewOption("Tailwind CSS v4", "Tailwind"),
					huh.NewOption("Django REST Framework API", "REST Framework"),
				).
				Limit(4).
				Value(&m.selectedOptions),
		),
	).WithTheme(theme)

	m.appNameForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter app name").
				Value(&m.appName),
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

func (m *Model) processFormData() {
	if m.djangoVersion == "" {
		m.djangoVersion = "latest"
	}
	m.createTemplates = false
	m.createAppTemplates = false
	m.initializeGit = false
	m.setupTailwind = false
	m.setupRestFramework = false

	for _, opt := range m.selectedOptions {
		switch opt {
		case "Standard Django Project":
			m.createTemplates = true
			m.createAppTemplates = true
		case "Initialize Git":
			m.initializeGit = true
		case "Tailwind":
			m.setupTailwind = true
		case "REST Framework":
			m.setupRestFramework = true
		}
	}
	m.stepMessages = append(m.stepMessages, "Project name: "+m.projectName)
	m.stepMessages = append(m.stepMessages, "Django version: "+m.djangoVersion)
	if m.appName != "" {
		m.stepMessages = append(m.stepMessages, "App name: "+m.appName)
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Selected options: %v", m.selectedOptions))
}
