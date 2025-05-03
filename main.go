package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type step int

const (
	stepProjectName step = iota
	stepDjangoVersion
	stepFeatures
	stepSetup
	stepCreateApp
	stepServerOption
)

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

	inputForm     *huh.Form
	versionForm   *huh.Form
	featureForm   *huh.Form
	appForm       *huh.Form
	serverForm    *huh.Form
	appName       string
	runServer     bool
	stepMessages  []string // Store each step's message
}

func NewModel() *Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	m := &Model{
		spinner:  s,
		progress: p,
		doneChan: make(chan bool),
		step:     stepProjectName,
		features: []string{"vanilla"}, // Initialize features with a default value
		runServer: false,
	}

	theme := huh.ThemeBase()
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))

	m.inputForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project name").
				Description("How would you like to name your project?").
				Placeholder("my_django_project").
				Value(&m.projectName).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("Project name cannot be empty")
					}
					// Add validation for valid Python package name
					return nil
				}),
		),
	).WithTheme(theme).WithShowHelp(true).WithShowErrors(true)

	m.versionForm = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Django version").
				Description("Which Django version would you like to use?").
				Placeholder("5.2.0").
				Value(&m.djangoVersion),
		),
	).WithTheme(theme).WithShowHelp(true)

	m.featureForm = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Setup Type").
				Description("Choose your Django setup type").
				Options(
					huh.NewOption("Vanilla Setup 🍦", "vanilla"),
				).
				Value(&m.features[0]), // Access safely after ensuring slice is not empty
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
	return m.inputForm.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.done {
		return m, tea.Quit
	}

	// Add more keyboard shortcuts
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+b":
			// Go back to previous step if possible
			if m.step > stepProjectName {
				m.step--
				// Handle going back logic...
				return m, nil
			}
		}
	}

	switch msg := msg.(type) {
	case progressMsg:
		if float64(msg) >= 1.0 {
			m.progress.SetPercent(1.0)
			m.step = stepCreateApp
			return m, m.appForm.Init()
		}
		m.progress.SetPercent(float64(msg))
		return m, m.updateProgress()
	}

	switch m.step {
	case stepProjectName:
		form, cmd := m.inputForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.inputForm = f
			if f.State == huh.StateCompleted {
				m.step = stepDjangoVersion
				m.stepMessages = append(m.stepMessages, "Project name selected: "+m.projectName)
				return m, m.versionForm.Init()
			}
			return m, cmd
		}
	case stepDjangoVersion:
		form, cmd := m.versionForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.versionForm = f
			if f.State == huh.StateCompleted {
				m.step = stepFeatures
				m.stepMessages = append(m.stepMessages, "Django version selected: "+m.djangoVersion)
				return m, m.featureForm.Init()
			}
			return m, cmd
		}
	case stepFeatures:
		form, cmd := m.featureForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.featureForm = f
			if f.State == huh.StateCompleted {
				m.step = stepSetup
				m.stepMessages = append(m.stepMessages, "Features selected: "+fmt.Sprint(m.features))
				go m.createProject() // Run setup in background
				return m, tea.Batch(m.spinner.Tick, m.updateProgress()) // Start the spinner and progress update
			}
			return m, cmd
		}
	case stepCreateApp:
		form, cmd := m.appForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.appForm = f
			if f.State == huh.StateCompleted {
				m.step = stepServerOption
				if m.appName != "" {
					// Get absolute path to project
					wd, err := os.Getwd()
					if err != nil {
						m.error = fmt.Errorf("failed to get working directory: %v", err)
						return m, nil
					}
					projectPath := filepath.Join(wd, m.projectName)
					pythonPath := getPythonPath(projectPath)

					// Verify Python path exists
					if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
						m.error = fmt.Errorf("Python executable not found at %s: %v", pythonPath, err)
						return m, nil
					}

					// Check if manage.py exists
					managePyPath := filepath.Join(projectPath, "manage.py")
					if _, err := os.Stat(managePyPath); os.IsNotExist(err) {
						m.error = fmt.Errorf("manage.py not found in project directory. Project setup may be incomplete.")
						return m, nil
					}

					// Create app with better error handling
					createAppCmd := exec.Command(pythonPath, "manage.py", "startapp", m.appName)
					createAppCmd.Dir = projectPath

					// Capture command output for better error reporting
					output, err := createAppCmd.CombinedOutput()
					if err != nil {
						m.error = fmt.Errorf("failed to create app: %v\nOutput: %s", err, output)
						return m, nil
					 }

					// Register app in settings.py
					settingsPath := filepath.Join(projectPath, m.projectName, "settings.py")
					settingsContent, err := os.ReadFile(settingsPath)
					if err != nil {
						m.error = fmt.Errorf("failed to read settings.py: %v", err)
						return m, nil
					}

					// Find INSTALLED_APPS section and add the new app
					settingsStr := string(settingsContent)
					installedAppsIndex := strings.Index(settingsStr, "INSTALLED_APPS = [")
					if installedAppsIndex == -1 {
						m.error = fmt.Errorf("could not find INSTALLED_APPS in settings.py")
						return m, nil
					}

					// Find the closing bracket of INSTALLED_APPS
					closeBracketIndex := strings.Index(settingsStr[installedAppsIndex:], "]")
					if closeBracketIndex == -1 {
						m.error = fmt.Errorf("malformed INSTALLED_APPS in settings.py")
						return m, nil
					}

					// Insert the new app
					newSettingsContent := settingsStr[:installedAppsIndex+closeBracketIndex] +
						"    '" + m.appName + "',\n" +
						settingsStr[installedAppsIndex+closeBracketIndex:]

					if err := os.WriteFile(settingsPath, []byte(newSettingsContent), 0644); err != nil {
						m.error = fmt.Errorf("failed to update settings.py: %v", err)
						return m, nil
					}

					m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Created and registered Django app: %s", m.appName))
				}
				return m, m.serverForm.Init()
			}
			return m, cmd
		}
	case stepServerOption:
		form, cmd := m.serverForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.serverForm = f
			if f.State == huh.StateCompleted {
				if m.runServer {
					wd, err := os.Getwd()
					if err != nil {
						m.error = fmt.Errorf("failed to get working directory: %v", err)
						return m, nil
					}
					projectPath := filepath.Join(wd, m.projectName)
					pythonPath := getPythonPath(projectPath)

					// Check if Python executable exists
					if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
						m.error = fmt.Errorf("Python executable not found at %s: %v", pythonPath, err)
						return m, nil
					}

					cmd := exec.Command(pythonPath, "manage.py", "runserver")
					cmd.Dir = projectPath
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Start(); err != nil {
						m.error = fmt.Errorf("failed to start development server: %v", err)
						return m, nil
					}
					m.stepMessages = append(m.stepMessages, "✨ Development server started at http://127.0.0.1:8000/")
				}
				m.done = true
				return m, nil
			}
			return m, cmd
		}
	}

	return m, nil
}

