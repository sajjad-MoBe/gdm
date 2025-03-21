package tui

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/sajjad-mobe/gdm/internal/manager"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Global Variables
var counterForForms = 0
var regForConcurrent = regexp.MustCompile(`^([1-9][0-9]{0,2}|200)$`)
var regForMaxBW = regexp.MustCompile(`^[1-9]\d*$|0`)
var regForHHMMFormat = regexp.MustCompile(`^(?:[01]?[0-9]|2[0-3]):([0-5]?[0-9])$|^$`)

// Define your table columns for the Downloads tab
var downloadColumns = []table.Column{
	{Title: "Download ID", Width: 13},
	{Title: "Queue ID", Width: 10},
	{Title: "URL", Width: 50},
	{Title: "Status", Width: 10},
	{Title: "Progress", Width: 10},
	{Title: "Speed", Width: 10},
	{Title: "Retries", Width: 10},
}

// Define your table columns for the Queues tab
var queueColumns = []table.Column{
	{Title: "Queue ID", Width: 10},
	{Title: "SaveDir", Width: 50},
	{Title: "Max Concurrent", Width: 16},
	{Title: "Max Bandwidth", Width: 15},
	{Title: "Max Retries", Width: 13},
	{Title: "Active Start Time", Width: 19},
	{Title: "Active End Time", Width: 19},
}

// Tabs constants
const (
	tabAddDownload = iota
	tabDownloads
	tabQueues
	tabHelp // New help page tab
)

const (
	minWidth  = 160
	minHeight = 41
)

// Model for the table content in Downloads tab
type Model struct {
	// Existing fields
	currentTab            int
	inputURL              textinput.Model
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
	maxRetriesPerDLInput  textinput.Model
	activeStartTimeInput  textinput.Model
	activeEndTimeInput    textinput.Model
	focusedFieldForQueues int
	dataStore             *manager.DataStore
	maxQueueID            int
	maxDownloadID         int
	downloadmanager       *manager.DownloadManager
	width, height         int
}

// Define styles using LipGloss
var (
	greenTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Italic(true)
	tabActiveStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("4")).Padding(0, 2).Bold(true)
	tabInactiveStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(0, 2)
	cursorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	checkmark        = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("✔")
	redErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
)

// Init initializes the UI
func (m *Model) Init() tea.Cmd {
	return tea.Batch(tickToUpdateDownloadTable())
	// return textInput.Blink
}

func tickToUpdateDownloadTable() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg { return updateDownloadMsg{} })
}

type updateDownloadMsg struct{}

