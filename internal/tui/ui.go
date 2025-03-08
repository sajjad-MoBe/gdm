package tui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Tabs constants
const (
	tabAddDownload = iota
	tabDownloads
	tabQueues
)

// Model defines the UI state
type Model struct {
	currentTab int
	inputURL   textinput.Model
	table      table.Model
	list       list.Model
	typing     bool
	loading    bool
	err        error

	// Additional state for downloads and queues
	downloads list.Model // Displaying ongoing downloads
	queues    list.Model // Displaying queues
}

// Define color styles using lipgloss
var (
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // Yellow
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4")) // Blue
	magentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")) // Magenta

	// Tab Styles
	inactiveTabStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 1)
	activeTabStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("4")).Padding(0, 1)
)

// Init initializes the UI
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages (keypresses)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q": // Quit program
			return m, tea.Quit
		case "left": // Left arrow: Previous tab
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right": // Right arrow: Next tab
			if m.currentTab < tabQueues {
				m.currentTab++
			}
		case "enter": // Handle enter key
			switch m.currentTab {
			case tabAddDownload:
				query := m.inputURL.Value()
				if query != "" {
					// Add the download to the list (simulating)
					m.loading = true
					m.err = nil
				}
			default:
				//Do Nothing
			}
		case "esc": // Escape key
			if !m.typing && !m.loading {
				m.typing = true
				m.err = nil
			}
		case "d": // Delete selected download
			if m.currentTab == tabDownloads {
				// Delete download logic here
			}
		case "p": // Pause/Resume download
			if m.currentTab == tabDownloads {
				// Pause/Resume logic here
			}
		case "r": // Retry failed download
			if m.currentTab == tabDownloads {
				// Retry logic here
			}
		}
	}

	// Update the components based on the current tab
	if m.typing {
		var cmd tea.Cmd
		m.inputURL, cmd = m.inputURL.Update(msg)
		return m, cmd
	}

	if m.loading {
		// Simulate loading action (this could be customized with a spinner or something similar)
		m.loading = false
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	// Handle tabs rendering
	var renderedTabs []string
	for i := tabAddDownload; i <= tabQueues; i++ {
		var style lipgloss.Style
		if i == m.currentTab {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		var tabName string
		switch i {
		case tabAddDownload:
			tabName = "Add Download"
		case tabDownloads:
			tabName = "Downloads"
		case tabQueues:
			tabName = "Queues"
		}
		renderedTabs = append(renderedTabs, style.Render(tabName))
	}

	// Join the tabs horizontally
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Render content for the active tab
	var content string
	switch m.currentTab {
	case tabAddDownload:
		content = fmt.Sprintf("%s\n\nURL: %s\n\n[Press ←/→ to switch tabs]", greenStyle.Render("Add Download"), m.inputURL.View())
	case tabDownloads:
		content = yellowStyle.Render("Downloads List") + "\n\n[Press ←/→ to switch tabs]"
		if len(m.downloads.Items()) > 0 {
			content += "\n\n" + m.downloads.View()
		} else {
			content += "\n\n" + redStyle.Render("No downloads available.")
		}
	case tabQueues:
		content = magentaStyle.Render("Queues List") + "\n\n[Press ←/→ to switch tabs]"
		if len(m.queues.Items()) > 0 {
			content += "\n\n" + m.queues.View()
		} else {
			content += "\n\n" + redStyle.Render("No queues available.")
		}
	}

	// Render the final UI with borders and styles
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(fmt.Sprintf("%s\n\n%s", tabsRow, content))
}

// NewModel initializes the UI model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Download URL..."
	ti.Focus()

	// Set up the lists for downloads and queues (can be expanded)
	downloads := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	queues := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	// If file doesn't exist, create default file and load data
	if _, err := os.Stat("data.json"); os.IsNotExist(err) {
		// Create and load default data
		// Note: This is a placeholder and should be expanded with actual data creation logic.
	}

	return Model{
		currentTab: tabAddDownload,
		inputURL:   ti,
		typing:     true,
		downloads:  downloads,
		queues:     queues,
	}
}
