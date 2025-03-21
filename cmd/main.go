package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	// Use tea.WithAltScreen() to enable full terminal usage
	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error starting TUI: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
