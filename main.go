package main

import (
	"os"

	"cli-casino/casino"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(casino.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