// Update method to handle new key presses
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case updateDownloadMsg:
		if m.currentTab == tabDownloads {
			m.updateDownloadTable()
		}
		return m, tickToUpdateDownloadTable()
	case tea.KeyMsg:
		if m.width < minWidth || m.height < minHeight {
			if msg.String() != "*" {
				// Ignore any key other than "*" until the window is resized.
				return m, nil
			} else if msg.String() == "*" {
				m.dataStore.Save()
				return m, tea.Quit
			}
		}
		switch msg.String() {
		case "*":
			m.dataStore.Save()
			return m, tea.Quit
		case "shift+left":
			m.handleTabLeft()
		case "shift+right":
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
				m.focusedFieldForQueues = (m.focusedFieldForQueues + 1) % 6
				m.updateFocusedFieldForTab3()
			}
		/*case " ":
		if m.currentTab == tabAddDownload {
			m.handleSpaceKey()
		}*/
		case "d": // Remove selected download
			if m.currentTab == tabDownloads {
				m.removeDownload()
			} else if m.currentTab == tabQueues {
				m.removeQueue()
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
			if counterForForms == 0 && m.currentTab == tabQueues {
				m.handleSwitchToAddQueueForm()
			}

		case "e": // Press E to edit a selected queue
			if counterForForms == 0 && m.currentTab == tabQueues {
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
	m.clearMessages()

	return m, cmd
}

func (m *Model) View() string {
	//tempVariable := m.currentTab
	if m.width < minWidth || m.height < minHeight {
		//m.currentTab = -1
		return fmt.Sprintf(
			"Your window is too small: current size %dx%d.\nPlease resize your window to at least %dx%d(width x height)."+
				"\nOr quit the  application by pressing * ",
			m.width, m.height, minWidth, minHeight,
		)
	} //else if  {

	//}

	var renderedTabs []string
	for i := tabAddDownload; i <= tabHelp; i++ {
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
		case tabHelp:
			tabName = "Help"
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
		content = m.renderQueuesTab(tabsRow)
	case tabHelp:
		content = m.renderHelpPage(tabsRow)
	}

	// Compute margins as 5% of the terminal dimensions
	marginWidth := int(float64(m.width) * 0.05)
	marginHeight := int(float64(m.height) * 0.05)

	// Calculate the area available for the TUI content
	contentWidth := m.width - 2*marginWidth
	contentHeight := m.height - 2*marginHeight

	// Place the content in the computed area, aligning to top-left
	placedContent := lipgloss.Place(contentWidth, contentHeight, lipgloss.Center, lipgloss.Center, content)

	// Wrap the placed content with a border
	borderedContent := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(placedContent)

	// Pad the bordered content with the computed margins so it appears 5% from the edges
	finalView := lipgloss.NewStyle().Padding(marginHeight, marginWidth).Render(borderedContent)

	return finalView
}

func (m *Model) renderHelpPage(tabsRow string) string {
	// Style the tabs row with a bold font, background color, and padding.
	helpContent := fmt.Sprintf("%s\n", tabsRow)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#3498db")).
		Foreground(lipgloss.Color("15")).
		BorderForeground(lipgloss.Color("#2980b9"))
	// Define a style for non-header texts: italic with a chosen color.
	textStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#FFC0CB"))

	// Define a different style for Global Keys non-header texts.
	globalTextStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))

	// Add Download Tab section.
	helpContent += headerStyle.Render("Add Download Tab:") + "\n"
	helpContent += textStyle.Render("  Enter: Submits the new download form.") + "\n"
	helpContent += textStyle.Render("  Up/Down Arrows: Navigate through the queue list.") + "\n"
	helpContent += textStyle.Render("  Tab: Cycles focus among URL, queue selection, and output file name.") + "\n"
	helpContent += textStyle.Render("  \"-\": Resets focus back to the URL input field.") + "\n"

	// Downloads Tab section.
	helpContent += headerStyle.Render("Downloads Tab:") + "\n"
	helpContent += textStyle.Render("  Up/Down Arrows: Navigate through the list of downloads.") + "\n"
	helpContent += textStyle.Render("  D: Removes the selected download.") + "\n"
	helpContent += textStyle.Render("  P: Pauses or resumes the selected download.") + "\n"
	helpContent += textStyle.Render("  R: Retries the selected download if it has failed.") + "\n"

	// Queues Tab section.
	helpContent += headerStyle.Render("Queues Tab:") + "\n"
	helpContent += textStyle.Render("  Up/Down Arrows: Navigate through the list of queues.") + "\n"
	helpContent += textStyle.Render("  N: Opens the form for adding a new queue.") + "\n"
	helpContent += textStyle.Render("  E: Opens the form for editing the currently selected queue.") + "\n"
	helpContent += textStyle.Render("  Enter: Submits the queue form (new or edit).") + "\n"
	helpContent += textStyle.Render("  Tab: Cycles through the fields in the queue form.") + "\n"
	helpContent += textStyle.Render("  \"-\": Cancels the current queue form and resets the fields.") + "\n"
	helpContent += textStyle.Render("  D: Removes the selected queue.") + "\n"

	// Global keys section.
	helpContent += headerStyle.Render("Global Keys:") + "\n"
	helpContent += globalTextStyle.Render("  *: Exit help mode when active.") + "\n"
	helpContent += globalTextStyle.Render("  shift+right/left: Navigate through the tabs") + "\n"

	return helpContent
}

