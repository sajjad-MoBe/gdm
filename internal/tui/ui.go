package tui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sajjad-mobe/gdm/internal/manager"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Global Variables
var counterForForms = 0

// Define your table columns for the Downloads tab
var downloadColumns = []table.Column{
	{Title: "Download ID", Width: 15},
	{Title: "URL", Width: 60},
	{Title: "Status", Width: 15},
	{Title: "Progress", Width: 10},
	{Title: "Speed", Width: 15},
}

// Sample rows for the Downloads table
// var downloadRows = []table.Row{
// 	{"1", "https://example.com/file1.zip", "Downloading", "50%", "1.2 MB/s"},
// 	{"2", "https://example.com/file2.zip", "Completed", "100%", "N/A"},
// 	{"3", "https://example.com/file3.zip", "Paused", "20%", "800 KB/s"},
// 	{"4", "https://example.com/file4.zip", "Failed", "N/A", "N/A"},
// }

// Define your table columns for the Queues tab
var queueColumns = []table.Column{
	{Title: "Queue ID", Width: 15},
	{Title: "SaveDir", Width: 60},
	{Title: "Max Concurrent", Width: 15},
	{Title: "Max Bandwidth", Width: 15},
	{Title: "Active Start Time", Width: 20},
	{Title: "Active End Time", Width: 20},
}

// Sample rows for the Queues table
// var queueRows = []table.Row{
// 	{"1", "/path/to/dir1", "5", "100 MB/s", "08:00 AM", "06:00 PM"},
// 	{"2", "/path/to/dir2", "3", "50 MB/s", "09:00 AM", "05:00 PM"},
// 	{"3", "/path/to/dir3", "2", "30 MB/s", "10:00 AM", "04:00 PM"},
// }

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
	queueSelect           list.Model
	outputFileName        textinput.Model
	selectedQueueRowIndex int       // Tracks selected pages
	focusedField          int       // 0 for inputURL, 1 for queueSelect, 2 for outputFileName
	confirmationMessage   string    // Holds the confirmation message
	errorMessage          string    // Holds the error message (if URL is empty)
	confirmationTime      time.Time // Time when confirmation message was set
	errorTime             time.Time // Time when error message was set
	downloadsTable        table.Model
	selectedRow           int
	queuesTable           table.Model // Add the queuesTable field
	// editingQueue          *manager.Queue // Holds the queue currently being edited (nil if no queue is being edited)
	newQueueForm  bool // Flag to indicate if the form for adding a new queue is open
	editQueueForm bool // Flag to indicate if the form for adding a new queue is open
	// newQueueData          *manager.Queue // Temporarily holds new queue data while filling out the form
	saveDirInput          textinput.Model
	maxConcurrentInput    textinput.Model
	maxBandwidthInput     textinput.Model
	activeStartTimeInput  textinput.Model
	activeEndTimeInput    textinput.Model
	focusedFieldForQueues int
	queues                map[string]*manager.Queue
	downloads             map[string]*manager.Download

	downloadmanager *manager.DownloadManager
}

// QueueItem is the custom type to represent a queue
// type QueueItem struct {
// 	ID              string
// 	SaveDir         string
// 	MaxConcurrent   string
// 	MaxBandwidth    string
// 	ActiveStartTime string
// 	ActiveEndTime   string
// }

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
	return tea.Batch(tickToUpdateDownloadTable())
	// return textinput.Blink
}

func tickToUpdateDownloadTable() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg { return updateDownloadMsg{} })
}

type updateDownloadMsg struct{}

