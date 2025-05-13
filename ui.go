package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205")).
		Padding(1, 2). 
		Margin(1, 0)  
	borderStyle = borderStyle.MaxWidth(90)

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		PaddingLeft(0). 
		MarginBottom(1).
		Underline(true)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		MarginBottom(1) 

	quitStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		PaddingLeft(2).
		MarginTop(1) 

	var contentForBorder string 

	for _, msg := range m.stepMessages {
		contentForBorder += msgStyle.Render(msg) + "\n"
	}

	var currentStepPrimaryContent string
	if m.done {
		successTextStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true).
			MarginTop(1)   
		houston := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(`
  ╭─────╮ 
  │ ◠ ◡ ◠    Houston: Good luck out there, Djangonaut! 🚀 
  ╰─────╯ 
`)
		currentStepPrimaryContent = successTextStyle.Render("✅ Django project setup complete!") + "\n\n" + houston + "\n\n"
	} else if m.error != nil {
		errorTextStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			MarginTop(1)   

		currentStepPrimaryContent = errorTextStyle.Render("❌ " + m.error.Error())
	} else {
		switch m.step {
		case stepSplashScreen:
			astroStyle := lipgloss.NewStyle().Background(lipgloss.Color("10")).Foreground(lipgloss.Color("0")).Bold(true).Padding(0, 1)
			houstonStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
			versionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
			rocketStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
			countdownStyle := lipgloss.NewStyle().MarginLeft(2)
			rocketArt := `  
  ╭─────╮
  │ ◠ ◡ ◠ 
  ╰─────╯`
			countdownText := fmt.Sprintf("Starting Django project creator in %d seconds...\n", m.splashCountdown)
			horizontalContent := lipgloss.JoinHorizontal(
				lipgloss.Center,
				rocketStyle.Render(rocketArt),
				countdownStyle.Render(countdownText),
			)
			splashArt := astroStyle.Render("Django") + " Launch sequence initiated.\n\n" +
				houstonStyle.Render("Houston:") + "\n" +
				"Welcome to " + astroStyle.Render("Django") + " " + versionStyle.Render("v5.2.1") + "!\n\n" +
				horizontalContent + "\n"
			currentStepPrimaryContent = splashArt
		case stepProjectName:
			currentStepPrimaryContent = headerStyle.Render("✨ Django Project Creator - Step 1/4") + "\n\n" + m.inputForm.View()
		case stepDjangoVersion:
			currentStepPrimaryContent = headerStyle.Render("✨ Django Project Creator - Step 2/4") + "\n\n" + m.versionForm.View()
		case stepFeatures:
			currentStepPrimaryContent = headerStyle.Render("✨ Django Project Creator - Step 3/4") + "\n\n" + m.featureForm.View()
		case stepTemplates:
			currentStepPrimaryContent = headerStyle.Render("✨ Django Project Creator - Step 4/4") + "\n\n" + m.templateForm.View()
		case stepSetup:
			progressStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
			currentStepPrimaryContent = headerStyle.Render("🚧 Setting up your project...") + "\n" +
				m.spinner.View() + "\n" +
				progressStyle.Render(m.progress.View())
		case stepCreateApp:
			currentStepPrimaryContent = headerStyle.Render("✨ Django Project Creator - Optional Step") + "\n\n" + m.appForm.View()
		case stepAppTemplates:
			currentStepPrimaryContent = headerStyle.Render("✨ App Templates Setup") + "\n\n" + m.appTemplateForm.View()
		case stepServerOption:
			currentStepPrimaryContent = headerStyle.Render("✨ Run Development Server?") + "\n\n" + m.serverForm.View()
		case stepGitInit:
			currentStepPrimaryContent = headerStyle.Render("✨ Initialize Git Repository?") + "\n\n" + m.gitForm.View()
		default:
			currentStepPrimaryContent = headerStyle.Render("✨ Unknown state")
		}
	}

	contentForBorder += currentStepPrimaryContent
	borderedView := borderStyle.Render(contentForBorder)
	finalView := borderedView + "\n" + quitStyle.Render("Press 'q' to quit at any time.")

	return finalView
}