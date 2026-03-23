package ui

import "github.com/charmbracelet/lipgloss"

var (
	MainColor      = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#04B575")
	AccentColor    = lipgloss.Color("#EF910F")
	ErrorColor     = lipgloss.Color("#F25D94")
	MutedColor     = lipgloss.Color("#626262")

	TitleStyle = lipgloss.NewStyle().
			Foreground(MainColor).
			Bold(true).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(MainColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	PromptStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	StatusStyle = lipgloss.NewStyle().
			Foreground(MutedColor)
)
