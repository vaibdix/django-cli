package main

type step int

const (
	stepSplashScreen = iota
	stepProjectName
	stepDjangoVersion
	stepProjectConfig
	stepAppName
	stepFeatures
	stepTemplates
	stepSetup
	stepCreateApp
	stepAppTemplates
	stepServerOption
	stepGitInit
	stepDevServerPrompt
	stepComplete
)

type projectProgressMsg struct {
	percent float64
	status  string
}
type projectCreationDoneMsg struct {
	err error
}
type tickMsg struct{}
