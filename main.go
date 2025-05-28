package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen()) // Using AltScreen for a cleaner exit

	// Pass the program instance to the model so it can send messages
	// from goroutines (like CreateProject)
	m.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}