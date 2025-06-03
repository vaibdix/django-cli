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
    // Clean title without background - vibrant gradient-like colors
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.AdaptiveColor{
            Light: "#FF6B35", // Vibrant orange
            Dark:  "#FF8A50", // Brighter orange for dark mode
        }).
        MarginBottom(1)

    // Sophisticated subtitle with better contrast
    subtitleStyle = lipgloss.NewStyle().
        Italic(true).
        Foreground(lipgloss.AdaptiveColor{
            Light: "#7C3AED", // Rich purple
            Dark:  "#A78BFA", // Lighter purple for dark mode
        }).
        MarginBottom(1)

    // Clean error style without background - just colored text
    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#DC2626", // Red
            Dark:  "#EF4444", // Lighter red for dark mode
        }).
        Bold(true).
        MarginBottom(1)

    // Clean success style without background
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#059669", // Green
            Dark:  "#10B981", // Lighter green for dark mode
        }).
        Bold(true).
        MarginBottom(1)

    // Clean warning style without background
    warningStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#D97706", // Amber
            Dark:  "#F59E0B", // Lighter amber for dark mode
        }).
        Bold(true).
        MarginBottom(1)

    // Refined footer with better readability
    footerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#64748B", // Sophisticated gray
            Dark:  "#CBD5E1", // Light gray for dark mode
        }).
        Italic(true).
        MarginTop(1)

    // Enhanced content box with better contrast
    contentBox = lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.AdaptiveColor{
            Light: "#7C3AED", // Rich purple
            Dark:  "#A855F7", // Bright purple
        })

    // Elegant highlight boxes without jarring backgrounds
    highlightBox = lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.AdaptiveColor{
            Light: "#06B6D4", // Cyan accent
            Dark:  "#22D3EE", // Bright cyan
        })

    accentBox = lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.AdaptiveColor{
            Light: "#EC4899", // Pink accent
            Dark:  "#F472B6", // Light pink
        })

    // Clean help text
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#6B7280", // Muted gray
            Dark:  "#9CA3AF", // Light muted gray
        }).Render

    // Interactive elements with yellow cursor
    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#000000")).  // Black text for contrast
        Background(lipgloss.AdaptiveColor{
            Light: "#FDE047", // Bright yellow
            Dark:  "#FACC15", // Golden yellow for dark mode
        }).
        Bold(true).
        Padding(0, 1)

    unselectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#4B5563", // Dark gray
            Dark:  "#E5E7EB", // Light gray
        }).
        Padding(0, 1)

    // Colorful list item styles
    listItemStyle1 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#059669", // Emerald
            Dark:  "#10B981", // Light emerald
        }).
        Bold(true)

    listItemStyle2 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#DC2626", // Red
            Dark:  "#EF4444", // Light red
        }).
        Bold(true)

    listItemStyle3 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#7C3AED", // Purple
            Dark:  "#A78BFA", // Light purple
        }).
        Bold(true)

    listItemStyle4 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#EA580C", // Orange
            Dark:  "#FB923C", // Light orange
        }).
        Bold(true)

    listItemStyle5 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#0891B2", // Cyan
            Dark:  "#22D3EE", // Light cyan
        }).
        Bold(true)

    listItemStyle6 = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#BE185D", // Pink
            Dark:  "#F472B6", // Light pink
        }).
        Bold(true)

    // Selected list item with yellow background
    selectedListItemStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#000000")).  // Black text for contrast
        Background(lipgloss.AdaptiveColor{
            Light: "#FDE047", // Bright yellow
            Dark:  "#FACC15", // Golden yellow for dark mode
        }).
        Bold(true).
        Padding(0, 1)

    // Status indicators with more sophistication
    activeStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#059669", // Rich green
            Dark:  "#34D399", // Bright green
        }).
        Bold(true)

    inactiveStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#9CA3AF", // Muted gray
            Dark:  "#6B7280", // Dark muted gray
        })

    // Enhanced accent palette
    primaryAccent   = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#A78BFA"} // Rich purple
    secondaryAccent = lipgloss.AdaptiveColor{Light: "#06B6D4", Dark: "#22D3EE"} // Vibrant cyan
    tertiaryAccent  = lipgloss.AdaptiveColor{Light: "#EC4899", Dark: "#F472B6"} // Electric pink
    
    // Special text styles for form labels and content
    labelStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#374151", // Dark gray
            Dark:  "#F3F4F6", // Light gray
        }).
        Bold(true)

    inputPromptStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#7C3AED", // Rich purple
            Dark:  "#A78BFA", // Light purple
        })

    // Enhanced progress and status styles
    progressStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#059669", // Rich green
            Dark:  "#10B981", // Emerald
        })

    // Text styles without backgrounds
    instructionTextStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#6B7280", // Gray
            Dark:  "#9CA3AF", // Light gray
        }).
        Italic(true)

    // Clean section header style
    sectionHeaderStyle = lipgloss.NewStyle().
        Foreground(lipgloss.AdaptiveColor{
            Light: "#EC4899", // Pink
            Dark:  "#F472B6", // Light pink
        }).
        Bold(true).
        MarginBottom(1)
)

