package main

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle spinner updates
	var spinnerCmd tea.Cmd
	m.spinner, spinnerCmd = m.spinner.Update(msg)
	cmds = append(cmds, spinnerCmd)

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyCtrlC || keyMsg.String() == "q" {
			if !m.done && m.error == nil {
				return m, tea.Quit
			}
		}
	}
	if m.error != nil || m.done {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyEnter || keyMsg.Type == tea.KeyCtrlC || keyMsg.String() == "q" || keyMsg.String() == "esc" {
				return m, tea.Quit
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = m.width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}
		return m, nil

	case tickMsg:
		if m.step == stepSplashScreen {
			m.splashCountdown--
			if m.splashCountdown <= 0 {
				m.step = stepProjectName
				cmds = append(cmds, m.projectNameForm.Init())
			} else {
				cmds = append(cmds, tea.Tick(1*time.Second, func(_ time.Time) tea.Msg {
					return tickMsg{}
				}))
			}
		}
		return m, tea.Batch(cmds...)

	case projectProgressMsg:
		if m.step == stepSetup {
			cmd := m.progress.SetPercent(msg.percent)
			m.progressStatus = msg.status
			m.stepMessages = append(m.stepMessages, "PROGRESS: "+msg.status)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(append(cmds, m.spinner.Tick)...)

	case projectCreationDoneMsg:
		if m.step == stepSetup {
			if msg.err != nil {
				m.error = msg.err
				m.progressStatus = "Error during project setup!"
				return m, nil
			}
			cmd := m.progress.SetPercent(1.0)
			m.progressStatus = "Django project setup complete!"
			m.stepMessages = append(m.stepMessages, "✅ Django project setup complete!")
			m.step = stepDevServerPrompt
			cmds = append(cmds, cmd, m.devServerForm.Init())
		}
		return m, tea.Batch(cmds...)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	activeForm := m.getActiveForm()
	if activeForm != nil {
		if activeForm.State != huh.StateCompleted {
			formModel, formCmd := activeForm.Update(msg)
			if castedForm, ok := formModel.(*huh.Form); ok {
				switch m.step {
				case stepProjectName:
					m.projectNameForm = castedForm
				case stepDjangoVersion:
					m.djangoVersionForm = castedForm
				case stepProjectConfig:
					m.projectConfigForm = castedForm
				case stepAppName:
					m.appNameForm = castedForm
				case stepDevServerPrompt:
					m.devServerForm = castedForm
				}
			}
			cmds = append(cmds, formCmd)
		} else {
			switch m.step {
			case stepProjectName:
				m.step = stepDjangoVersion
				cmds = append(cmds, m.djangoVersionForm.Init())
			case stepDjangoVersion:
				m.step = stepProjectConfig
				cmds = append(cmds, m.projectConfigForm.Init())
			case stepProjectConfig:
				m.processFormData()
				if m.createAppTemplates || m.setupRestFramework {
					m.step = stepAppName
					cmds = append(cmds, m.appNameForm.Init())
				} else {
					m.step = stepSetup
					m.totalSteps = m.calculateTotalSteps()
					m.progressStatus = "Starting project setup..."
					m.progress.SetPercent(0.0)
					go m.CreateProject()
					cmds = append(cmds,
						m.spinner.Tick,
						func() tea.Msg {
							return projectProgressMsg{
								percent: 0.0,
								status:  "Initializing project setup...",
							}
						},
						m.progress.SetPercent(0.0),
						func() tea.Msg {
							return progress.FrameMsg{}
						},
					)
				}
			case stepAppName:
				m.step = stepSetup
				m.totalSteps = m.calculateTotalSteps()
				m.progressStatus = "Starting project setup..."
				m.progress.SetPercent(0.0)
				go m.CreateProject()
				cmds = append(cmds,
					m.spinner.Tick,
					func() tea.Msg {
						return projectProgressMsg{
							percent: 0.0,
							status:  "Initializing project setup...",
						}
					},
					m.progress.SetPercent(0.0),
					func() tea.Msg {
						return progress.FrameMsg{}
					},
				)
			case stepDevServerPrompt:
				if m.startDevServer {
					go m.startDevelopmentEnvironment()
					m.done = true
					cmds = append(cmds, func() tea.Msg {
						return tea.KeyMsg{Type: tea.KeyEnter}
					})
				} else {
					m.step = stepComplete
				}
			}
			return m, tea.Batch(cmds...)
		}
	}

	return m, tea.Batch(cmds...)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Ternary[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}
