package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error starting TUI: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
