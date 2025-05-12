package main

type step int

const (
	stepSplashScreen = iota
	stepProjectName
	stepDjangoVersion
	stepFeatures
	stepTemplates
	stepSetup
	stepCreateApp
	stepAppTemplates    
	stepServerOption
)

type progressMsg float64