func (m *Model) renderQueuesTab(tabsRow string) string {
	if m.newQueueForm {
		return m.renderQueueForm()
	}
	if m.editQueueForm {
		return m.renderQueueFormForEdit()
	}

	columns := queueColumns

	// Define a table style with rounded borders, padding, and margin.
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")).
		Padding(1).
		Margin(1)

	// Build the header row with colorful cells.
	headerRow := ""
	for i, column := range columns {
		// Alternate between lemon yellow (odd) and sky blue (even)
		bgColor := "#87CEEB" // Sky blue for even columns
		if i%2 == 0 {
			bgColor = "#FFFF9F" // Lemon yellow for odd columns
		}

		headerCellStyle := lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(bgColor)).
			Foreground(lipgloss.Color("0")).
			Padding(0, 1)
		headerRow += headerCellStyle.Render(fmt.Sprintf("%-*s", column.Width, column.Title))
	}

	// Build table rows with colorful cells.
	var tableRows []string
	tableRows = append(tableRows, headerRow)

	for rowIndex, row := range m.queuesTable.Rows() {
		rowStr := ""
		for colIndex, cell := range row {
			if columns[colIndex].Title == "SaveDir" && len(cell) > 40 {
				cell = cell[:40] + "..."
			}
			if columns[colIndex].Title == "Max Bandwidth" {
				if cell == "0" {
					cell = "Unlimited"
				}
			}
			// Alternate between lemon yellow (odd) and sky blue (even)
			bgColor := "#87CEEB" // Sky blue for even columns
			if colIndex%2 == 0 {
				bgColor = "#FFFF9F" // Lemon yellow for odd columns
			}

			cellStyle := lipgloss.NewStyle().
				Background(lipgloss.Color(bgColor)).
				Foreground(lipgloss.Color("0")).
				Padding(0, 1)

			// If the row is selected, override with a distinct style.
			if rowIndex == m.selectedRow {
				cellStyle = cellStyle.Copy().
					Background(lipgloss.Color("4")).
					Foreground(lipgloss.Color("0"))
			}

			rowStr += cellStyle.Render(fmt.Sprintf("%-*s", columns[colIndex].Width, cell))
		}
		tableRows = append(tableRows, rowStr)
	}

	// Join all rows vertically.
	tableContent := lipgloss.JoinVertical(lipgloss.Left, tableRows...)

	// Assemble the final content.
	content := fmt.Sprintf("%s\n\n%s", tabsRow, tableStyle.Render(tableContent))
	if m.confirmationMessage != "" {
		content += fmt.Sprintf("\n\n%s", m.confirmationMessage)
	}
	if m.errorMessage != "" {
		content += fmt.Sprintf("\n\n%s", redErrorStyle.Render(m.errorMessage))
	}
	navigationStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))
	content += fmt.Sprintf("\n\n%s", navigationStyle.Render("	  Use shift+right/left to navigate through the tabs."))
	return content
}