func (m *Model) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).PaddingLeft(2)
	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).PaddingLeft(4)
	var view string

	// Add a welcome banner for the initial screen
	if m.step == stepProjectName {
		bannerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
		view += bannerStyle.Render(`
██████╗░░█████╗░██████╗░██╗██████╗░    ██████╗░░░░░░██╗
██╔══██╗██╔══██╗██╔══██╗██║██╔══██╗    ██╔══██╗░░░░░██║
██████╔╝███████║██████╔╝██║██║░░██║    ██║░░██║░░░░░██║
██╔══██╗██╔══██║██╔═══╝░██║██║░░██║    ██║░░██║██╗░░██║
██║░░██║██║░░██║██║░░░░░██║██████╔╝    ██████╔╝╚█████╔╝
╚═╝░░╚═╝╚═╝░░╚═╝╚═╝░░░░░╚═╝╚═════╝░    ╚═════╝░░╚════╝░
        `) + "\n\n"
		view += style.Render("✨ Welcome to Django Project Creator!") + "\n"
		view += msgStyle.Render("Let's build your Django project together.") + "\n\n"
	}

	// Display step messages with styling
	for _, msg := range m.stepMessages {
		view += msgStyle.Render(msg) + "\n"
	}

	// Progress indicator
	progressBar := ""
	totalSteps := 5 // Total number of possible steps
	for i := 0; i < totalSteps; i++ {
		if int(m.step) > i {
			// Completed step
			progressBar += lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("● ")
		} else if int(m.step) == i {
			// Current step
			progressBar += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("○ ")
		} else {
			// Future step
			progressBar += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("○ ")
		}
	}

	// Add this progress indicator before showing the current step's content
	if !m.done && m.error == nil {
		view += lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(progressBar) + "\n"
	}

	// Handle completion state
	if m.done {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).PaddingLeft(2)
		summaryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).PaddingLeft(4)
		houstonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

		// Add Houston face and rocket art
		houston := houstonStyle.Render(`
   ___
  /|_|\
 /_/_\_\
 \_\_/_/
  \|_|/

  ^---^
 /o   o\
 \  ω  /   Houston: Good luck out there, Djangonaut!
  \___/    Your project is ready for takeoff! 🚀
`)

		view += successStyle.Render("✅ Django project setup complete!") + "\n\n"
		view += houston + "\n\n"
		view += successStyle.Render("Project Summary:") + "\n"
		view += summaryStyle.Render(fmt.Sprintf("📁 Project: %s", m.projectName)) + "\n"
		view += summaryStyle.Render(fmt.Sprintf("📦 Django: %s", m.djangoVersion)) + "\n"

		if m.appName != "" {
			view += summaryStyle.Render(fmt.Sprintf("🧩 App: %s", m.appName)) + "\n"
		}

		// Add next steps guidance with rocket emojis
		view += "\n" + successStyle.Render("Next steps:") + "\n"
		view += summaryStyle.Render(fmt.Sprintf("1. 🚀 cd %s", m.projectName)) + "\n"
		view += summaryStyle.Render("2. 🚀 python manage.py runserver") + "\n"

		return view
	}

	// Handle error state
	if m.error != nil {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).PaddingLeft(2)
		view += errorStyle.Render("❌ " + m.error.Error())
		return view
	}

	// Handle ongoing steps
	switch m.step {
	case stepProjectName:
		view += style.Render("✨ Django Project Creator - Step 1/3\n\n") + m.inputForm.View()
	case stepDjangoVersion:
		view += style.Render("✨ Django Project Creator - Step 2/3\n\n") + m.versionForm.View()
	case stepFeatures:
		view += style.Render("✨ Django Project Creator - Step 3/3\n\n") + m.featureForm.View()
	case stepSetup:
		view += style.Render("🚧 Setting up your project...") + "\n" +
			m.spinner.View() + "\n" +
			m.progress.View()
	case stepCreateApp:
		view += style.Render("✨ Django Project Creator - Optional Step\n\n") + m.appForm.View()
	case stepServerOption:
		view += style.Render("✨ Almost done! Would you like to run the development server?") + "\n\n"
		view += m.serverForm.View()
	default:
		view += style.Render("✨ Unknown state")
	}

	// Show quit option
	quitStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).PaddingLeft(2)
	view += quitStyle.Render("\nPress 'q' to quit at any time.")

	return view
}

