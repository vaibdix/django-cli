package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m *Model) View() string {
	// ===== COLOR PALETTE =====
	primaryColor := lipgloss.Color("#7D56F4")     // Rich purple
	secondaryColor := lipgloss.Color("#04B575")   // Emerald green
	accentColor := lipgloss.Color("#F25D94")      // Pink
	bgColor := lipgloss.Color("#1A1B26")          // Deep navy
	textColor := lipgloss.Color("#FAFAFA")        // Soft white
	headingColor := lipgloss.Color("#FFC44C")     // Warm yellow
	successColor := lipgloss.Color("#73F991")     // Bright green
	errorColor := lipgloss.Color("#F03E4D")       // Vibrant red
	infoColor := lipgloss.Color("#61AFEF")        // Sky blue
	subtleColor := lipgloss.Color("#6C7086")      // Muted slate

	// ===== BASE STYLES =====
	// Base style for the entire view
	baseStyle := lipgloss.NewStyle().
		Width(100).
		Padding(1).
		Background(bgColor).
		Foreground(textColor)

	// Border style with gradient effect
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Margin(1, 2).
		Background(bgColor)

	// Header style with gradient appearance
	headerStyle := lipgloss.NewStyle().
		Foreground(headingColor).
		Bold(true).
		Italic(false).
		PaddingLeft(1).
		PaddingRight(1).
		MarginBottom(1).
		BorderStyle(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}).
		BorderForeground(primaryColor).
		Align(lipgloss.Center)

	// Message style
	msgStyle := lipgloss.NewStyle().
		Foreground(infoColor).
		MarginBottom(1).
		Italic(true)

	// Quit style
	quitStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Bold(true).
		PaddingLeft(2).
		MarginTop(2)

	// Success message style
	successTextStyle := lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		MarginTop(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(successColor).
		Padding(0, 1)

	// Error message style
	errorTextStyle := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		MarginTop(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(errorColor).
		Padding(0, 1)

	// Step indicator style
	stepIndicatorStyle := lipgloss.NewStyle().
		Foreground(accentColor).
		Background(bgColor).
		Bold(true).
		PaddingLeft(1).
		PaddingRight(1).
		MarginBottom(1)

	// Houston style
	houstonStyle := lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true)

	// Splash screen styles
	astroStyle := lipgloss.NewStyle().

		Foreground(textColor).
		Bold(true).
		Padding(0, 1)

	houstonStyleSplash := lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true)

	versionStyle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Bold(true)

	rocketStyle := lipgloss.NewStyle().
		Foreground(accentColor)

	countdownStyle := lipgloss.NewStyle().
		MarginLeft(2).
		Foreground(secondaryColor)

	// Progress style
	progressStyle := lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1)
		
	// Form container style
	formContainerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(80)

	// ===== ROCKET ART =====
	rocketArt := `
  ╭─────╮ 
  │ ◠ ◡ ◠ 
  ╰─────╯ 
`
	// ===== CONTENT ASSEMBLY =====
	var contentForBorder strings.Builder

	// Add step messages with styling
	for _, msg := range m.stepMessages {
		contentForBorder.WriteString(msgStyle.Render(msg) + "\n")
	}


	var currentStepPrimaryContent string
	if m.done {
		currentStepPrimaryContent = successTextStyle.Render("✅ Django project setup complete!") + "\n\n" + 
			houstonStyle.Render(`
  ╭────────────────────────────────────────────────╮ 
  │   ◠◡◠   Houston: Good luck out there!          │
  ╰────────────────────────────────────────────────╯ 
`) + "\n\n" +
		lipgloss.NewStyle().Foreground(successColor).Render(``)
	} else if m.error != nil {
		currentStepPrimaryContent = errorTextStyle.Render("❌ " + m.error.Error()) + "\n\n" +
			lipgloss.NewStyle().Foreground(subtleColor).Render("Try again or press 'q' to quit.")
	} else {
		switch m.step {
		case stepSplashScreen:
			// Create a stylish splash screen
			countdownText := fmt.Sprintf("Starting Django project creator in %d seconds...", m.splashCountdown)
			
			// Top border for the splash container
			topBorder := lipgloss.NewStyle().Foreground(accentColor).Render("╭" + strings.Repeat("─", 70) + "╮")
			// Bottom border for the splash container
			bottomBorder := lipgloss.NewStyle().Foreground(accentColor).Render("╰" + strings.Repeat("─", 70) + "╯")


			// Create horizontal content with rocket animation and countdown
			horizontalContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				rocketStyle.Render(rocketArt),

				countdownStyle.Render(countdownText),
			)
			
			// Create the splash art
			splashArt := topBorder + "\n" +
				lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).Render(

					astroStyle.Render(" Django Launch Sequence Initiated ") + "\n\n" +
					houstonStyleSplash.Render("Houston:") + "\n" +
					"Welcome to " + astroStyle.Render("Django") + " " + versionStyle.Render("v5.2.1") + "!\n\n" +
					horizontalContent + "\n",
				) +
				bottomBorder
				
			currentStepPrimaryContent = splashArt
		case stepProjectName:
			stepIndicator := stepIndicatorStyle.Render("STEP 1/4: PROJECT SETUP")
			title := headerStyle.Render("Django Project Creator")
			formContainer := formContainerStyle.Render(m.inputForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepDjangoVersion:
			stepIndicator := stepIndicatorStyle.Render("STEP 2/4: DJANGO VERSION")
			title := headerStyle.Render("Choose Django Version")
			formContainer := formContainerStyle.Render(m.versionForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepFeatures:
			stepIndicator := stepIndicatorStyle.Render("STEP 3/4: FEATURES")
			title := headerStyle.Render("Select Project Features")
			formContainer := formContainerStyle.Render(m.featureForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepTemplates:
			stepIndicator := stepIndicatorStyle.Render("STEP 4/4: TEMPLATES")
			title := headerStyle.Render("Configure Templates")
			formContainer := formContainerStyle.Render(m.templateForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepSetup:
			title := headerStyle.Render("🚧 Project Initialization")
			
			// Create a decorated progress view
			progressView := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(infoColor).
				Padding(1, 2).
				Render(
					"Building your Django project...\n\n" +
					m.spinner.View() + " " + 
					lipgloss.NewStyle().Foreground(textColor).Render("Setting up project structure") + 
					"\n\n" +
					progressStyle.Render(m.progress.View()),
				)
				
			currentStepPrimaryContent = title + "\n\n" + progressView
		case stepCreateApp:
			stepIndicator := stepIndicatorStyle.Render("OPTIONAL: CREATE APP")
			title := headerStyle.Render("Django App Creator")
			formContainer := formContainerStyle.Render(m.appForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepAppTemplates:
			stepIndicator := stepIndicatorStyle.Render("OPTIONAL: APP TEMPLATES")
			title := headerStyle.Render("Configure App Templates")
			formContainer := formContainerStyle.Render(m.appTemplateForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepServerOption:
			stepIndicator := stepIndicatorStyle.Render("OPTIONAL: DEV SERVER")
			title := headerStyle.Render("Development Server")
			formContainer := formContainerStyle.Render(m.serverForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		case stepGitInit:
			stepIndicator := stepIndicatorStyle.Render("OPTIONAL: GIT REPOSITORY")
			title := headerStyle.Render("Initialize Git Repository")
			formContainer := formContainerStyle.Render(m.gitForm.View())
			
			currentStepPrimaryContent = stepIndicator + "\n" + title + "\n\n" + formContainer
		default:
			currentStepPrimaryContent = headerStyle.Render("✨ Unknown state")
		}
	}

	// Join the content and apply the base style
	contentForBorder.WriteString(currentStepPrimaryContent)
	
	// Footer with helpful commands
	helpFooter := ""
	if !m.done && m.error == nil && m.step != stepSplashScreen && m.step != stepSetup {
		helpFooter = "\n\n" + lipgloss.NewStyle().
			Foreground(subtleColor).
			Render("Navigate: [Tab] next field  [Shift+Tab] previous field  [Enter] confirm")
	}
	
	contentForBorder.WriteString(helpFooter)
	
	borderedContent := baseStyle.Render(contentForBorder.String())
	borderedView := borderStyle.Render(borderedContent)
	
	// Create a fancy quit indicator
	quitIndicator := quitStyle.Render("┃ Press 'q' to quit at any time ┃")
	
	// Add a footer with version info and quit message
	footer := lipgloss.JoinHorizontal(
		lipgloss.Center,
		quitIndicator,
		lipgloss.NewStyle().Width(20).Render(""),
	)

	finalView := borderedView + "\n" + footer

	return finalView
}