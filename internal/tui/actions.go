package tui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sajjad-mobe/gdm/internal/manager"
)

// Helper functions

func (m *Model) handleTabLeft() {
	if m.currentTab > 0 && !m.newQueueForm && !m.editQueueForm {
		m.currentTab--
		m.selectedRow = 0
	}
}

func (m *Model) handleTabRight() {
	if m.currentTab < tabHelp && !m.newQueueForm && !m.editQueueForm {
		m.currentTab++
		m.selectedRow = 0
	}
}

// Handle the submission of a new download form
func (m *Model) handleNewDownloadSubmit() {
	if m.inputURL.Value() == "" {
		m.showURLValidationError()
	} else if len(m.dataStore.Queues) == 0 {
		m.showCreateQueueError()
	} else {
		// Create a new download with the data entered in fields
		downloadURL := m.inputURL.Value()
		// validate download url

		outputFile := m.outputFileName.Value()
		// should be validated
		queue := m.dataStore.Queues[m.queuesTable.Rows()[m.selectedQueueRowIndex][0]]
		m.maxDownloadID++
		newDwnload := manager.Download{
			ID:         m.maxDownloadID,
			URL:        downloadURL,
			QueueID:    queue.ID,
			Queue:      queue,
			OutputFile: outputFile,
			Status:     "pending",
		}

		m.addNewDownload(&newDwnload)

		// Reset the form after submission
		m.inputURL.Reset()
		m.outputFileName.Reset()

		m.showDownloadConfirmation()
	}
}

// Add a new download
func (m *Model) addNewDownload(download *manager.Download) {
	m.dataStore.AddDownload(download)
	m.downloadmanager.AddDownload(download)

	newRow := table.Row{
		strconv.Itoa(download.ID),
		strconv.Itoa(download.QueueID),
		download.URL,
		download.Status,
		"N/A",
		"N/A",
		"0",
	}

	// Add the row to the downloadsTable
	m.downloadsTable = table.New(
		table.WithColumns(downloadColumns),                                      // Keep the existing columns
		table.WithRows(append([]table.Row{newRow}, m.downloadsTable.Rows()...)), // Add the new row
	)
}

func (m *Model) resetFieldsForTab1() {
	m.inputURL.SetValue("")
	m.outputFileName.SetValue("")
	m.selectedQueueRowIndex = 0
}

func (m *Model) resetFieldsForTab3() {
	m.saveDirInput.SetValue("")
	m.maxBandwidthInput.SetValue("")
	m.maxRetriesPerDLInput.SetValue("")
	m.maxConcurrentInput.SetValue("")
	m.activeStartTimeInput.SetValue("")
	m.activeEndTimeInput.SetValue("")
}

func (m *Model) handleUpArrowForTab1() {
	if m.focusedField == 1 {
		m.selectedQueueRowIndex = (m.selectedQueueRowIndex + len(m.queuesTable.Rows()) - 1) % len(m.queuesTable.Rows())
	}
}

