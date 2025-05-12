package main

import (
	"fmt"
	"time"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type splashDoneMsg struct{}
type tickMsg struct{}

type Model struct {
	step           step
	projectName    string
	djangoVersion  string
	features       []string
	spinner        spinner.Model
	progress       progress.Model
	progressStatus string
	error          error
	done           bool
	doneChan       chan bool

	inputForm      *huh.Form
	versionForm    *huh.Form
	featureForm    *huh.Form
	templateForm   *huh.Form
	appForm        *huh.Form
	appTemplateForm *huh.Form
	serverForm     *huh.Form
	appName        string
	createTemplates bool
	createAppTemplates bool
	runServer      bool
	stepMessages   []string 
	splashCountdown int      
	width          int       
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
		progress.WithoutPercentage(),
	)

	m := &Model{
		spinner:  s,
		progress: p,
		doneChan: make(chan bool),
		step:     stepSplashScreen, 
		splashCountdown: 3, 
		features: []string{"vanilla"}, 
		createTemplates: false,
		createAppTemplates: false,
		runServer: false,
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
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Project name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme)

	m.versionForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Django version").
				Description("Press Enter to use default version (5.2.0)").
				Value(&m.djangoVersion),
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

	m.appForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Create Django App").
				Description("Enter app name (optional, press Enter to skip)").
				Value(&m.appName),
		),
	).WithTheme(theme)

	m.templateForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Django Templates").
				Description("Would you like to set up template directories?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.createTemplates),
		),
	).WithTheme(theme)

	m.appTemplateForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("App Templates").
				Description("Would you like to set up templates for this app?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.createAppTemplates),
		),
	).WithTheme(theme)

	m.serverForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Title("Run Development Server").
				Description("Would you like to start the development server?").
				Options(
					huh.NewOption("Yes", true),
					huh.NewOption("No", false),
				).
				Value(&m.runServer),
		),
	).WithTheme(theme)

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		tea.Tick(1*time.Second, func(_ time.Time) tea.Msg { 
			return tickMsg{}
		}),
		m.inputForm.Init(), 
	)
}

func (m *Model) updateProgress() tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case <-m.doneChan:
				return progressMsg(1.0)
			default:
				time.Sleep(100 * time.Millisecond)
				return progressMsg(m.progress.Percent() + 0.1)
			}
		}
	}
}