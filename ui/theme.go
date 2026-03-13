package ui

import "github.com/charmbracelet/lipgloss"

var (
	Primary = lipgloss.Color("#FFD700")
	Win     = lipgloss.Color("#00FF88")
	Lose    = lipgloss.Color("#FF4444")
	Muted   = lipgloss.Color("#666666")
	Border  = lipgloss.Color("#FFFFFF")
)

var (
	HeaderStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Align(lipgloss.Center)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(Border).
			Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Foreground(Primary).
				Background(lipgloss.Color("#333333")).
				Padding(0, 2).
				Bold(true)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(Muted)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	WinStyle = lipgloss.NewStyle().
			Foreground(Win).
			Bold(true)

	LoseStyle = lipgloss.NewStyle().
			Foreground(Lose).
			Bold(true)

	WalletStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Padding(0, 1)
)

func Width(s string) int {
	return lipgloss.Width(s)
}
