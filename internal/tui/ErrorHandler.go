package tui

import "time"

func (m *Model) showURLValidationError() {
	m.errorMessage = "URL is required!"
	m.confirmationMessage = ""
	m.errorTime = time.Now()

	m.resetFieldsForTab1()
	m.focusedField = 0
	m.inputURL.Focus()
	m.outputFileName.Blur()
}

func (m *Model) handleBWError() {
	m.errorMessage = "Invalid Max Bandwidth Input!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleTimeError() {
	m.errorMessage = "Invalid Time Input!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleConcurrentError() {
	m.errorMessage = "Invalid Max Concurrent Input!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleAllErrors() {
	// Logic for when all errors are present
	m.errorMessage = "Max Concurrent, Max Bandwidth, and Time inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleConcurrentAndBWErrors() {
	// Logic for concurrent + bandwidth errors
	m.errorMessage = "Max Concurrent and Max Bandwidth inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleBWAndTimeErrors() {
	// Logic for bandwidth + time errors
	m.errorMessage = "Max Bandwidth and Time format inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleConcurrentAndTimeErrors() {
	// Logic for concurrent + time errors
	m.errorMessage = "Max Concurrent and Time format inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) setupsAfterErrorForQueues() {
	m.errorTime = time.Now()
	m.resetFieldsForTab3()
	m.focusedFieldForQueues = 0
	m.saveDirInput.Focus()
	m.maxConcurrentInput.Blur()
	m.maxConcurrentInput.Blur()
	m.activeStartTimeInput.Blur()
	m.activeEndTimeInput.Blur()
}

func (m *Model) showDownloadConfirmation() {
	m.confirmationMessage = "Download has been added!"
	m.confirmationTime = time.Now()

	m.resetFieldsForTab1()
	m.focusedField = 0
	m.inputURL.Focus()
	m.outputFileName.Blur()
}

func (m *Model) showAddQConfirmation() {
	m.confirmationMessage = "Queue has been added successfully!"
	m.confirmationTime = time.Now()

	m.resetFieldsForTab3()
	m.focusedFieldForQueues = 0
	m.saveDirInput.Focus()
	m.maxConcurrentInput.Blur()
	m.maxConcurrentInput.Blur()
	m.activeStartTimeInput.Blur()
	m.activeEndTimeInput.Blur()
}

func (m *Model) showEditQConfirmation() {
	m.confirmationMessage = "Queue has been edited successfully!"
	m.confirmationTime = time.Now()

	m.resetFieldsForTab3()
	m.focusedFieldForQueues = 0
	m.saveDirInput.Focus()
	m.maxConcurrentInput.Blur()
	m.maxConcurrentInput.Blur()
	m.activeStartTimeInput.Blur()
	m.activeEndTimeInput.Blur()
}

func (m *Model) CheckErrorCodes() int {
	var concurrentError, bwError, timeError bool

	if !regForConcurrent.MatchString(m.maxConcurrentInput.Value()) {
		concurrentError = true
	}
	if !regForMaxBW.MatchString(m.maxBandwidthInput.Value()) {
		bwError = true
	}
	if !regForHHMMFormat.MatchString(m.activeStartTimeInput.Value()) || !regForHHMMFormat.MatchString(m.activeEndTimeInput.Value()) {
		timeError = true
	}
	if timeError && concurrentError && bwError {
		m.handleAllErrors()
		return 1
	} else if timeError && concurrentError {
		m.handleConcurrentAndTimeErrors()
		return 1
	} else if bwError && concurrentError {
		m.handleConcurrentAndBWErrors()
		return 1
	} else if bwError && timeError {
		m.handleBWAndTimeErrors()
		return 1
	} else if timeError {
		m.handleTimeError()
		return 1
	} else if bwError {
		m.handleBWError()
		return 1
	} else if concurrentError {
		m.handleConcurrentError()
		return 1
	}
	return 0
}