func (m *Model) handleDownArrowForTab1() {
	if m.focusedField == 1 {
		m.selectedQueueRowIndex = (m.selectedQueueRowIndex + 1) % len(m.queuesTable.Rows())
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

func (m *Model) updateFocusedFieldForTab1() {
	if m.currentTab == tabAddDownload {
		url := m.inputURL.Value()
		if m.focusedField == 0 && len(m.outputFileName.Value()) == 0 {
			if outputFileName, err := manager.GetFileNameFromURL(url); err == nil {
				m.outputFileName.SetValue(outputFileName)
			}
		}
		m.focusedField = (m.focusedField + 1) % 3
		m.updateFieldFocus()
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

// Update focus for the queue form fields
func (m *Model) updateFocusedFieldForTab3() {
	m.saveDirInput.Blur()
	m.maxConcurrentInput.Blur()
	m.maxBandwidthInput.Blur()
	m.maxRetriesPerDLInput.Blur()
	m.activeStartTimeInput.Blur()
	m.activeEndTimeInput.Blur()
	switch m.focusedFieldForQueues {
	case 0:
		m.saveDirInput.Focus()
	case 1:
		m.maxConcurrentInput.Focus()
	case 2:
		m.maxBandwidthInput.Focus()
	case 3:
		m.maxRetriesPerDLInput.Focus()
	case 4:
		m.activeStartTimeInput.Focus()
	case 5:
		m.activeEndTimeInput.Focus()
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
func (m *Model) togglePauseDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Check current state of the download
		state := m.downloadsTable.Rows()[m.selectedRow][3]

		if state == "downloading" {
			// Pause the download
			m.downloadsTable.Rows()[m.selectedRow][3] = "paused" // Update the state to "Paused"
			download := m.dataStore.Downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.PauseDownload(download)
		} else if state == "paused" {
			// Resume the download
			m.downloadsTable.Rows()[m.selectedRow][3] = "pending" // Update the state to "pending"
			download := m.dataStore.Downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.ResumeDownload(download)
		}
	}
}

// Modify the removeDownload method to remove the selected row
func (m *Model) removeDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Remove the row from the table by slicing the rows

		index := m.downloadsTable.Rows()[m.selectedRow][0]
		download := m.dataStore.Downloads[index]
		m.downloadmanager.RemoveDownload(download)
		m.dataStore.RemoveDownload(download)

		newRows := append(m.downloadsTable.Rows()[:m.selectedRow], m.downloadsTable.Rows()[m.selectedRow+1:]...)

		// Update the downloadsTable with the new rows
		m.downloadsTable = table.New(
			table.WithColumns(downloadColumns), // Keep the existing columns
			table.WithRows(newRows),            // Set the new rows
		)

		// Adjust the selected row to prevent out of bounds error if the last row is removed
		if m.selectedRow >= len(newRows) {
			m.selectedRow = len(newRows) - 1
		}
	}
}

func (m *Model) removeQueue() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.queuesTable.Rows()) {
		// Remove the row from the table by slicing the rows

		index := m.queuesTable.Rows()[m.selectedRow][0]
		queue := m.dataStore.Queues[index]
		m.downloadmanager.RemoveQueue(queue)
		for _, download := range queue.Downloads {
			m.dataStore.RemoveDownload(download)
		}
		m.dataStore.RemoveQueue(queue)

		// Update the queuesTable with the new rows
		newRows := append(m.queuesTable.Rows()[:m.selectedRow], m.queuesTable.Rows()[m.selectedRow+1:]...)
		m.queuesTable = table.New(
			table.WithColumns(queueColumns), // Keep the existing columns
			table.WithRows(newRows),         // Set the new rows
		)

		// Adjust the selected row to prevent out of bounds error if the last row is removed
		if m.selectedRow >= len(newRows) {
			m.selectedRow = len(newRows) - 1
		}
		m.selectedQueueRowIndex = 0
	}
}

// Add a method to handle Retry action (only if the state is "Failed")
func (m *Model) retryDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Check the state of the selected row
		state := m.downloadsTable.Rows()[m.selectedRow][3]

		if state == "failed" {
			// Retry the download
			m.downloadsTable.Rows()[m.selectedRow][3] = "retrying" // Update status to "Retrying"
			download := m.dataStore.Downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.RetryDownload(download)
		}
	}
}

func (m *Model) updateFocusedField(msg tea.Msg) {
	if m.focusedField == 0 {
		m.inputURL.Update(msg)
		// } else if m.focusedField == 1 {
		// 	m.queueSelect.Update(msg)
	} else if m.focusedField == 2 {
		m.outputFileName.Update(msg)
	}
}

