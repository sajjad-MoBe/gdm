package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"time"
)

// Define your table columns for the Downloads tab
var downloadColumns = []table.Column{
	{Title: "Queue ID", Width: 10},
	{Title: "URL", Width: 30},
	{Title: "Status", Width: 15},
	{Title: "Progress", Width: 10},
	{Title: "Speed", Width: 15},
}

// Sample rows for the Downloads table
var downloadRows = []table.Row{
	{"1", "https://example.com/file1.zip", "Downloading", "50%", "1.2 MB/s"},
	{"2", "https://example.com/file2.zip", "Completed", "100%", "N/A"},
	{"3", "https://example.com/file3.zip", "Paused", "20%", "800 KB/s"},
	{"4", "https://example.com/file4.zip", "Failed", "N/A", "N/A"},
}

// Tabs constants
const (
	tabAddDownload = iota
	tabDownloads
	tabQueues
)

// Model for the table content in Downloads tab
type Model struct {
	// Existing fields
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
	downloadsTable      table.Model
	selectedRow         int
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
	greenTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Italic(true)
	yellowTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true).Italic(true)
	tabActiveStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("4")).Padding(0, 2).Bold(true)
	tabInactiveStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 2)
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	checkmark        = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")
	redErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
)

// Init initializes the UI
func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update method to handle new key presses
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "*":
			return m, tea.Quit
		case "left":
			m.handleTabLeft()
		case "right":
			m.handleTabRight()
		case "enter":
			if m.currentTab == tabAddDownload {
				m.handleEnterPress()
			}
		case "esc":
			if m.currentTab == tabAddDownload {
				m.focusedField = 0
			}
		case "up", "k":
			if m.currentTab == tabAddDownload {
				m.handleUpArrowForTab1()
			}
			if m.currentTab == tabDownloads {
				m.handleUpArrowForTab2()
			}

		case "down", "j":
			if m.currentTab == tabAddDownload {
				m.handleDownArrowForTab1()
			}
			if m.currentTab == tabDownloads {
				m.handleDownArrowForTab2()
			}
		case "tab":
			if m.currentTab == tabAddDownload {
				m.handleTabKey()
			}
		case " ":
			if m.currentTab == tabAddDownload {
				m.handleSpaceKey()
			}
		case "d": // Delete selected download
			if m.currentTab == tabDownloads {
				m.deleteDownload(m.selectedRow)
			}
		case "p": // Pause/Resume selected download
			if m.currentTab == tabDownloads {
				m.togglePauseDownload(m.selectedRow)
			}
		case "r": // Retry selected download if failed
			if m.currentTab == tabDownloads {
				m.retryDownload(m.selectedRow)
			}
		}
	}

	// Update the text inputs based on focus
	if m.currentTab == tabAddDownload {
		if m.focusedField == 0 {
			m.inputURL, cmd = m.inputURL.Update(msg)
		} else if m.focusedField == 1 {
			m.pageSelect, cmd = m.pageSelect.Update(msg)
		} else if m.focusedField == 2 {
			m.outputFileName, cmd = m.outputFileName.Update(msg)
		}
		// Clear messages if necessary
		m.clearMessages()
		// Update the focused field accordingly
		m.updateFocusedField(msg)
	}

	return m, cmd
}

// Helper functions

func (m *Model) handleTabLeft() {
	if m.currentTab > 0 {
		m.currentTab--
	}
}

func (m *Model) handleTabRight() {
	if m.currentTab < tabQueues {
		m.currentTab++
	}
}

func (m *Model) handleEnterPress() {
	if m.currentTab == tabAddDownload {
		m.processURLInput()
	}
}

func (m *Model) processURLInput() {
	if m.inputURL.Value() == "" {
		m.showURLValidationError()
	} else {
		m.showDownloadConfirmation()
	}
}

func (m *Model) showURLValidationError() {
	m.errorMessage = "URL is required!"
	m.confirmationMessage = ""
	m.errorTime = time.Now()

	m.resetFields()
	m.focusedField = 0
	m.inputURL.Focus()
}

func (m *Model) showDownloadConfirmation() {
	m.confirmationMessage = "Download has been added!"
	m.confirmationTime = time.Now()

	m.resetFields()
	m.focusedField = 0
	m.inputURL.Focus()
}

func (m *Model) resetFields() {
	m.inputURL.SetValue("")
	m.pageSelect.ResetSelected()
	m.outputFileName.SetValue("")
	m.selectedFiles = make(map[int]struct{})
}

func (m *Model) handleUpArrowForTab1() {
	if m.focusedField == 1 && m.selectedPage > 0 {
		m.selectedPage--
	}
}

func (m *Model) handleDownArrowForTab1() {
	if m.focusedField == 1 && m.selectedPage < len(m.pageSelect.Items())-1 {
		m.selectedPage++
	}
}

func (m *Model) handleUpArrowForTab2() {
	if m.selectedRow > 0 {
		m.selectedRow--
	}
}

func (m *Model) handleDownArrowForTab2() {
	if m.selectedRow < len(m.downloadsTable.Rows())-1 {
		m.selectedRow++
	}
}

