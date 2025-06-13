package main

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorLavender = "#e0d4f0"
	ColorPink     = "#ff9ac4"
	ColorSkyBlue  = "#b8e8ff"
	ColorMint     = "#a8ddc8"
	ColorPurple   = "#c9b8f5"
	ColorOrchid   = "#f0b8ec"
	ColorPeach    = "#ffcab0"
	ColorLilac    = "#f8f0ff"
	ColorSlate    = "#8b7394"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(ColorPeach)).
			Foreground(lipgloss.Color(ColorLilac)).
			MarginBottom(1)

	// Use ColorLavender for subtitles - second lightest
	SubtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color(ColorLavender)).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{
			Light: "#DC2626",
			Dark:  ColorPink,
		}).
		Bold(true).
		MarginBottom(1)

	// Use ColorSkyBlue for footer text
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSkyBlue)).
			Italic(true).
			MarginTop(1)

	// Use ColorMint for borders
	ContentBox = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSlate))

	// Use ColorPeach for progress indicators
	ProgressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPeach)).
			Bold(true)
)

func GetTheme() *huh.Theme {
	theme := huh.ThemeBase()

	// Focused state - use your brightest colors
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorOrchid))
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorLilac))
	theme.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLavender)).Italic(true)
	theme.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPurple))
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPink))

	// Blurred state - use medium brightness colors
	theme.Blurred.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorOrchid))
	theme.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSkyBlue)).Bold(true)
	theme.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMint)).Italic(true)
	theme.Blurred.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))

	return theme
}

func GetSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorOrchid)). // Use ColorOrchid for spinner
		Bold(true).
		MarginRight(2)
	return s
}