// Handle the submission of a new queue form
func (m *Model) handleNewOrEditQueueFormSubmit() {
	if m.newQueueForm || m.editQueueForm {

		errorCode := m.CheckErrorCodes()
		if errorCode == 1 {
			m.handleCancel()
			return
		}
		MaxConcurrentDownloads, err := strconv.Atoi(m.maxConcurrentInput.Value())
		if err != nil {
			MaxConcurrentDownloads = 10
		} else if MaxConcurrentDownloads < 0 {
			MaxConcurrentDownloads = 1
		}
		MaxBandwidth, _ := strconv.Atoi(m.maxBandwidthInput.Value())
		if err != nil {
			MaxBandwidth = 0
		} else if MaxBandwidth < 0 {
			MaxBandwidth = 0
		}
		MaxRetries, _ := strconv.Atoi(m.maxRetriesPerDLInput.Value())
		if err != nil {
			MaxRetries = 5
		} else if MaxRetries < 0 {
			MaxRetries = 0
		}
		if m.activeStartTimeInput.Value() == "" {
			m.activeStartTimeInput.SetValue("00:00")
		}

		if m.activeEndTimeInput.Value() == "" {
			m.activeEndTimeInput.SetValue("23:59")
		}

		if m.editQueueForm {
			if m.selectedRow >= 0 && m.selectedRow < len(m.queuesTable.Rows()) {
				oldQueueRow := m.queuesTable.Rows()[m.selectedRow]

				thisQueue := m.dataStore.Queues[oldQueueRow[0]]
				thisQueue.SaveDir = m.saveDirInput.Value()
				thisQueue.MaxConcurrentDownloads = MaxConcurrentDownloads
				thisQueue.MaxRetries = MaxRetries

				thisQueue.ActiveStartTime = m.activeStartTimeInput.Value()
				thisQueue.ActiveEndTime = m.activeEndTimeInput.Value()

				if thisQueue.MaxBandwidth != MaxBandwidth {
					thisQueue.SetBandwith(MaxBandwidth)
				}

				m.editQueue(oldQueueRow, thisQueue)
				m.newQueueForm = false
				m.editQueueForm = false
				m.showEditQConfirmation()
			}
		} else {
			m.maxQueueID++
			newQueue := manager.Queue{
				ID:                     m.maxQueueID,
				SaveDir:                m.saveDirInput.Value(),
				MaxConcurrentDownloads: MaxConcurrentDownloads,
				MaxBandwidth:           MaxBandwidth,
				MaxRetries:             MaxRetries,
				ActiveStartTime:        m.activeStartTimeInput.Value(),
				ActiveEndTime:          m.activeEndTimeInput.Value(),
			}
			// Adding a new queue
			m.addNewQueue(&newQueue)
			m.newQueueForm = false
			m.editQueueForm = false
			m.showAddQConfirmation()
		}

		// Reset the form after submission

		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
		m.maxRetriesPerDLInput.Reset()
		m.activeStartTimeInput.Reset()
		m.activeEndTimeInput.Reset()
		m.saveDirInput.Focus()
		m.maxConcurrentInput.Blur()
		m.maxBandwidthInput.Blur()
		m.maxRetriesPerDLInput.Blur()
		m.activeStartTimeInput.Blur()
		m.activeEndTimeInput.Blur()

		m.focusedFieldForQueues = 0
		counterForForms = 0
	}
}

// Add a new queue
func (m *Model) addNewQueue(queue *manager.Queue) {
	m.dataStore.AddQueue(queue)
	m.downloadmanager.AddQueue(queue)

	newRow := table.Row{
		strconv.Itoa(queue.ID),
		queue.SaveDir,
		strconv.Itoa(queue.MaxConcurrentDownloads),
		strconv.Itoa(queue.MaxBandwidth),
		strconv.Itoa(queue.MaxRetries),
		queue.ActiveStartTime,
		queue.ActiveEndTime,
	}

	// Add the row to the queuesTable
	m.queuesTable = table.New(
		table.WithColumns(queueColumns),                                      // Keep the existing columns
		table.WithRows(append([]table.Row{newRow}, m.queuesTable.Rows()...)), // Add the new row
	)
}

// Edit an existing queue
func (m *Model) editQueue(oldQueueRow table.Row, queue *manager.Queue) {
	m.dataStore.Save()
	// Update the selected queue with new values
	m.queuesTable.Rows()[m.selectedRow] = table.Row{
		oldQueueRow[0],
		queue.SaveDir,
		strconv.Itoa(queue.MaxConcurrentDownloads),
		strconv.Itoa(queue.MaxBandwidth),
		strconv.Itoa(queue.MaxRetries),
		queue.ActiveStartTime,
		queue.ActiveEndTime,
	}

	// Update the table
	m.queuesTable = table.New(
		table.WithColumns(queueColumns),
		table.WithRows(m.queuesTable.Rows()),
	)

}

