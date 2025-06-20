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
			Background(lipgloss.Color(ColorPurple)).
			Foreground(lipgloss.Color(ColorLilac)).
			MarginBottom(1).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color(ColorLavender)).
			MarginBottom(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPink)).
			Bold(true).
			MarginBottom(1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSkyBlue)).
			Italic(true).
			MarginTop(1)

	ContentBox = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorSlate))

	ProgressStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPeach)).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorMint)).
			Bold(true)
)

func GetTheme() *huh.Theme {
	theme := huh.ThemeBase()

	// Focused state - active/selected elements
	theme.Focused.Base = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorOrchid))
	theme.Focused.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorLilac))
	theme.Focused.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLavender)).Italic(true)
	theme.Focused.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))
	theme.Focused.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPink))
	theme.Focused.TextInput.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorLilac))

	// Selected options styling
	theme.Focused.SelectedOption = lipgloss.NewStyle().
		Background(lipgloss.Color(ColorPurple)).
		Foreground(lipgloss.Color(ColorLilac)).
		Bold(true)

	theme.Focused.UnselectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorLavender))

	// Blurred state - inactive elements
	theme.Blurred.Base = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))
	theme.Blurred.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSkyBlue)).Bold(true)
	theme.Blurred.Description = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMint)).Italic(true)
	theme.Blurred.TextInput.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))
	theme.Blurred.TextInput.Cursor = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))
	theme.Blurred.TextInput.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))

	theme.Blurred.SelectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSlate))

	theme.Blurred.UnselectedOption = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSlate))

	// Help text styling
	theme.Help.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPeach)).Bold(true)
	theme.Help.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSkyBlue))
	theme.Help.ShortSeparator = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))
	theme.Help.Ellipsis = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSlate))

	return theme
}

func GetSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorOrchid)).
		Bold(true).
		MarginRight(2)
	return s
}

// Additional utility functions for consistent theming
func GetInputStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorLilac)).
		Background(lipgloss.Color(ColorPurple)).
		Padding(0, 1).
		MarginBottom(1)
}

func GetHighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorLilac)).
		Background(lipgloss.Color(ColorOrchid)).
		Bold(true).
		Padding(0, 1)
}

func GetWarningStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorPeach)).
		Bold(true).
		MarginBottom(1)
}
