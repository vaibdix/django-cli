package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
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

	if m.error != nil {
		errMsg := ErrorStyle.Render(fmt.Sprintf("❌ ERROR: %s", m.error.Error()))
		s.WriteString(errMsg + "\n\n")
		s.WriteString("Press Enter or Q to exit.")
		return baseStyle.Render(s.String())
	}

	if m.done {
		if m.startDevServer {
			s.WriteString(TitleStyle.Render("🚀 Development Environment Started!") + "\n\n")
			s.WriteString("VS Code has been opened and the development server is starting...\n\n")
			if m.setupTailwind {
				s.WriteString("Two terminals have been opened:\n")
				s.WriteString("• Terminal 1: Tailwind CSS watcher (npm run watch:css)\n")
				s.WriteString("• Terminal 2: Django development server\n\n")
			} else {
				s.WriteString("Django development server terminal has been opened.\n\n")
			}
			s.WriteString("   ╭─────╮\n")
			s.WriteString("   │ ◠ ◡ ◠        happy coding 🚀 \n")
			s.WriteString("   ╰─────╯\n")
		} else {
			s.WriteString(TitleStyle.Render("✅ Django Project Setup Complete!") + "\n\n")
			s.WriteString(SubtitleStyle.Render("What's Next:") + "\n")
			s.WriteString(fmt.Sprintf("1. Navigate to your project directory:\n   cd %s\n\n", m.projectName))

			projectAbsPath, _ := filepath.Abs(m.projectName)
			pythonVenvPath := getPythonPath(projectAbsPath)
			s.WriteString(fmt.Sprintf("2. Start the development server:\n   %s manage.py runserver\n\n", pythonVenvPath))
		}

		s.WriteString("\n" + FooterStyle.Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
	}

	if m.step == stepSplashScreen {
		djangoDisplayVersion := m.djangoVersion
		if djangoDisplayVersion == "" || djangoDisplayVersion == "latest" {
			djangoDisplayVersion = "latest stable"
		}
		s.WriteString(TitleStyle.Render("Django Forge CLI 🚀") + "\n\n")
		s.WriteString(fmt.Sprintf("Welcome! Initializing Django project creator with Django %s.\n", djangoDisplayVersion))
		s.WriteString(fmt.Sprintf("Starting in %d seconds...\n\n", m.splashCountdown))
		s.WriteString(SubtitleStyle.Render("Crafting your Django project, one step at a time."))
		return baseStyle.Render(s.String())
	}

	activeForm := m.getActiveForm()

	switch m.step {
	case stepSetup:
		s.WriteString(TitleStyle.Render("🚧 Project Initialization 🚧") + "\n\n")
		s.WriteString(fmt.Sprintf("%s %s\n\n", m.spinner.View(), m.progressStatus))

		m.progress.Width = contentWidth - 8
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}

		pad := strings.Repeat(" ", padding)
		s.WriteString(pad + m.progress.View() + "\n\n")

		percentage := int(m.progress.Percent() * 100)
		s.WriteString(pad + ProgressStyle.Render(fmt.Sprintf("Progress: %d%%", percentage)) + "\n")

		if len(m.stepMessages) > 0 {
			s.WriteString("\n" + pad + "Recent steps:\n")
			start := len(m.stepMessages)
			if start > 3 {
				start = len(m.stepMessages) - 3
			}
			for _, msg := range m.stepMessages[start:] {
				s.WriteString(pad + "• " + msg + "\n")
			}
		}

		return ContentBox.Width(contentWidth).Render(s.String())

	case stepProjectName:
		if activeForm != nil {
			s.WriteString(TitleStyle.Render("🚀 Django Project Configuration") + "\n")
			s.WriteString(activeForm.View())
		}

	case stepDjangoVersion:
		if activeForm != nil {
			s.WriteString(TitleStyle.Render("🐍 Django Version Selection") + "\n")
			s.WriteString(activeForm.View())
		}

	case stepProjectConfig:
		if activeForm != nil {
			s.WriteString(TitleStyle.Render("⚙️ Project Features") + "\n")
			s.WriteString(activeForm.View())
		}

	case stepAppName:
		if activeForm != nil {
			s.WriteString(TitleStyle.Render("📱 App Configuration") + "\n")
			s.WriteString(activeForm.View())
		}

	case stepDevServerPrompt:
		if activeForm != nil {
			s.WriteString(TitleStyle.Render("🎉 Project Setup Complete!") + "\n\n")
			s.WriteString(activeForm.View())
		}

	case stepComplete:
		s.WriteString(TitleStyle.Render("✅ All Done!") + "\n\n")
		s.WriteString("   ╭─────╮\n")
		s.WriteString("   │ ◠ ◡ ◠           happy coding\n")
		s.WriteString("   ╰─────╯\n\n")
		s.WriteString(SubtitleStyle.Render("Manual Steps:") + "\n")
		s.WriteString(fmt.Sprintf("1. Navigate to your project: cd %s\n", m.projectName))
		projectAbsPath, _ := filepath.Abs(m.projectName)
		pythonVenvPath := getPythonPath(projectAbsPath)
		if m.setupTailwind {
			s.WriteString("2. Start CSS watching: npm run watch:css\n")
			s.WriteString(fmt.Sprintf("3. Start Django server: %s manage.py runserver\n", pythonVenvPath))
		} else {
			s.WriteString(fmt.Sprintf("2. Start Django server: %s manage.py runserver\n", pythonVenvPath))
		}
		s.WriteString("\n" + FooterStyle.Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
	}

	quitHelp := FooterStyle.Render("Press 'q' or 'Ctrl+C' to quit.")
	var navHelp string

	s.WriteString("\n")
	if navHelp != "" {
		s.WriteString(navHelp + "\n")
	}
	s.WriteString(quitHelp)

	return ContentBox.Width(contentWidth).Render(s.String())
}