func (m *Model) renderAddDownloadTab(tabsRow string) string {
	// File Address field
	var urlCursor string
	if m.focusedField == 0 {
		urlCursor = cursorStyle.Render("> ")
	} else {
		urlCursor = cursorStyle.Render("  ")
	}

	// Initialize content
	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n%s\n\n",
		tabsRow,
		greenTitleStyle.Render("File Address:"),
		urlCursor+m.inputURL.View(),
		greenTitleStyle.Render("Queue Selection:"),
	)

	// Render queues with selection checkboxes
	for i, item := range m.queuesTable.Rows() {
		cursor := " "
		checkbox := "[ ]"

		if m.selectedQueueRowIndex == i {
			if m.focusedField == 1 {
				cursor = ">" // Highlight selected row
			}
			checkbox = "[" + checkmark + "]" // Display checkmark for selected queue
		}

		// Add each queue to content with indentation for clarity
		content += fmt.Sprintf("  %s %s Queue %s\n", cursor, checkbox, item[0])
	}

	// Output File Name field
	var outnameCursor string
	if m.focusedField == 2 {
		outnameCursor = cursorStyle.Render("> ")
	} else {
		outnameCursor = cursorStyle.Render("  ")
	}

	// Add output file name input to content
	content += fmt.Sprintf(
		"\n%s\n%s\n\n",
		greenTitleStyle.Render("Output File Name (optional):"),
		outnameCursor+m.outputFileName.View(),
	)

	// Display error message (if any)
	if m.errorMessage != "" {
		content += fmt.Sprintf("\n\n%s", redErrorStyle.Render(m.errorMessage))
	}

	// Display confirmation message (if any)
	if m.confirmationMessage != "" {
		content += fmt.Sprintf("\n\n%s", greenTitleStyle.Render(m.confirmationMessage))
	}

	// Add navigation message with stylish formatting
	navigationStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))

	content += fmt.Sprintf("\n\n%s", navigationStyle.Render("    Use shift+right/left to navigate through the tabs."))

	return content
}

func (m *Model) renderDownloadListTab(tabsRow string) string {
	// Get columns from global definition
	columns := downloadColumns

	// Define a table style with rounded borders, padding, and margin.
	tableStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")).
		Padding(1).
		Margin(1)

	// Build the header row with colorful cells.
	headerRow := ""
	for i, column := range columns {
		// Alternate between lemon yellow (odd) and sky blue (even)
		bgColor := "#87CEEB" // Sky blue for even columns
		if i%2 == 0 {
			bgColor = "#FFFF9F" // Lemon yellow for odd columns
		}

		headerCellStyle := lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(bgColor)).
			Foreground(lipgloss.Color("0")).
			Padding(0, 1)
		headerRow += headerCellStyle.Render(fmt.Sprintf("%-*s", column.Width, column.Title))
	}

	// Build table rows with colorful cells.
	var tableRows []string
	tableRows = append(tableRows, headerRow)

	for rowIndex, row := range m.downloadsTable.Rows() {
		rowStr := ""
		for colIndex, cell := range row {
			if columns[colIndex].Title == "URL" && len(cell) > 40 {
				cell = cell[:40] + "..."
			}
			// Alternate between lemon yellow (odd) and sky blue (even)
			bgColor := "#87CEEB" // Sky blue for even columns
			if colIndex%2 == 0 {
				bgColor = "#FFFF9F" // Lemon yellow for odd columns
			}

			cellStyle := lipgloss.NewStyle().
				Background(lipgloss.Color(bgColor)).
				Foreground(lipgloss.Color("0")).
				Padding(0, 1)

			// If the row is selected, override with a distinct style.
			if rowIndex == m.selectedRow {
				cellStyle = cellStyle.Copy().
					Background(lipgloss.Color("4")).
					Foreground(lipgloss.Color("0"))
			}

			rowStr += cellStyle.Render(fmt.Sprintf("%-*s", columns[colIndex].Width, cell))
		}
		tableRows = append(tableRows, rowStr)
	}

	// Join all rows vertically.
	tableContent := lipgloss.JoinVertical(lipgloss.Left, tableRows...)

	// Assemble the final content.
	content := fmt.Sprintf("%s\n\n%s", tabsRow, tableStyle.Render(tableContent))
	if m.confirmationMessage != "" {
		content += fmt.Sprintf("\n\n%s", m.confirmationMessage)
	}
	if m.errorMessage != "" {
		content += fmt.Sprintf("\n\n%s", redErrorStyle.Render(m.errorMessage))
	}

	// Apply the navigation style (italic and #F39C12 color)
	navigationStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))

	content += fmt.Sprintf("\n\n%s", navigationStyle.Render("    Use shift+right/left to navigate through the tabs."))

	return content
}

