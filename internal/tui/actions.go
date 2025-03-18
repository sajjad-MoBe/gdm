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

// Handle the submission of a new download form
func (m *Model) handleNewDownloadSubmit() {
	if m.inputURL.Value() == "" {
		m.showURLValidationError()
	} else {
		// Create a new download with the data entered in the fields
		downloadURL := m.inputURL.Value()
		// validate download url

		outputFile := m.outputFileName.Value()
		// should be validated
		queue := m.queues[m.queuesTable.Rows()[m.selectedQueueRowIndex][0]]

		newDwnload := manager.Download{
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
	if err := manager.Create(&download); err != nil {
		//	show error
		fmt.Printf("failed to create download: %v", err)
		return
	}
	m.downloadmanager.AddDownload(download)
	m.downloads[strconv.Itoa(download.ID)] = download

	newRow := table.Row{
		strconv.Itoa(download.ID),
		download.URL,
		download.Status,
		"N/A",
		"N/A",
	}

	// Add the row to the downloadsTable
	m.downloadsTable = table.New(
		table.WithColumns(downloadColumns),                      // Keep the existing columns
		table.WithRows(append(m.downloadsTable.Rows(), newRow)), // Add the new row
	)
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
	m.confirmationMessage = "Download has been added, press right button to show Downloads!"
	m.confirmationTime = time.Now()

	m.resetFieldsForTab1()
	m.focusedField = 0
	m.inputURL.Focus()
}

func (m *Model) resetFieldsForTab1() {
	m.inputURL.SetValue("")
	m.queueSelect.ResetSelected()
	m.outputFileName.SetValue("")
	m.selectedQueueRowIndex = 0
}

// unused
// func (m *Model) resetFieldsForTab3() {
// 	m.saveDirInput.SetValue("")
// 	m.maxConcurrentInput.SetValue("")
// 	m.maxBandwidthInput.SetValue("")
// }

func (m *Model) handleUpArrowForTab1() {
	if m.focusedField == 1 {
		m.selectedQueueRowIndex = (m.selectedQueueRowIndex - 1) % len(m.queuesTable.Rows())
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

func (m *Model) handleSpaceKey() {
	if m.currentTab == tabAddDownload && m.focusedField == 1 {
		m.selectedQueueRowIndex++
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

		if state == "downloading" {
			// Pause the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "paused" // Update the state to "Paused"
			download := m.downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.PauseDownload(download)
		} else if state == "paused" {
			// Resume the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "pending" // Update the state to "pending"
			download := m.downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.ResumeDownload(download)
		}
	}
}

// Modify the deleteDownload method to delete the selected row
func (m *Model) deleteDownload() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.downloadsTable.Rows()) {
		// Remove the row from the table by slicing the rows
		newRows := append(m.downloadsTable.Rows()[:m.selectedRow], m.downloadsTable.Rows()[m.selectedRow+1:]...)

		index := m.downloadsTable.Rows()[m.selectedRow][0]
		download := m.downloads[index]
		m.downloadmanager.DeleteDownload(download)
		manager.Delete(download)
		delete(m.downloads, index)

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

func (m *Model) deleteQueue() {
	if m.selectedRow >= 0 && m.selectedRow < len(m.queuesTable.Rows()) {
		// Remove the row from the table by slicing the rows
		newRows := append(m.queuesTable.Rows()[:m.selectedRow], m.queuesTable.Rows()[m.selectedRow+1:]...)

		index := m.queuesTable.Rows()[m.selectedRow][0]
		queue := m.queues[index]
		m.downloadmanager.DeleteQueue(queue)
		manager.Delete(queue)
		delete(m.queues, index)

		// Update the queuesTable with the new rows
		m.queuesTable = table.New(
			table.WithColumns(queueColumns), // Keep the existing columns
			table.WithRows(newRows),         // Set the new rows
		)

		// Adjust the selected row to prevent out of bounds error if the last row is deleted
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
		state := m.downloadsTable.Rows()[m.selectedRow][2]

		if state == "failed" {
			// Retry the download
			m.downloadsTable.Rows()[m.selectedRow][2] = "retrying" // Update status to "Retrying"
			download := m.downloads[m.downloadsTable.Rows()[m.selectedRow][0]]
			m.downloadmanager.RetryDownload(download)
		}
	}
}

func (m *Model) updateFocusedField(msg tea.Msg) {
	if m.focusedField == 0 {
		m.inputURL.Update(msg)
	} else if m.focusedField == 1 {
		m.queueSelect.Update(msg)
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

		if m.editQueueForm {
			if m.selectedRow >= 0 && m.selectedRow < len(m.queuesTable.Rows()) {
				oldQueue := m.queuesTable.Rows()[m.selectedRow]

				thisQueue := m.queues[oldQueue[0]]
				thisQueue.SaveDir = m.saveDirInput.Value()
				thisQueue.MaxConcurrentDownloads = MaxConcurrentDownloads
				if thisQueue.MaxBandwidth != MaxBandwidth {
					thisQueue.SetBandwith(MaxBandwidth)
				}
				// Editing this existing queue

				m.editQueue(oldQueue, thisQueue)

			}
		} else {
			newQueue := manager.Queue{
				SaveDir:                m.saveDirInput.Value(),
				MaxConcurrentDownloads: MaxConcurrentDownloads,
				MaxBandwidth:           MaxBandwidth,
			}
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
		return
	}
	m.downloadmanager.AddQueue(queue)
	m.queues[strconv.Itoa(queue.ID)] = queue

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
func (m *Model) editQueue(oldQueue table.Row, queue *manager.Queue) {
	manager.Save(queue)
	// Update the selected queue with new values
	m.queuesTable.Rows()[m.selectedRow] = table.Row{
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
			thisQueue := m.queues[queueID]
			m.saveDirInput.SetValue(thisQueue.SaveDir)
			m.maxConcurrentInput.SetValue(strconv.Itoa(thisQueue.MaxConcurrentDownloads))
			m.maxBandwidthInput.SetValue(strconv.Itoa(thisQueue.MaxBandwidth))

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
			m.queueSelect, _ = m.queueSelect.Update(msg)
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