// Helper function to get colorful list item style based on index
func getListItemStyle(index int, isSelected bool) lipgloss.Style {
    if isSelected {
        return selectedListItemStyle
    }
    
    styles := []lipgloss.Style{
        listItemStyle1, // Emerald
        listItemStyle2, // Red
        listItemStyle3, // Purple
        listItemStyle4, // Orange
        listItemStyle5, // Cyan
        listItemStyle6, // Pink
    }
    
    return styles[index%len(styles)]
}

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
		if m.startDevServer {
			s.WriteString(titleStyle.Render("🚀 Development Environment Started!") + "\n\n")
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
			s.WriteString(titleStyle.Render("✅ Django Project Setup Complete!") + "\n\n")
			s.WriteString(subtitleStyle.Render("What's Next:") + "\n")
			s.WriteString(fmt.Sprintf("1. Navigate to your project directory:\n   cd %s\n\n", m.projectName))

			projectAbsPath, _ := filepath.Abs(m.projectName)
			pythonVenvPath := getPythonPath(projectAbsPath)
			s.WriteString(fmt.Sprintf("2. Start the development server:\n   %s manage.py runserver\n\n", pythonVenvPath))
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
		
		// Show percentage with enhanced styling
		percentage := int(m.progress.Percent() * 100)
		s.WriteString(pad + progressStyle.Render(fmt.Sprintf("Progress: %d%%", percentage)) + "\n")

	case stepProjectName:
		if activeForm != nil {
			s.WriteString(titleStyle.Render("🚀 Django Project Configuration") + "\n")
			s.WriteString(subtitleStyle.Render("Configure your Django project with all options in one place") + "\n\n")
			s.WriteString(activeForm.View())
		}
	
	case stepDevServerPrompt:
		if activeForm != nil {
			s.WriteString(titleStyle.Render("🎉 Project Setup Complete!") + "\n\n")
			s.WriteString(activeForm.View())
		}
	
	case stepComplete:
		s.WriteString(titleStyle.Render("✅ All Done!") + "\n\n")
		s.WriteString("   ╭─────╮\n")
		s.WriteString("   │ ◠ ◡ ◠           happy coding\n")
		s.WriteString("   ╰─────╯\n\n")
		s.WriteString(subtitleStyle.Render("Manual Steps:") + "\n")
		s.WriteString(fmt.Sprintf("1. Navigate to your project: cd %s\n", m.projectName))
		projectAbsPath, _ := filepath.Abs(m.projectName)
		pythonVenvPath := getPythonPath(projectAbsPath)
		if m.setupTailwind {
			s.WriteString("2. Start CSS watching: npm run watch:css\n")
			s.WriteString(fmt.Sprintf("3. Start Django server: %s manage.py runserver\n", pythonVenvPath))
		} else {
			s.WriteString(fmt.Sprintf("2. Start Django server: %s manage.py runserver\n", pythonVenvPath))
		}
		s.WriteString("\n" + footerStyle.Render("Press Enter or Q to exit."))
		return baseStyle.Render(s.String())
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

// Helper function to render colorful list items with proper form integration
func (m *Model) renderColorfulListItem(index int, text string, isChecked bool, isSelected bool) string {
	style := getListItemStyle(index, isSelected)
	
	var checkbox string
	if isChecked {
		checkbox = "[•]"
	} else {
		checkbox = "[ ]"
	}
	
	var prefix string
	if isSelected {
		prefix = "> "
	} else {
		prefix = "  "
	}
	
	return prefix + style.Render(checkbox + " " + text)
}

// Example integration - replace this with your actual form logic
func (m *Model) renderFormWithColorfulItems() string {
	var s strings.Builder
	
	// You'll need to replace these with your actual form data
	// Example of how to integrate with your real form:
	/*
	form := m.getActiveForm() // Your actual form
	for i, item := range form.Items {
		isSelected := i == form.cursor
		isChecked := form.selectedItems[i]
		
		colorfulItem := m.renderColorfulListItem(i, item.text, isChecked, isSelected)
		s.WriteString(colorfulItem + "\n")
	}
	*/
	
	return s.String()
}