func (m *Model) renderQueueForm() string {
	var content string
	navigationStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))
	// Define a style for yellow italic text
	italicYellowStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#FFFF00"))

	// Display the form header
	content += italicYellowStyle.Render("New Queue\n")

	// Display the fields for Save Directory, Max Concurrent, Max Bandwidth and Max Retries Per Download
	content += fmt.Sprintf(
		"%s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n\n",
		"",
		italicYellowStyle.Render("Save Directory"),
		m.saveDirInput.View(),
		italicYellowStyle.Render("Max Concurrent"),
		m.maxConcurrentInput.View(),
		italicYellowStyle.Render("Max Bandwidth"),
		m.maxBandwidthInput.View(),
		italicYellowStyle.Render("Max Retries per download"),
		m.maxRetriesPerDLInput.View(),
		italicYellowStyle.Render("Active Start Time"),
		m.activeStartTimeInput.View(),
		italicYellowStyle.Render("Active End Time"),
		m.activeEndTimeInput.View(),
	)

	// Add instructions with the same style
	content += fmt.Sprintf("\n%s\n", italicYellowStyle.Render("Press Enter to submit, or \"-\" to cancel."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Note:"))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Time must be in HH:MM format."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Max Concurrent must be an integer from 1 to 200."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Max Bandwidth must be an integer in KB/S. 0 makes no limit."))

	return content
}

