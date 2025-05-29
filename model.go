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
	djangoVersion string
	features        []string // Currently only "vanilla"
	spinner         spinner.Model
	progress        progress.Model
	progressStatus  string
	error           error
	done            bool

	program *tea.Program

	inputForm    *huh.Form
	versionForm  *huh.Form
	featureForm  *huh.Form
	templateForm *huh.Form // For global templates

	// For app creation and its templates
	appNameInput      *huh.Input        // Store the input field for appName
	appForm           *huh.Form         // Form that contains appNameInput
	appTemplateSelect *huh.Select[bool] // Store the select field for app templates
	appTemplateForm   *huh.Form         // Form that contains appTemplateSelect

	serverForm *huh.Form
	gitForm    *huh.Form

	appName            string
	createTemplates    bool // For global templates/static
	createAppTemplates bool // For app-specific templates
	runServer          bool
	initializeGit      bool

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
		progress.WithGradient("#7D56F4", "#41E296"),
		progress.WithWidth(50),
		// Removed progress.WithoutPercentage() to enable percentage display
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
	}

	theme := huh.ThemeBase()
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	theme.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	theme.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true)
	theme.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

	m.inputForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("Enter a name for your Django project").
				Value(&m.projectName).
				Validate(validateProjectName), // Use the validation function
		),
	).WithTheme(theme)

	m.versionForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Django version").
				Description("Press Enter to use default version (e.g., 5.2.0 or latest stable).").
				Placeholder("latest").
				Value(&m.djangoVersion).
				Validate(validateDjangoVersion), // Use the validation function
		),
	).WithTheme(theme)

	m.featureForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setup Type").
				Description("Choose your Django setup type").
				Options(
					huh.NewOption("Vanilla Setup 🍦", "vanilla"),
				).
				Value(&m.features[0]),
		),
	).WithTheme(theme)

	m.templateForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Global Templates & Static").
				Description("Set up global 'templates' and 'static' directories? (Recommended)").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.createTemplates),
		),
	).WithTheme(theme)

	// App Name Input Field and Form
	m.appNameInput = huh.NewInput().
		Title("Create Initial Django App"). // Static title for the input field itself
		Description("Enter app name (optional, press Enter to skip)").
		Value(&m.appName)
	m.appForm = huh.NewForm(
		huh.NewGroup(m.appNameInput),
	).WithTheme(theme)

	// App Template Select Field and Form
	m.appTemplateSelect = huh.NewSelect[bool]().
		Title("App Templates"). // Generic initial title, will be updated
		Description("Set up basic templates and views for this app?").
		Options(
			huh.NewOption("Yes", true),
			huh.NewOption("No", false),
		).
		Value(&m.createAppTemplates)
	m.appTemplateForm = huh.NewForm(
		huh.NewGroup(m.appTemplateSelect),
	).WithTheme(theme)

	m.serverForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Run Development Server").
				Description("Start the Django development server after setup? (python manage.py runserver)").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.runServer),
		),
	).WithTheme(theme)

	m.gitForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Initialize Git Repository").
				Description("Initialize a Git repository and create a .gitignore file?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.initializeGit),
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