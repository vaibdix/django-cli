package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
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

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
)

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

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
		
		// Enhanced animated progress bar
		m.progress.Width = contentWidth - 8
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}
		
		pad := strings.Repeat(" ", padding)
		s.WriteString("\n" + pad + m.progress.View() + "\n\n")
		
		// Show percentage
		percentage := int(m.progress.Percent() * 100)
		s.WriteString(pad + fmt.Sprintf("Progress: %d%%\n", percentage))

	case stepProjectName:
		if activeForm != nil {
			s.WriteString(titleStyle.Render("🚀 Django Project Configuration") + "\n")
			s.WriteString(subtitleStyle.Render("Configure your Django project with all options in one place") + "\n\n")
			s.WriteString(activeForm.View())
		}
	}

	quitHelp := footerStyle.Render("Press 'q' or 'Ctrl+C' to quit.")

	var navHelp string
	if activeForm != nil && m.step == stepProjectName {
		navHelp = footerStyle.Render("Navigate: ↑/↓ or Tab/Shift+Tab  |  Select: Space/Enter  |  Submit: Enter")
	}

	s.WriteString("\n")
	if navHelp != "" {
		s.WriteString(navHelp + "\n")
	}
	s.WriteString(quitHelp)

	return contentBox.Width(contentWidth).Render(s.String())
}
