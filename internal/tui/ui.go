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

// Global Variables
var counterForForms = 0

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

// Define your table columns for the Queues tab
var queueColumns = []table.Column{
	{Title: "Queue ID", Width: 10},
	{Title: "SaveDir", Width: 20},
	{Title: "Max Concurrent", Width: 15},
	{Title: "Max Bandwidth", Width: 15},
	{Title: "Active Start Time", Width: 20},
	{Title: "Active End Time", Width: 20},
}

// Sample rows for the Queues table
var queueRows = []table.Row{
	{"1", "/path/to/dir1", "5", "100 MB/s", "08:00 AM", "06:00 PM"},
	{"2", "/path/to/dir2", "3", "50 MB/s", "09:00 AM", "05:00 PM"},
	{"3", "/path/to/dir3", "2", "30 MB/s", "10:00 AM", "04:00 PM"},
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
	currentTab            int
	inputURL              textinput.Model
	pageSelect            list.Model
	outputFileName        textinput.Model
	selectedPage          int
	selectedFiles         map[int]struct{} // Tracks selected pages
	focusedField          int              // 0 for inputURL, 1 for pageSelect, 2 for outputFileName
	confirmationMessage   string           // Holds the confirmation message
	errorMessage          string           // Holds the error message (if URL is empty)
	confirmationTime      time.Time        // Time when confirmation message was set
	errorTime             time.Time        // Time when error message was set
	downloadsTable        table.Model
	selectedRow           int
	queuesTable           table.Model // Add the queuesTable field
	editingQueue          *QueueItem  // Holds the queue currently being edited (nil if no queue is being edited)
	newQueueForm          bool        // Flag to indicate if the form for adding a new queue is open
	editQueueForm         bool        // Flag to indicate if the form for adding a new queue is open
	newQueueData          *QueueItem  // Temporarily holds new queue data while filling out the form
	saveDirInput          textinput.Model
	maxConcurrentInput    textinput.Model
	maxBandwidthInput     textinput.Model
	focusedFieldForQueues int // Use focusedFieldForQueues instead of focusedField
}

// QueueItem is the custom type to represent a queue
type QueueItem struct {
	ID              string
	SaveDir         string
	MaxConcurrent   string
	MaxBandwidth    string
	ActiveStartTime string
	ActiveEndTime   string
}

// String implements the list.Item interface
func (q QueueItem) String() string {
	return q.ID
}

// FilterValue implements the list.Item interface
func (q QueueItem) FilterValue() string {
	return q.ID
}

