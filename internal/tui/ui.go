package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Download struct {
	URL      string
	Queue    string
	Status   string
	Speed    string
	Progress string
}

type App struct {
	downloads   []Download
	currentPage string
}

func NewApp() *App {
	return &App{
		downloads: []Download{},
	}
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "f1":
			a.currentPage = "Add Download"
		case "f2":
			a.currentPage = "Downloads List"
		case "f3":
			a.currentPage = "Queues List"
		case "a":
			// Switch to Add Download page
			a.currentPage = "Add Download"
		case "d":
			// Handle delete action
		case "e":
			// Handle edit action
		}
	}

	return a, nil
}

func (a *App) View() string {
	var output string
	switch a.currentPage {
	case "Add Download":
		output = a.renderAddDownloadPage()
	case "Downloads List":
		output = a.renderDownloadsListPage()
	case "Queues List":
		output = a.renderQueuesListPage()
	}

	return output + "\nF1: Add Download | F2: Downloads List | F3: Queues List"
}

func (a *App) renderAddDownloadPage() string {
	return "Add Download Page\nEnter the URL and Output Name."
}

func (a *App) renderDownloadsListPage() string {
	var downloadsList string
	for _, d := range a.downloads {
		downloadsList += fmt.Sprintf("URL: %s | Queue: %s | Status: %s | Speed: %s | Progress: %s\n",
			d.URL, d.Queue, d.Status, d.Speed, d.Progress)
	}
	return "Downloads List:\n" + downloadsList
}

func (a *App) renderQueuesListPage() string {
	// This can be customized with your own queue logic
	return "Queues List:\nNo queues available."
}

func (a *App) Run() (tea.Model, error) {
	p := tea.NewProgram(a)
	return p.Run() // Run returns both tea.Model and error
}
