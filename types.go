package main

// step defines the current stage of the CLI interaction.
type step int

const (
	stepSplashScreen = iota
	stepProjectName
	stepDjangoVersion
	stepFeatures
	stepTemplates    // For global templates
	stepSetup        // Main project setup (venv, django install, project creation)
	stepCreateApp    // Prompting for initial app name
	stepAppTemplates // For app-specific templates, if an app is created
	stepServerOption
	stepGitInit
)

// projectProgressMsg is sent by CreateProject to update the progress bar and status.
type projectProgressMsg struct {
	percent float64
	status  string
}

// projectCreationDoneMsg is sent by CreateProject when it's finished (successfully or with error).
type projectCreationDoneMsg struct {
	err error
}

// tickMsg is used for the splash screen countdown.
type tickMsg struct{}