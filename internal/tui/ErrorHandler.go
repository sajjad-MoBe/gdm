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

func (m *Model) handleRetriesError() {
	m.errorMessage = "Invalid Max Retries Input!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleAllErrors() {
	// Logic for when all errors are present
	m.errorMessage = "Max Concurrent, Max Bandwidth, Max Retries ,and Time inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleBWAndConcurrentAndTimeErrors() {
	m.errorMessage = "Max Concurrent, Max Bandwidth, and Time inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleRetriesAndConcurrentAndTimeErrors() {
	m.errorMessage = "Max Concurrent, Max Retries, and Time inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleRetriesAndConcurrentAndBWErrors() {
	m.errorMessage = "Max Concurrent, Max Bandwidth, and Max Retries inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleRetriesAndBWAndTimeErrors() {
	m.errorMessage = "Max Retries, Max Bandwidth, and Time inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleBWAndTimeErrors() {
	// Logic for bandwidth + time errors
	m.errorMessage = "Max Bandwidth and Time format inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleBWAndConcurrentErrors() {
	m.errorMessage = "Max Bandwidth and Max Concurrent inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleConcurrentAndTimeErrors() {
	// Logic for concurrent + time errors
	m.errorMessage = "Max Concurrent and Time format inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleConcurrentAndRetriesErrors() {
	// Logic for concurrent + time errors
	m.errorMessage = "Max Concurrent and Max Retries inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleTimeAndRetriesErrors() {
	// Logic for concurrent + time errors
	m.errorMessage = "Max Retries and Time format inputs are invalid!"
	m.confirmationMessage = ""
	m.setupsAfterErrorForQueues()
}

func (m *Model) handleBWAndRetriesErrors() {
	// Logic for concurrent + time errors
	m.errorMessage = "Max BandWidth and Max Retries inputs are invalid!"
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
	var concurrentError, bwError, retriesError, timeError bool

	if !regForConcurrent.MatchString(m.maxConcurrentInput.Value()) {
		concurrentError = true
	}
	if !regForMaxBW.MatchString(m.maxBandwidthInput.Value()) {
		bwError = true
	}
	if !regForMaxBW.MatchString(m.maxRetriesPerDLInput.Value()) {
		retriesError = true
	}
	if !regForHHMMFormat.MatchString(m.activeStartTimeInput.Value()) || !regForHHMMFormat.MatchString(m.activeEndTimeInput.Value()) {
		timeError = true
	}

	if timeError && concurrentError && bwError && retriesError {
		m.handleAllErrors()
		return 1
	} else if timeError && concurrentError && bwError {
		m.handleBWAndConcurrentAndTimeErrors()

		return 1
	} else if timeError && concurrentError && retriesError {
		m.handleRetriesAndConcurrentAndTimeErrors()
		return 1
	} else if retriesError && concurrentError && bwError {
		m.handleRetriesAndConcurrentAndBWErrors()
		return 1
	} else if timeError && bwError && retriesError {
		m.handleRetriesAndBWAndTimeErrors()
		return 1
	} else if timeError && concurrentError {
		m.handleConcurrentAndTimeErrors()
		return 1
	} else if bwError && retriesError {
		m.handleBWAndRetriesErrors()
		return 1
	} else if timeError && bwError {
		m.handleBWAndTimeErrors()
		return 1
	} else if timeError && retriesError {
		m.handleTimeAndRetriesErrors()
		return 1
	} else if bwError && concurrentError {
		m.handleBWAndConcurrentErrors()
		return 1
	} else if retriesError && concurrentError {
		m.handleConcurrentAndRetriesErrors()
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
	} else if retriesError {
		m.handleRetriesError()
		return 1
	}
	return 0
}
