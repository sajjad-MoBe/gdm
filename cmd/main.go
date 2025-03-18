package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sajjad-mobe/gdm/internal/tui"
)

func main() {
	// Initialize the program with a pointer to Model

	p := tea.NewProgram(tui.NewModel()) // Now it's fine to use it like this
	if _, err := p.Run(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error starting TUI: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
