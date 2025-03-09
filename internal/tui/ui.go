package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
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
	currentTab          int
	inputURL            textinput.Model
	pageSelect          list.Model
	outputFileName      textinput.Model
	selectedPage        int
	selectedFiles       map[int]struct{} // Tracks selected pages
	focusedField        int              // 0 for inputURL, 1 for pageSelect, 2 for outputFileName
	confirmationMessage string           // Holds the confirmation message
	errorMessage        string           // Holds the error message (if URL is empty)
	confirmationTime    time.Time        // Time when confirmation message was set
	errorTime           time.Time        // Time when error message was set
}

// QueueItem is the custom type to represent a queue
type QueueItem struct {
	name string
}

// String implements the list.Item interface
func (q QueueItem) String() string {
	return q.name
}

// FilterValue implements the list.Item interface
func (q QueueItem) FilterValue() string {
	return q.name
}

// Define styles using Lipgloss
var (
	greenTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true).
			Italic(true)

	yellowTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")).
				Bold(true).
				Italic(true)

	tabActiveStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("4")).
			Padding(0, 2).
			Bold(true)

	tabInactiveStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				Padding(0, 2)

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	checkmark   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")

	redErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
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
		case "*":
			return m, tea.Quit
		case "left":
			if m.currentTab > 0 {
				m.currentTab--
			}
		case "right":
			if m.currentTab < tabQueues {
				m.currentTab++
			}
		case "enter":
			if m.currentTab == tabAddDownload {
				// Check if the URL is empty
				if m.inputURL.Value() == "" {
					// URL is empty, reset and show error message
					m.errorMessage = "URL is required!"
					m.confirmationMessage = "" // Clear any previous confirmation message
					m.errorTime = time.Now()   // Set the error message time

					// Reset the fields
					m.inputURL.SetValue("")                  // Clear inputURL field
					m.pageSelect.ResetSelected()             // Reset page selection
					m.outputFileName.SetValue("")            // Clear outputFileName field
					m.selectedFiles = make(map[int]struct{}) // Reset selected pages

					// Set focus back to the URL input field
					m.focusedField = 0
					m.inputURL.Focus()

				} else {
					// URL is not empty, proceed with regular flow
					m.confirmationMessage = "Download has been added!"
					m.confirmationTime = time.Now() // Set the confirmation message time

					// Reset the fields after adding the download
					m.inputURL.SetValue("")       // Clear inputURL field
					m.pageSelect.ResetSelected()  // Reset page selection
					m.outputFileName.SetValue("") // Clear outputFileName field

					// Reset the selected files (deselect each page manually)
					m.selectedFiles = make(map[int]struct{})

					// Optionally, you can set focus back to the first field (inputURL) after resetting
					m.focusedField = 0
					m.inputURL.Focus()
				}
			}
		case "esc":
			m.focusedField = 0
		case "up", "k":
			if m.currentTab == tabAddDownload && m.focusedField == 1 && m.selectedPage > 0 {
				m.selectedPage--
			}
		case "down", "j":
			if m.currentTab == tabAddDownload && m.focusedField == 1 && m.selectedPage < len(m.pageSelect.Items())-1 {
				m.selectedPage++
			}
		case "tab":
			if m.currentTab == tabAddDownload {
				m.focusedField = (m.focusedField + 1) % 3
				if m.focusedField == 0 {
					m.inputURL.Focus()
					m.outputFileName.Blur()
				} else if m.focusedField == 1 {
					m.inputURL.Blur()
					m.outputFileName.Blur()
				} else if m.focusedField == 2 {
					m.outputFileName.Focus()
					m.inputURL.Blur()
				}
			}
		case " ":
			if m.currentTab == tabAddDownload && m.focusedField == 1 {
				if _, exists := m.selectedFiles[m.selectedPage]; exists {
					delete(m.selectedFiles, m.selectedPage)
				} else {
					m.selectedFiles[m.selectedPage] = struct{}{}
				}
			}
		}
	}

	// Check if 3 seconds have passed since the message was set, and clear the messages if so
	if m.errorMessage != "" && time.Since(m.errorTime) > 3*time.Second {
		m.errorMessage = ""
	}

	if m.confirmationMessage != "" && time.Since(m.confirmationTime) > 3*time.Second {
		m.confirmationMessage = ""
	}

	var cmd tea.Cmd
	if m.focusedField == 0 {
		m.inputURL, cmd = m.inputURL.Update(msg)
	} else if m.focusedField == 1 {
		m.pageSelect, cmd = m.pageSelect.Update(msg)
	} else if m.focusedField == 2 {
		m.outputFileName, cmd = m.outputFileName.Update(msg)
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	var renderedTabs []string
	for i := tabAddDownload; i <= tabQueues; i++ {
		var style lipgloss.Style
		if i == m.currentTab {
			style = tabActiveStyle
		} else {
			style = tabInactiveStyle
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

	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	var content string
	switch m.currentTab {
	case tabAddDownload:
		content = fmt.Sprintf(
			"%s\n\n%s\n%s\n%s\n\n",
			tabsRow,
			greenTitleStyle.Render("File Address:"),
			cursorStyle.Render("> ")+m.inputURL.View(),
			greenTitleStyle.Render("Page Selection:"),
		)

		for i, item := range m.pageSelect.Items() {
			cursor := " "
			checkbox := "[ ]"

			if m.selectedPage == i {
				cursor = ">"
			}
			if _, selected := m.selectedFiles[i]; selected {
				checkbox = "[" + checkmark + "]"
			}
			if queueItem, ok := item.(QueueItem); ok {
				content += fmt.Sprintf("%s %s %s\n", cursor, checkbox, queueItem.String())
			}
		}

		content += fmt.Sprintf(
			"\n%s\n%s\n\n%s\n\n",
			greenTitleStyle.Render("Output File Name (optional):"),
			cursorStyle.Render("> ")+m.outputFileName.View(),
			yellowTitleStyle.Render("Press Enter to confirm, Space to select pages, ESC to \n"+
				"cancel/reset, or * to quit the download manager."),
		)

		// Show the error message in red if it exists
		if m.errorMessage != "" {
			content += fmt.Sprintf("\n\n%s", redErrorStyle.Render(m.errorMessage))
		}

		// Show the confirmation message if it exists
		if m.confirmationMessage != "" {
			content += fmt.Sprintf("\n\n%s", m.confirmationMessage)
		}

	case tabDownloads:
		content = tabActiveStyle.Render("Downloads List") + "\n\n[Press ←/→ to switch tabs]"
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(content)
}

// NewModel initializes the UI model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Download URL..."
	ti.Focus()

	pageSelect := list.New([]list.Item{
		QueueItem{name: "Page 1"},
		QueueItem{name: "Page 2"},
		QueueItem{name: "Page 3"},
	}, list.NewDefaultDelegate(), 0, 0)

	outputFileName := textinput.New()
	outputFileName.Placeholder = "Optional output file name"
	outputFileName.Blur()

	return Model{
		currentTab:          tabAddDownload,
		inputURL:            ti,
		pageSelect:          pageSelect,
		outputFileName:      outputFileName,
		selectedPage:        0,
		selectedFiles:       make(map[int]struct{}),
		focusedField:        0,
		confirmationMessage: "",
		errorMessage:        "",
		confirmationTime:    time.Time{},
		errorTime:           time.Time{},
	}
}