func (m *Model) renderQueueFormForEdit() string {
	var content string

	// Define the style for italic yellow text
	italicYellowStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#FFFF00"))

	// Define the style for italic golden yellow text for instructions
	navigationStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("#F39C12"))

	// Display the form header
	content += italicYellowStyle.Render(fmt.Sprintf("Edit Queue %d\n\n", m.selectedRow))

	// Display the fields for Save Directory, Max Concurrent, Max Bandwidth  Max Retries Per Download
	content += fmt.Sprintf(
		"%s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n\n",
		"",
		italicYellowStyle.Render("Save Directory"),
		m.saveDirInput.View(),
		italicYellowStyle.Render("Max Concurrent"),
		m.maxConcurrentInput.View(),
		italicYellowStyle.Render("Max Bandwidth"),
		m.maxBandwidthInput.View(),
		italicYellowStyle.Render("Max Retries per download"),
		m.maxRetriesPerDLInput.View(),
		italicYellowStyle.Render("Active Start Time"),
		m.activeStartTimeInput.View(),
		italicYellowStyle.Render("Active End Time"),
		m.activeEndTimeInput.View(),
	)

	// Add instructions with the same styling
	content += fmt.Sprintf("\n%s\n", navigationStyle.Render("Press Enter to submit, or \"-\" to cancel."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Note:"))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Time must be in HH:MM format."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Max Concurrent must be an integer from 1 to 200."))
	content += fmt.Sprintf("%s\n", navigationStyle.Render("Max Bandwidth must be an integer in KB/S. 0 makes no limit."))

	return content
}

type ByDescending []string

func (s ByDescending) Len() int {
	return len(s)
}
func (s ByDescending) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDescending) Less(i, j int) bool {
	return s[i] > s[j]
}
func NewModel() *Model {
	dataStore := manager.LoadData()
	MaxParts := 10 // Maximum number of parts for one download
	PartSize := 10 // create new part downloader per each PartSize mb
	downloadmanager := manager.NewManager(MaxParts, PartSize)

	ti := textinput.New()
	ti.Placeholder = "Enter Download URL..."
	ti.Focus()

	outputFileName := textinput.New()
	outputFileName.Placeholder = "Optional output file name"
	outputFileName.Blur()

	keys := make([]string, 0, len(dataStore.Queues))
	for key := range dataStore.Queues {
		keys = append(keys, key)
	}
	sort.Sort(ByDescending(keys))

	maxQueueID := 0
	queueRows := []table.Row{}
	for _, key := range keys {
		row := dataStore.Queues[key]
		maxBandwidth := row.MaxBandwidth
		downloadmanager.AddQueue(row)
		if row.ID > maxQueueID {
			maxQueueID = row.ID
		}
		queueRows = append(queueRows, table.Row{
			strconv.Itoa(row.ID),
			row.SaveDir,
			strconv.Itoa(row.MaxConcurrentDownloads),
			strconv.Itoa(maxBandwidth),
			strconv.Itoa(row.MaxRetries),
			row.ActiveStartTime,
			row.ActiveEndTime,
		})
	}
	queuesTable := table.New(
		table.WithColumns(queueColumns), // Specify columns with WithColumns
		table.WithRows(queueRows),       // Specify rows
	)

	keys = make([]string, 0, len(dataStore.Downloads))
	for key := range dataStore.Downloads {
		keys = append(keys, key)
	}
	sort.Sort(ByDescending(keys))
	maxDownloadID := 0

	downloadRows := []table.Row{}
	for _, key := range keys {
		row := dataStore.Downloads[key]
		if row.ID > maxDownloadID {
			maxDownloadID = row.ID
		}
		if row.QueueID == 0 {
			dataStore.RemoveDownload(row)
			continue
		}
		row.Queue = dataStore.Queues[strconv.Itoa(row.QueueID)]

		downloadmanager.AddDownload(row)
		downloadRows = append(downloadRows, table.Row{
			strconv.Itoa(row.ID),
			strconv.Itoa(row.QueueID),
			row.URL,
			row.Status,
			"N/A",
			"N/A",
			"0",
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
	maxConcurrentInput.Placeholder = "Enter Max Concurrent Downloads"

	maxBandwidthInput := textinput.New()
	maxBandwidthInput.Placeholder = "Enter Max Bandwidth. "

	maxRetriesPerDLInput := textinput.New()
	maxRetriesPerDLInput.Placeholder = "Enter Max Retries Per Download. "

	activeStartTimeInput := textinput.New()
	activeStartTimeInput.Placeholder = "Default is 00:00"
	activeEndTimeInput := textinput.New()
	activeEndTimeInput.Placeholder = "Default is 23:59"

	return &Model{
		currentTab:            tabDownloads,
		inputURL:              ti,
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
		maxRetriesPerDLInput:  maxRetriesPerDLInput,
		activeStartTimeInput:  activeStartTimeInput,
		activeEndTimeInput:    activeEndTimeInput,
		focusedFieldForQueues: 0, // Focus on Save Directory initially
		dataStore:             dataStore,
		maxQueueID:            maxQueueID,
		maxDownloadID:         maxDownloadID,
		downloadmanager:       downloadmanager,
	}
}

func (m *Model) updateDownloadTable() {
	var downloadRows []table.Row

	for _, row := range m.downloadsTable.Rows() {
		download := m.dataStore.Downloads[row[0]]
		if download == nil || download.IsRemoved || download.Queue == nil {
			continue
		}
		row[3] = download.GetStatus()

		if download.IsPartial {
			switch download.Status {
			case "finished":
				row[4] = "100%"
			case "initializing":
				row[4] = "N/A"
			default:
				row[4] = strconv.Itoa(download.GetProgress()) + "%"

			}
		} else {
			row[4] = "?"
		}

		if download.GetStatus() != "downloading" {
			row[5] = "-"
		} else {
			speed := download.GetSpeed()
			if speed > 10240 {
				row[5] = fmt.Sprintf("%.1f", float32(speed)/1024) + "Mb/s"
			} else {
				row[5] = strconv.Itoa(download.GetSpeed()) + "Kb/s"
			}

		}
		row[6] = strconv.Itoa(download.Temps.Retries)

		downloadRows = append(downloadRows, row)
	}
	m.downloadsTable.SetRows(downloadRows)
}
