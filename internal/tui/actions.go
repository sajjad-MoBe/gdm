package tui

import (
	"fmt"
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
	if m.currentTab < tabQueues && !m.newQueueForm && !m.editQueueForm {
		m.currentTab++
		m.selectedRow = 0
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

	m.resetFieldsForTab1()
	m.focusedField = 0
	m.inputURL.Focus()
}

func (m *Model) showDownloadConfirmation() {
	m.confirmationMessage = "Download has been added!"
	m.confirmationTime = time.Now()

	m.resetFieldsForTab1()
	m.focusedField = 0
	m.inputURL.Focus()
}

func (m *Model) resetFieldsForTab1() {
	m.inputURL.SetValue("")
	m.pageSelect.ResetSelected()
	m.outputFileName.SetValue("")
	m.selectedFiles = make(map[int]struct{})
}

// unused
// func (m *Model) resetFieldsForTab3() {
// 	m.saveDirInput.SetValue("")
// 	m.maxConcurrentInput.SetValue("")
// 	m.maxBandwidthInput.SetValue("")
// }

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

func (m *Model) updateFocusedFieldForTab1() {
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

// Update focus for the queue form fields
func (m *Model) updateFocusedFieldForTab3() {
	switch m.focusedFieldForQueues {
	case 0:
		m.saveDirInput.Focus()
		m.maxConcurrentInput.Blur()
		m.maxBandwidthInput.Blur()
	case 1:
		m.saveDirInput.Blur()
		m.maxConcurrentInput.Focus()
		m.maxBandwidthInput.Blur()
	case 2:
		m.saveDirInput.Blur()
		m.maxConcurrentInput.Blur()
		m.maxBandwidthInput.Focus()
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
		state := m.downloadsTable.Rows()[m.selectedRow][2]

		if state == "Downloading" {
			// Pause the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "Paused" // Update the state to "Paused"
		} else if state == "Paused" {
			// Resume the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "Downloading" // Update the state to "Downloading"
		}
	}
}

// Modify the deleteDownload method to delete the selected row
func (m *Model) deleteDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Remove the row from the table by slicing the rows
		newRows := append(m.downloadsTable.Rows()[:m.selectedRow], m.downloadsTable.Rows()[m.selectedRow+1:]...)

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
func (m *Model) retryDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Check the state of the selected row
		state := m.downloadsTable.Rows()[m.selectedRow][2]

		if state == "Failed" {
			// Retry the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "Retrying" // Update status to "Retrying"
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

// unused
// func (m *Model) updateFocusedFieldForQueue(msg tea.Msg) {
// 	if m.focusedFieldForQueues == 0 {
// 		m.saveDirInput.Update(msg)
// 	} else if m.focusedFieldForQueues == 1 {
// 		m.maxConcurrentInput.Update(msg)
// 	} else if m.focusedFieldForQueues == 2 {
// 		m.maxBandwidthInput.Update(msg)
// 	}
// }

// Handle the submission of a new queue form
func (m *Model) handleNewOrEditQueueFormSubmit() {
	if m.newQueueForm || m.editQueueForm {
		// Create a new QueueItem with the data entered in the fields
		MaxConcurrentDownloads, err := strconv.Atoi(m.maxConcurrentInput.Value())
		if err != nil {
			// show error
			fmt.Println("Error:", err)
			return
		}
		MaxBandwidth, err := strconv.Atoi(m.maxBandwidthInput.Value())
		if err != nil {
			// show error
			fmt.Println("Error:", err)
			return
		}
		newQueue := manager.Queue{
			SaveDir:                m.saveDirInput.Value(),
			MaxConcurrentDownloads: MaxConcurrentDownloads,
			MaxBandwidth:           MaxBandwidth,
		}

		if m.editQueueForm {
			// Editing an existing queue
			m.editQueue(m.selectedRow, &newQueue)
		} else {
			// Adding a new queue
			m.addNewQueue(&newQueue)
		}

		// Reset the form after submission
		m.newQueueForm = false
		m.editQueueForm = false
		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
		m.saveDirInput.Focus()
		m.maxConcurrentInput.Blur()
		m.maxBandwidthInput.Blur()

		m.focusedFieldForQueues = 0
		counterForForms = 0
	}
}

// Add a new queue
func (m *Model) addNewQueue(queue *manager.Queue) {
	if err := manager.Create(&queue); err != nil {
		//	show error
		fmt.Printf("failed to create queue: %v", err)
	}
	newRow := table.Row{
		strconv.Itoa(queue.ID),
		strconv.Itoa(queue.MaxConcurrentDownloads),
		strconv.Itoa(queue.MaxBandwidth),
		queue.SaveDir,
		queue.ActiveStartTime,
		queue.ActiveEndTime,
	}

	// Add the row to the queuesTable
	m.queuesTable = table.New(
		table.WithColumns(queueColumns),                      // Keep the existing columns
		table.WithRows(append(m.queuesTable.Rows(), newRow)), // Add the new row
	)
}

// Edit an existing queue
func (m *Model) editQueue(index int, queue *manager.Queue) {
	if index >= 0 && index < len(m.queuesTable.Rows()) {
		oldQueue := m.queuesTable.Rows()[index]
		// Update the selected queue with new values
		m.queuesTable.Rows()[index] = table.Row{
			oldQueue[0],
			queue.SaveDir,
			strconv.Itoa(queue.MaxConcurrentDownloads),
			strconv.Itoa(queue.MaxBandwidth),
			oldQueue[4],
			oldQueue[5],
		}

		// Update the table
		m.queuesTable = table.New(
			table.WithColumns(queueColumns),
			table.WithRows(m.queuesTable.Rows()),
		)
	}
}

func (m *Model) handleCancel() {
	if m.newQueueForm {
		// If adding a new queue, cancel and reset the form
		m.newQueueForm = false
		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
		counterForForms = 0
		m.focusedFieldForQueues = 0
	}
	if m.editQueueForm {
		// If editing a queue, cancel the edit and return to the queue list
		m.editQueueForm = false
		m.saveDirInput.Reset()
		m.maxConcurrentInput.Reset()
		m.maxBandwidthInput.Reset()
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
			m.newQueueData = &manager.Queue{}
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
			m.newQueueData = m.queues[queueID]
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
		} else if m.focusedField == 1 {
			m.pageSelect, _ = m.pageSelect.Update(msg)
		} else if m.focusedField == 2 {
			m.outputFileName, _ = m.outputFileName.Update(msg)
		}
		// Clear messages if necessary
		m.clearMessages()
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
		}
		m.updateFocusedFieldForTab1()
	}
}
