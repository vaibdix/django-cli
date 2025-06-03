package main

import (
	"fmt"
	"time"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

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
		m.progress.Width = msg.Width - padding*2 - 4
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
				cmds = append(cmds, m.mainForm.Init())
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

	if m.step == stepProjectName && m.mainForm != nil {
		if m.mainForm.State != huh.StateCompleted {
			formModel, formCmd := m.mainForm.Update(msg)
			if castedForm, ok := formModel.(*huh.Form); ok {
				m.mainForm = castedForm
			}
			cmds = append(cmds, formCmd)
		} else {
			m.processFormData()
			m.step = stepSetup
			m.progressStatus = "Starting project setup..."
			cmd := m.progress.SetPercent(0.0)
			cmds = append(cmds, cmd)
			go m.CreateProject()
			cmds = append(cmds, m.spinner.Tick)
			return m, tea.Batch(cmds...)
		}
	}

	if m.step == stepDevServerPrompt && m.devServerForm != nil {
		if m.devServerForm.State != huh.StateCompleted {
			formModel, formCmd := m.devServerForm.Update(msg)
			if castedForm, ok := formModel.(*huh.Form); ok {
				m.devServerForm = castedForm
			}
			cmds = append(cmds, formCmd)
		} else {
			if m.startDevServer {
				go m.startDevelopmentEnvironment()
				m.done = true
			} else {
				m.step = stepComplete
			}
			return m, tea.Batch(cmds...)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) getActiveForm() *huh.Form {
	switch m.step {
	case stepProjectName:
		return m.mainForm
	case stepDevServerPrompt:
		return m.devServerForm
	}
	return nil
}

func (m *Model) processFormData() {
	if m.djangoVersion == "" {
		m.djangoVersion = "latest"
	}
	m.createTemplates = contains(m.selectedOptions, "Global Templates")
	m.createAppTemplates = contains(m.selectedOptions, "App Templates")
	m.initializeGit = contains(m.selectedOptions, "Initialize Git")
	m.setupTailwind = contains(m.selectedOptions, "Tailwind")
	m.stepMessages = append(m.stepMessages, "Project name: "+m.projectName)
	m.stepMessages = append(m.stepMessages, "Django version: "+m.djangoVersion)
	if m.appName != "" {
		m.stepMessages = append(m.stepMessages, "App name: "+m.appName)
	}
	m.stepMessages = append(m.stepMessages, fmt.Sprintf("Selected options: %v", m.selectedOptions))
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