// Define styles using LipGloss
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
			if m.currentTab == tabQueues {
				m.handleNewOrEditQueueFormSubmit()
			}
		case "-":
			if m.currentTab == tabAddDownload {
				m.focusedField = 0
			}
			if m.currentTab == tabQueues {
				m.handleCancel()
			}

		case "up":
			if m.currentTab == tabAddDownload {
				m.handleUpArrowForTab1()
			}
			if m.currentTab == tabDownloads {
				m.handleUpArrowForTab2()
			}
			if m.currentTab == tabQueues {
				m.handleUpArrowForTab3()
			}

		case "down":
			if m.currentTab == tabAddDownload {
				m.handleDownArrowForTab1()
			}
			if m.currentTab == tabDownloads {
				m.handleDownArrowForTab2()
			}
			if m.currentTab == tabQueues {
				m.handleDownArrowForTab3()
			}
		case "tab":
			if m.currentTab == tabAddDownload {
				m.updateFocusedFieldForTab1()
			}
			if m.currentTab == tabQueues {
				m.focusedFieldForQueues = (m.focusedFieldForQueues + 1) % 3
				m.updateFocusedFieldForTab3()
			}
		case " ":
			if m.currentTab == tabAddDownload {
				m.handleSpaceKey()
			}
		case "d": // Delete selected download
			if m.currentTab == tabDownloads {
				m.deleteDownload()
			}
		case "p": // Pause/Resume selected download
			if m.currentTab == tabDownloads {
				m.togglePauseDownload()
			}
		case "r": // Retry selected download if failed
			if m.currentTab == tabDownloads {
				m.retryDownload()
			}
		case "n": // Press N to add a new queue
			if counterForForms == 0 {
				m.handleSwitchToAddQueueForm()
			}

		case "e": // Press E to edit a selected queue
			if counterForForms == 0 {
				m.handleSwitchToEditQueueForm()
			}
		}
	}

	m.updateBasedOnInputForTab1(msg, cmd)

	if counterForForms > 1 {
		m.updateBasedOnInputForTab3(msg, cmd)
	}
	if counterForForms == 1 {
		counterForForms = counterForForms + 1
	}

	return m, cmd
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
		content = m.renderDownloadListTab(tabsRow)
	case tabQueues:
		content = m.renderQueuesTab(tabsRow) // Ensure this renders when in the Queues tab
	}

	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Render(content)
}
func (m *Model) renderQueuesTab(tabsRow string) string {
	if m.newQueueForm {
		return m.renderQueueForm() // Show the form if the user is adding/editing a queue
	}
	if m.editQueueForm {
		return m.renderQueueFormForEdit()
	}

	// Custom styles for the table
	columns := queueColumns
	tableStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("5")).Padding(1)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))                           // Header color
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))                                         // Row color
	alternateRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))                                // Alternate row color
	selectedRowStyle := lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")) // Highlighted row

	content := fmt.Sprintf(
		"%s\n\n",
		tabsRow,
	)

	// Render the headers
	tableContent := ""
	for _, column := range columns {
		tableContent += headerStyle.Render(fmt.Sprintf("%-*s", column.Width, column.Title))
	}
	tableContent += "\n"

	// Render the rows with their states
	for rowIndex, row := range m.queuesTable.Rows() {
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
		"%s\n\n",
		tabsRow,
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

// Render the form for adding or editing a queue
func (m *Model) renderQueueForm() string {
	var content string

	// Display the form header
	content += "New Queue\n\n"

	// Display the fields for Save Directory, Max Concurrent, and Max Bandwidth
	content += fmt.Sprintf(
		"Save Directory: %s\nMax Concurrent: %s\nMax Bandwidth: %s\n\n",
		m.saveDirInput.View(),
		m.maxConcurrentInput.View(),
		m.maxBandwidthInput.View(),
	)

	// Add instructions
	content += "\nPress Enter to submit, ESC to cancel.\n"

	return content
}

// Render the form for adding or editing a queue
func (m *Model) renderQueueFormForEdit() string {
	var content string

	// Display the form header
	content += "Edit Queue\n\n"

	// Display the fields for Save Directory, Max Concurrent, and Max Bandwidth
	content += fmt.Sprintf(
		"Save Directory: %s\nMax Concurrent: %s\nMax Bandwidth: %s\n\n",
		m.saveDirInput.View(),
		m.maxConcurrentInput.View(),
		m.maxBandwidthInput.View(),
	)

	// Add instructions
	content += "\nPress Enter to submit, ESC to cancel.\n"

	return content
}

func NewModel() *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter Download URL..."
	ti.Focus()

	pageSelect := list.New([]list.Item{
		QueueItem{ID: "Page 1"},
		QueueItem{ID: "Page 2"},
		QueueItem{ID: "Page 3"},
	}, list.NewDefaultDelegate(), 0, 0)

	outputFileName := textinput.New()
	outputFileName.Placeholder = "Optional output file name"
	outputFileName.Blur()

	// Initialize the Downloads table using WithColumns option
	downloadsTable := table.New(
		table.WithColumns(downloadColumns), // Specify columns with WithColumns
		table.WithRows(downloadRows),       // Specify rows
	)

	queuesTable := table.New(
		table.WithColumns(queueColumns), // Specify columns with WithColumns
		table.WithRows(queueRows),       // Specify rows
	)

	saveDirInput := textinput.New()
	saveDirInput.Placeholder = "Enter Save Directory"

	maxConcurrentInput := textinput.New()
	maxConcurrentInput.Placeholder = "Enter Max Concurrent"

	maxBandwidthInput := textinput.New()
	maxBandwidthInput.Placeholder = "Enter Max Bandwidth"

	return &Model{
		currentTab:            tabAddDownload,
		inputURL:              ti,
		pageSelect:            pageSelect,
		outputFileName:        outputFileName,
		selectedPage:          0,
		selectedFiles:         make(map[int]struct{}),
		focusedField:          0,
		confirmationMessage:   "",
		errorMessage:          "",
		confirmationTime:      time.Time{},
		errorTime:             time.Time{},
		downloadsTable:        downloadsTable, // Set the Downloads table
		queuesTable:           queuesTable,    // Set the Queues table
		saveDirInput:          saveDirInput,
		maxConcurrentInput:    maxConcurrentInput,
		maxBandwidthInput:     maxBandwidthInput,
		focusedFieldForQueues: 0, // Focus on Save Directory initially
	}
}
