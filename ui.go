package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#A0A0A0")).
			MarginBottom(1)

	errorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#FF0000")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(1, 2).
			MarginBottom(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			MarginTop(1)

	contentBox = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#1E1E1E"))
)

func (m *Model) View() string {
	viewWidth := m.width
	if viewWidth <= 0 {
		viewWidth = 100
	}
	contentWidth := viewWidth - 8
	if contentWidth < 70 {
		contentWidth = 70
	}

	baseStyle := lipgloss.NewStyle().
		Width(contentWidth).
		Padding(1, 4)

	var s strings.Builder

	// Error state
	if m.error != nil {
		errMsg := errorStyle.Render(fmt.Sprintf("❌ ERROR: %s", m.error.Error()))
		s.WriteString(errMsg + "\n\n")
		s.WriteString("Press Enter or Q to exit.")
		return baseStyle.Render(s.String())
	}

	// Done state
	if m.done {
		s.WriteString(titleStyle.Render("✅ Django Project Setup Complete!") + "\n\n")
		s.WriteString(subtitleStyle.Render("What's Next:") + "\n")
		s.WriteString(fmt.Sprintf("1. Navigate to your project directory:\n   cd %s\n\n", m.projectName))

		projectAbsPath, _ := filepath.Abs(m.projectName)
		pythonVenvPath := getPythonPath(projectAbsPath)
		if !m.runServer {
			s.WriteString(fmt.Sprintf("2. Start the development server:\n   %s manage.py runserver\n\n", pythonVenvPath))
		} else {
			s.WriteString(fmt.Sprintf("2. If the server is not running, start it with:\n   %s manage.py runserver\n\n", pythonVenvPath))
		}

		s.WriteString(subtitleStyle.Render("Detailed Setup Log:") + "\n")
		for _, msg := range m.stepMessages {
			if strings.HasPrefix(msg, "✓") || strings.HasPrefix(msg, "PROGRESS") {
				s.WriteString(fmt.Sprintf("✓ %s\n", strings.TrimPrefix(msg, "✓ ")))
			} else if strings.HasPrefix(msg, "•") || strings.HasPrefix(msg, "  •") {
				s.WriteString(fmt.Sprintf("%s\n", msg))
			} else if strings.HasPrefix(msg, "To start the server:") {
				s.WriteString(fmt.Sprintf("\n%s\n", msg))
			} else {
				s.WriteString(fmt.Sprintf("• %s\n", msg))
			}
		}

		s.WriteString("\n" + footerStyle.Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
	}

	// Splash screen
	if m.step == stepSplashScreen {
		djangoDisplayVersion := m.djangoVersion
		if djangoDisplayVersion == "" || djangoDisplayVersion == "latest" {
			djangoDisplayVersion = "latest stable"
		}
		s.WriteString(titleStyle.Render("🚀 Django Forge CLI 🚀") + "\n\n")
		s.WriteString(fmt.Sprintf("Welcome! Initializing Django project creator with Django %s.\n", djangoDisplayVersion))
		s.WriteString(fmt.Sprintf("Starting in %d seconds...\n\n", m.splashCountdown))
		s.WriteString(subtitleStyle.Render("Crafting your Django project, one step at a time."))
		return baseStyle.Render(s.String())
	}

	// Main content area
	activeForm := m.getActiveForm()

	switch m.step {
	case stepSetup:
		s.WriteString(titleStyle.Render("🚧 Project Initialization 🚧") + "\n\n")
		s.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), m.progressStatus))
		m.progress.Width = contentWidth - 8
		s.WriteString(m.progress.View())

	default:
		if activeForm != nil {
			var stepTitle, stepDescription string
			switch m.step {
			case stepProjectName:
				stepTitle = "Step 1: Project Name"
				stepDescription = "Enter a memorable name for your new Django project."
			case stepDjangoVersion:
				stepTitle = "Step 2: Django Version"
				stepDescription = "Choose the Django version for your project. 'latest' is recommended."
			case stepFeatures:
				stepTitle = "Step 3: Setup Type"
				stepDescription = "Select a setup type to include common project structures."
			case stepTemplates:
				stepTitle = "Step 4: Global Templates"
				stepDescription = "Include standard global template directories (e.g., 'templates')?"
			case stepCreateApp:
				stepTitle = "Optional: Create App"
				stepDescription = "Do you want to create an initial Django app within your project?"
			case stepAppTemplates:
				stepTitle = fmt.Sprintf("Optional: App Templates for '%s'", m.appName)
				stepDescription = "Include standard template directories for your app (e.g., 'templates/<app_name>')?"
			case stepServerOption:
				stepTitle = "Optional: Development Server"
				stepDescription = "Automatically start the Django development server after setup?"
			case stepGitInit:
				stepTitle = "Optional: Git Repository"
				stepDescription = "Initialize a Git repository for version control?"
			}

			if stepTitle != "" {
				s.WriteString(titleStyle.Render(stepTitle) + "\n")
				s.WriteString(subtitleStyle.Render(stepDescription) + "\n\n")
			}
			s.WriteString(activeForm.View())
		}
	}

	quitHelp := footerStyle.Render("Press 'q' or 'Ctrl+C' to quit.")

	var navHelp string
	if activeForm != nil && m.step != stepSetup && m.step != stepSplashScreen {
		navHelp = footerStyle.Render("Navigate: ↑/↓ or Tab/Shift+Tab  |  Select: Enter  |  Change: ←/→")
	}

	s.WriteString("\n")
	if navHelp != "" {
		s.WriteString(navHelp + "\n")
	}
	s.WriteString(quitHelp)

	return contentBox.Width(contentWidth).Render(s.String())
}