type progressMsg float64

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

func (m *Model) createProject() {
	if m.projectName == "" {
		m.error = fmt.Errorf("project name cannot be empty")
		return
	}

	// Project path
	projectPath := m.projectName
	if !filepath.IsAbs(projectPath) {
		wd, err := os.Getwd()
		if err != nil {
			m.error = fmt.Errorf("failed to get working directory: %v", err)
			return
		}
		projectPath = filepath.Join(wd, m.projectName)
	}

	// Create project directory
	m.progressStatus = "Creating project directory..."
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		m.error = fmt.Errorf("failed to create project directory: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "📁 Project directory created")

	// Create virtual environment
	m.progressStatus = "Setting up Python virtual environment..."
	cmd := exec.Command("uv", "venv", ".venv")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to create virtual environment: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "🐍 Virtual environment ready")

	// Install Django
	version := m.djangoVersion
	if version == "" {
		version = "5.2.0"
	}
	m.progressStatus = fmt.Sprintf("📦 Installing Django %s...", version)
	cmd = exec.Command("uv", "pip", "install", "django=="+version)
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to install Django: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("✅ Django %s installed", version))

	// Create Django project
	m.progressStatus = "Creating Django project..."
	pythonPath := getPythonPath(projectPath)

	// Check if Python executable exists
	if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
		m.error = fmt.Errorf("Python executable not found at %s: %v", pythonPath, err)
		return
	}

	// Use django-admin startproject with the correct structure
	cmd = exec.Command(pythonPath, "-m", "django", "startproject", m.projectName, ".")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		m.error = fmt.Errorf("failed to create Django project: %v", err)
		return
	}
	m.stepMessages = append(m.stepMessages, "Django project created.")

	// Verify manage.py exists
	managePyPath := filepath.Join(projectPath, "manage.py")
	if _, err := os.Stat(managePyPath); os.IsNotExist(err) {
		m.error = fmt.Errorf("Django project structure invalid: manage.py not found")
		return
	}

	// Using vanilla setup by default
	m.stepMessages = append(m.stepMessages, "Using vanilla Django setup")

	// Mark the project as setup completed
	m.stepMessages = append(m.stepMessages, "✅ Project setup finished!")
	m.doneChan <- true
}

// Helper function to get the correct Python path based on OS
func getPythonPath(projectPath string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(projectPath, ".venv", "Scripts", "python.exe")
	}
	return filepath.Join(projectPath, ".venv", "bin", "python")
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