func (m *Model) handleTabKey() {
	if m.currentTab == tabAddDownload {
		m.focusedField = (m.focusedField + 1) % 3
		m.updateFieldFocus()
	}
}

func (m *Model) handleSpaceKey() {
	if m.currentTab == tabAddDownload && m.focusedField == 1 {
		if _, exists := m.selectedFiles[m.selectedPage]; exists {
			delete(m.selectedFiles, m.selectedPage)
		} else {
			m.selectedFiles[m.selectedPage] = struct{}{}
		}
	}
}

func (m *Model) updateFieldFocus() {
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

func (m *Model) clearMessages() {
	if m.errorMessage != "" && time.Since(m.errorTime) > 3*time.Second {
		m.errorMessage = ""
	}

	if m.confirmationMessage != "" && time.Since(m.confirmationTime) > 3*time.Second {
		m.confirmationMessage = ""
	}
}

// Add a method to handle Pause/Resume action
func (m *Model) togglePauseDownload(index int) {
	if index >= 0 && index < len(m.downloadsTable.Rows()) {
		// Check current state of the download
		state := m.downloadsTable.Rows()[index][2]

		if state == "Downloading" {
			// Pause the download
			m.downloadsTable.Rows()[index][2] = "Paused" // Update the state to "Paused"
		} else if state == "Paused" {
			// Resume the download
			m.downloadsTable.Rows()[index][2] = "Downloading" // Update the state to "Downloading"
		}
	}
}

// Modify the deleteDownload method to delete the selected row
func (m *Model) deleteDownload(index int) {
	if index >= 0 && index < len(m.downloadsTable.Rows()) {
		// Remove the row from the table by slicing the rows
		newRows := append(m.downloadsTable.Rows()[:index], m.downloadsTable.Rows()[index+1:]...)

		// Update the downloadsTable with the new rows
		m.downloadsTable = table.New(
			table.WithColumns(downloadColumns), // Keep the existing columns
			table.WithRows(newRows),            // Set the new rows
		)

		// Adjust the selected row to prevent out of bounds error if the last row is deleted
		if m.selectedRow >= len(newRows) {
			m.selectedRow = len(newRows) - 1
		}
	}
}

// Add a method to handle Retry action (only if the state is "Failed")
func (m *Model) retryDownload(index int) {
	if index >= 0 && index < len(m.downloadsTable.Rows()) {
		// Check the state of the selected row
		state := m.downloadsTable.Rows()[index][2]

		if state == "Failed" {
			// Retry the download
			m.downloadsTable.Rows()[index][2] = "Retrying" // Update status to "Retrying"
			// Optionally, trigger actual retry logic here (e.g., retry network request)
		}
	}
}

func (m *Model) updateFocusedField(msg tea.Msg) {
	if m.focusedField == 0 {
		m.inputURL.Update(msg)
	} else if m.focusedField == 1 {
		m.pageSelect.Update(msg)
	} else if m.focusedField == 2 {
		m.outputFileName.Update(msg)
	}
}

func (m *Model) View() string {
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
		content = m.renderAddDownloadTab(tabsRow)
	case tabDownloads:
		content = m.renderDownloadListTab(tabsRow) // Use the new function to render the table
	}

	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(content)
}

func (m *Model) renderAddDownloadTab(tabsRow string) string {
	content := fmt.Sprintf(
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

	return content
}

// Render the downloads table
func (m *Model) renderDownloadListTab(tabsRow string) string {
	// Get columns from global definition
	columns := downloadColumns

	// Custom styles for the table
	tableStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("5")).Padding(1)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))                           // Header color
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))                                         // Row color
	alternateRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))                                // Alternate row color
	selectedRowStyle := lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")) // Highlighted row

	content := fmt.Sprintf(
		"%s\n\n%s\n\n",
		tabsRow,
		tabActiveStyle.Render("Downloads List"),
	)

	// Render the headers
	tableContent := ""
	for _, column := range columns {
		tableContent += headerStyle.Render(fmt.Sprintf("%-*s", column.Width, column.Title))
	}
	tableContent += "\n"

	// Render the rows with their states
	for rowIndex, row := range m.downloadsTable.Rows() {
		var styledRow string
		for colIndex, cell := range row {
			// Highlight the selected row
			if rowIndex == m.selectedRow {
				styledRow += selectedRowStyle.Render(fmt.Sprintf("%-*s", columns[colIndex].Width, cell))
			} else {
				// Alternate row colors for better readability
				if rowIndex%2 == 0 {
					styledRow += rowStyle.Render(fmt.Sprintf("%-*s", columns[colIndex].Width, cell))
				} else {
					styledRow += alternateRowStyle.Render(fmt.Sprintf("%-*s", columns[colIndex].Width, cell))
				}
			}
		}
		tableContent += styledRow + "\n"
	}

	// Wrap the table content in a box style and add it to the final view
	content += tableStyle.Render(tableContent)

	return content
}

func NewModel() *Model {
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

	// Initialize the Downloads table using WithColumns option
	downloadsTable := table.New(
		table.WithColumns(downloadColumns), // Specify columns with WithColumns
		table.WithRows(downloadRows),       // Specify rows
	)

	return &Model{
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
		downloadsTable:      downloadsTable, // Set the Downloads table
	}
}