func (m *Model) handleCancel() {
	if m.newQueueForm {
		// If adding a new queue, cancel and reset the form
		m.newQueueForm = false
		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
		m.maxRetriesPerDLInput.Reset()
		m.activeStartTimeInput.Reset()
		m.activeEndTimeInput.Reset()
		counterForForms = 0
		m.focusedFieldForQueues = 0
	}
	if m.editQueueForm {
		// If editing a queue, cancel the edit and return to the queue list
		m.editQueueForm = false
		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
		m.maxRetriesPerDLInput.Reset()
		m.activeStartTimeInput.Reset()
		m.activeEndTimeInput.Reset()
		counterForForms = 0
		m.focusedFieldForQueues = 0
	}
	// Ensure we are in the "Queues" tab and re-render it
	m.currentTab = tabQueues
}

func (m *Model) handleUpArrowForTab3() {
	if m.selectedRow > 0 {
		m.selectedRow-- // Navigate up in the queue list
	}
}

func (m *Model) handleDownArrowForTab3() {
	if m.selectedRow < len(m.queuesTable.Rows())-1 {
		m.selectedRow++ // Navigate down in the queue list
	}
}

func (m *Model) handleSwitchToAddQueueForm() {
	if m.currentTab == tabQueues {
		if !m.newQueueForm && !m.editQueueForm {
			counterForForms = counterForForms + 1
			m.newQueueForm = true
			m.editQueueForm = false
			m.updateFocusedFieldForTab3()
		}
	}
}

func (m *Model) handleSwitchToEditQueueForm() {
	if m.currentTab == tabQueues && m.selectedRow >= 0 {
		if !m.newQueueForm && !m.editQueueForm {
			counterForForms = counterForForms + 1

			m.editQueueForm = true
			m.newQueueForm = false
			queueID := m.queuesTable.Rows()[m.selectedRow][0]
			thisQueue := m.dataStore.Queues[queueID]
			m.saveDirInput.SetValue(thisQueue.SaveDir)
			m.maxConcurrentInput.SetValue(strconv.Itoa(thisQueue.MaxConcurrentDownloads))
			m.maxBandwidthInput.SetValue(strconv.Itoa(thisQueue.MaxBandwidth))
			m.maxRetriesPerDLInput.SetValue(strconv.Itoa(thisQueue.MaxRetries))

			// queue[0]
			m.updateFocusedFieldForTab3()
		}
	}
}

func (m *Model) updateBasedOnInputForTab1(msg tea.Msg, _ tea.Cmd) {
	// Update the text inputs based on focus
	if m.currentTab == tabAddDownload {
		if m.focusedField == 0 {
			m.inputURL, _ = m.inputURL.Update(msg)
			// } else if m.focusedField == 1 {
			// m.selectedQueueRowIndex = 0
		} else if m.focusedField == 2 {
			m.outputFileName, _ = m.outputFileName.Update(msg)
		}
		// Update the focused field accordingly
		m.updateFocusedField(msg)
	}
}

func (m *Model) updateBasedOnInputForTab3(msg tea.Msg, _ tea.Cmd) {
	if m.newQueueForm || m.editQueueForm {
		counterForForms = counterForForms + 1
		switch m.focusedFieldForQueues {
		case 0:
			m.saveDirInput, _ = m.saveDirInput.Update(msg)
		case 1:
			m.maxConcurrentInput, _ = m.maxConcurrentInput.Update(msg)
		case 2:
			m.maxBandwidthInput, _ = m.maxBandwidthInput.Update(msg)
		case 3:
			m.maxRetriesPerDLInput, _ = m.maxRetriesPerDLInput.Update(msg)
		case 4:
			m.activeStartTimeInput, _ = m.activeStartTimeInput.Update(msg)
		case 5:
			m.activeEndTimeInput, _ = m.activeEndTimeInput.Update(msg)
		}
	}
}
