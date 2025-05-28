package main

import (
	"fmt"
	"path/filepath" // Added import
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	// ===== COLOR PALETTE (Example) =====
	primaryColor := lipgloss.Color("#7D56F4")
	// secondaryColor := lipgloss.Color("#04B575")
	accentColor := lipgloss.Color("#F25D94")
	// bgColor := lipgloss.Color("#1A1B26")
	textColor := lipgloss.Color("#FAFAFA")
	headingColor := lipgloss.Color("#FFC44C")
	successColor := lipgloss.Color("#73F991")
	errorColor := lipgloss.Color("#F03E4D")
	infoColor := lipgloss.Color("#61AFEF")
	subtleColor := lipgloss.Color("#6C7086")

	// Dynamic width for base style
	viewWidth := m.width
	if viewWidth <= 0 {
		viewWidth = 100 // Default width if not set
	}
	// Ensure content area is slightly less than total view width for padding/margin
	contentWidth := viewWidth - 4
	if contentWidth < 50 {
		contentWidth = 50
	}

	baseStyle := lipgloss.NewStyle().
		Width(contentWidth). // Make width dynamic
		Padding(1, 2)       // Consistent padding

	headerStyle := lipgloss.NewStyle().
		Foreground(headingColor).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder(), false, false, true, false). // Bottom border
		BorderForeground(primaryColor).
		Align(lipgloss.Center).Width(contentWidth - 4)

	// Step indicator style
	stepIndicatorStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		MarginBottom(1).
		Padding(0, 1)

	var s strings.Builder

	if m.error != nil {
		errorStyled := lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Border(lipgloss.DoubleBorder(), true).
			BorderForeground(errorColor).
			Padding(1, 2).
			SetString(fmt.Sprintf("❌ ERROR: %s", m.error.Error())).String()
		s.WriteString(errorStyled)
		s.WriteString("\n\n")
		s.WriteString(lipgloss.NewStyle().Foreground(subtleColor).Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
	}

	if m.done {
		successMsg := "✅ Django project setup complete!\n\n"
		successMsg += "What's next:\n"
		successMsg += fmt.Sprintf("  1. Navigate to your project: cd %s\n", m.projectName)
		if !m.runServer { // Only show this if server wasn't "auto-started" (or rather, if we didn't give the command)
			projectAbsPath, _ := filepath.Abs(m.projectName) // Used filepath
			pythonVenvPath := getPythonPath(projectAbsPath)
			successMsg += fmt.Sprintf("  2. Start the development server: %s manage.py runserver\n", pythonVenvPath)
		} else {
			projectAbsPath, _ := filepath.Abs(m.projectName) // Used filepath
			pythonVenvPath := getPythonPath(projectAbsPath)
			successMsg += fmt.Sprintf("  2. If you chose to run the server, it might be starting or use: %s manage.py runserver\n", pythonVenvPath)
		}

		successStyled := lipgloss.NewStyle().
			Foreground(successColor).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(successColor).
			Padding(1, 2).
			SetString(successMsg).String()
		s.WriteString(successStyled)
		s.WriteString("\n\n")
		// Display step messages for review
		s.WriteString("Log:\n")
		for _, msg := range m.stepMessages {
			s.WriteString(lipgloss.NewStyle().Foreground(subtleColor).Render("  "+msg) + "\n")
		}
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(subtleColor).Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
	}

	// Splash Screen
	if m.step == stepSplashScreen {
		djangoDisplayVersion := m.djangoVersion
		if djangoDisplayVersion == "" || djangoDisplayVersion == "latest" {
			djangoDisplayVersion = "latest stable" // Default display
		}
		splashTitle := headerStyle.Render("🚀 Django Forge CLI 🚀")
		welcome := fmt.Sprintf("Welcome! Initializing Django project creator with Django %s.", djangoDisplayVersion)
		countdown := fmt.Sprintf("\nStarting in %d seconds...", m.splashCountdown)
		s.WriteString(splashTitle + "\n\n")
		s.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(welcome))
		s.WriteString(lipgloss.NewStyle().Foreground(textColor).Render(countdown))
		return baseStyle.Render(s.String())
	}

	// Main content area for forms and setup progress
	var currentStepPrimaryContent string
	activeForm := m.getActiveForm() // Get the form for the current step

	switch m.step {
	case stepSetup:
		title := headerStyle.Render("🚧 Project Initialization 🚧")
		status := lipgloss.NewStyle().Foreground(infoColor).Render(m.progressStatus)
		progressView := m.progress.View()

		s.WriteString(title + "\n\n")
		s.WriteString(m.spinner.View() + " " + status + "\n\n")
		s.WriteString(progressView)
		currentStepPrimaryContent = s.String() // Use the builder directly

	default: // For all form steps
		if activeForm != nil {
			// Render step indicator if available for the form
			var stepTitle string
			switch m.step {
			case stepProjectName:
				stepTitle = "Step 1: Project Name"
			case stepDjangoVersion:
				stepTitle = "Step 2: Django Version"
			case stepFeatures:
				stepTitle = "Step 3: Setup Type"
			case stepTemplates:
				stepTitle = "Step 4: Global Templates"
			case stepCreateApp:
				stepTitle = "Optional: Create App"
			case stepAppTemplates:
				stepTitle = fmt.Sprintf("Optional: App Templates for '%s'", m.appName)
			case stepServerOption:
				stepTitle = "Optional: Development Server"
			case stepGitInit:
				stepTitle = "Optional: Git Repository"
			}
			if stepTitle != "" {
				s.WriteString(stepIndicatorStyle.Render(stepTitle) + "\n")
			}
			s.WriteString(activeForm.View())
			currentStepPrimaryContent = s.String()
		} else {
			currentStepPrimaryContent = "Unknown step or no active form."
		}
	}

	// Footer with quit message
	quitHelp := "\n\n" + lipgloss.NewStyle().Foreground(subtleColor).Render("Press 'q' or 'Ctrl+C' to quit.")

	// For form steps, add navigation help
	if activeForm != nil && m.step != stepSetup {
		navHelp := lipgloss.NewStyle().Foreground(subtleColor).Italic(true).Render(
			"Navigate: [↑/↓] or [Tab/Shift+Tab]  |  Select/Confirm: [Enter]  |  Change field: [←/→] for some inputs",
		)
		return baseStyle.Render(currentStepPrimaryContent + "\n\n" + navHelp + quitHelp)
	}

	return baseStyle.Render(currentStepPrimaryContent + quitHelp)
}