// Update method to handle new key presses
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case updateDownloadMsg:
		if m.currentTab == tabDownloads {
			m.updateDownloadTable()
		}
		return m, tickToUpdateDownloadTable()
	case tea.KeyMsg:
		switch msg.String() {
		case "*":
			for _, row := range m.downloads {
				manager.Save(row)
			}

			return m, tea.Quit
		case "ctrl+left":
			m.handleTabLeft()
		case "ctrl+right":
			m.handleTabRight()
		case "enter":
			if m.currentTab == tabAddDownload {
				m.handleNewDownloadSubmit()
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
				m.focusedFieldForQueues = (m.focusedFieldForQueues + 1) % 5
				m.updateFocusedFieldForTab3()
			}
		case " ":
			if m.currentTab == tabAddDownload {
				m.handleSpaceKey()
			}
		case "d": // Delete selected download
			if m.currentTab == tabDownloads {
				m.deleteDownload()
			} else if m.currentTab == tabQueues {
				m.deleteQueue()
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
		return m.renderQueueForm()
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
	var urlCursor string
	if m.focusedField == 0 {
		urlCursor = cursorStyle.Render("> ")
	} else {
		urlCursor = cursorStyle.Render("  ")
	}
	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n%s\n\n",
		tabsRow,
		greenTitleStyle.Render("File Address:"),
		urlCursor+m.inputURL.View(),
		greenTitleStyle.Render("Queue Selection:"),
	)

	for i, item := range m.queuesTable.Rows() {
		cursor := " "
		checkbox := "[ ]"

		if m.selectedQueueRowIndex == i {
			if m.focusedField == 1 {
				cursor = ">"
			}
			checkbox = "[" + checkmark + "]"
		}

		content += fmt.Sprintf("%s %s Queue %s\n", cursor, checkbox, item[0])

	}

	var outnameCursor string
	if m.focusedField == 2 {
		outnameCursor = cursorStyle.Render("> ")
	} else {
		outnameCursor = cursorStyle.Render("  ")
	}

	content += fmt.Sprintf(
		"\n%s\n%s\n\n%s\n\n",
		greenTitleStyle.Render("Output File Name (optional):"),
		outnameCursor+m.outputFileName.View(),
		yellowTitleStyle.Render("Press Enter to confirm, up/down to select Queue, ESC to \n"+
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

func (m *Model) renderQueueForm() string {
	var content string

	// Display the form header
	content += "New Queue\n\n"

	// Display the fields for Save Directory, Max Concurrent, and Max Bandwidth
	content += fmt.Sprintf(
		"Save Directory: %s\nMax Concurrent: %s\nMax Bandwidth: %s\nActive Start Time: %s\nActive End Time: %s\n\n",
		m.saveDirInput.View(),
		m.maxConcurrentInput.View(),
		m.maxBandwidthInput.View(),
		m.activeStartTimeInput.View(),
		m.activeEndTimeInput.View(),
	)
	// Add instructions
	content += "\nPress Enter to submit, ESC to cancel.\n"

	return content
}

// Render the form for adding or editing a queue
func (m *Model) renderQueueFormForEdit() string {
	var content string

	// Display the form header
	content += fmt.Sprintf("Edit Queue %d\n\n", m.selectedRow)

	// Display the fields for Save Directory, Max Concurrent, and Max Bandwidth
	content += fmt.Sprintf(
		"Save Directory: %s\nMax Concurrent: %s\nMax Bandwidth: %s\nActive Start Time: %s\nActive End Time: %s\n\n",
		m.saveDirInput.View(),
		m.maxConcurrentInput.View(),
		m.maxBandwidthInput.View(),
		m.activeStartTimeInput.View(),
		m.activeEndTimeInput.View(),
	)

	// Add instructions
	content += "\nPress Enter to submit, ESC to cancel.\n"

	return content
}

func NewModel() *Model {
	manager.InitializeDB()
	MaxParts := 10 // Maximum number of parts for one download
	PartSize := 10 // create new part downloader per each PartSize mb
	downloadmanager := manager.NewManager(MaxParts, PartSize)

	ti := textinput.New()
	ti.Placeholder = "Enter Download URL..."
	ti.Focus()

	outputFileName := textinput.New()
	outputFileName.Placeholder = "Optional output file name"
	outputFileName.Blur()

	var queues []*manager.Queue
	if err := manager.GetAllQueues(&queues); err != nil {
		fmt.Printf("failed to get queues: %v", err)
	}

	var queuesList []list.Item
	queueRows := []table.Row{}
	queuesMap := make(map[string]*manager.Queue)
	for _, row := range queues {
		maxBandwidth := row.MaxBandwidth
		queuesList = append(queuesList, row)
		downloadmanager.AddQueue(row)
		queuesMap[strconv.Itoa(row.ID)] = row

		queueRows = append(queueRows, table.Row{
			strconv.Itoa(row.ID),
			row.SaveDir,
			strconv.Itoa(row.MaxConcurrentDownloads),
			strconv.Itoa(maxBandwidth),
			row.ActiveStartTime,
			row.ActiveEndTime,
		})
	}
	queuesTable := table.New(
		table.WithColumns(queueColumns), // Specify columns with WithColumns
		table.WithRows(queueRows),       // Specify rows
	)
	queueSelect := list.New(queuesList, list.NewDefaultDelegate(), 0, 0)

	var downloads []*manager.Download
	if err := manager.GetAllDownloads(&downloads); err != nil {
		fmt.Printf("failed to get downloads: %v", err)
	}
	downloadRows := []table.Row{}
	downloadsMap := make(map[string]*manager.Download)
	for _, row := range downloads {
		// row.QueueID = queuesMap[strconv.Itoa(row.QueueID)].ID
		if row.Queue == nil {
			manager.Delete(row)
			continue
		}
		row.Queue = queuesMap[strconv.Itoa(row.QueueID)]
		// manager.Save(row)

		downloadmanager.AddDownload(row)
		downloadsMap[strconv.Itoa(row.ID)] = row
		downloadRows = append(downloadRows, table.Row{
			strconv.Itoa(row.ID),
			row.URL,
			row.Status,
			"N/A",
			"N/A",
		})
	}
	// Initialize the Downloads table using WithColumns option
	downloadsTable := table.New(
		table.WithColumns(downloadColumns), // Specify columns with WithColumns
		table.WithRows(downloadRows),       // Specify rows
	)

	saveDirInput := textinput.New()
	saveDirInput.Placeholder = "Enter Save Directory"

	maxConcurrentInput := textinput.New()
	maxConcurrentInput.Placeholder = "Enter Max Concurrent"

	maxBandwidthInput := textinput.New()
	maxBandwidthInput.Placeholder = "Enter Max Bandwidth"

	activeStartTimeInput := textinput.New()
	activeStartTimeInput.Placeholder = "Enter Active Start Time"
	activeEndTimeInput := textinput.New()
	activeEndTimeInput.Placeholder = "Enter Active End Time"

	return &Model{
		currentTab:            tabAddDownload,
		inputURL:              ti,
		queueSelect:           queueSelect,
		outputFileName:        outputFileName,
		selectedQueueRowIndex: 0,
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
		activeStartTimeInput:  activeStartTimeInput,
		activeEndTimeInput:    activeEndTimeInput,
		focusedFieldForQueues: 0, // Focus on Save Directory initially
		queues:                queuesMap,
		downloads:             downloadsMap,
		downloadmanager:       downloadmanager,
	}
}

func (m *Model) updateDownloadTable() {
	var downloadRows []table.Row

	for _, row := range m.downloadsTable.Rows() {
		download := m.downloads[row[0]]
		if download == nil || download.IsDeleted || download.Queue == nil {
			continue
		}
		row[2] = download.GetStatus()

		if download.IsPartial {
			switch download.Status {
			case "finished":
				row[3] = "100%"
			case "initializing":
				row[3] = "N/A"
			default:
				row[3] = strconv.Itoa(download.GetProgress()) + "%"

			}
		} else {
			row[3] = "?"
		}

		if download.GetStatus() != "downloading" {
			row[4] = "-"
		} else {
			speed := download.GetSpeed()
			if speed > 1024 {
				row[4] = fmt.Sprintf("%.1f", float32(speed)/1024) + "Mb/s"
			} else {
				row[4] = strconv.Itoa(download.GetSpeed()) + "Kb/s"
			}

		}
		downloadRows = append(downloadRows, row)
	}
	m.downloadsTable.SetRows(downloadRows)
}
