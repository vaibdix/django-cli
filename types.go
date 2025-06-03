package main
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
	stepDevServerPrompt // Ask about running development server at the end
	stepComplete     // Show completion message
)
type projectProgressMsg struct {
	percent float64
	status  string
}
type projectCreationDoneMsg struct {
	err error
}
type tickMsg struct{}
