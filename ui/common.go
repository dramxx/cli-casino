package ui

import "github.com/charmbracelet/lipgloss"

func Centered(width int, content string) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(content)
}

func Padded(content string) string {
	return lipgloss.NewStyle().Padding(1).Render(content)
}

func Spacer(lines int) string {
	result := ""
	for i := 0; i < lines; i++ {
		result += "\n"
	}
	return result
}
