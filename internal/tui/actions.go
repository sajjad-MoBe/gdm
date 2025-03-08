package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AddDownloadMsg struct {
	URL      string
	Queue    string
	Filename string
}

type DeleteDownloadMsg struct {
	Index int
}

type PauseResumeDownloadMsg struct {
	Index int
}

type RetryDownloadMsg struct {
	Index int
}

type AddQueueMsg struct {
	Name            string
	Folder          string
	MaxConcurrent   int
	SpeedLimit      int
	TimeRestriction string
}

type DeleteQueueMsg struct {
	Index int
}

type EditQueueMsg struct {
	Index              int
	NewName            string
	NewFolder          string
	NewMaxConcurrent   int
	NewSpeedLimit      int
	NewTimeRestriction string
}

func AddDownload(url, queue, filename string) tea.Cmd {
	return func() tea.Msg {
		return AddDownloadMsg{URL: url, Queue: queue, Filename: filename}
	}
}

func DeleteDownload(index int) tea.Cmd {
	return func() tea.Msg {
		return DeleteDownloadMsg{Index: index}
	}
}

func PauseResumeDownload(index int) tea.Cmd {
	return func() tea.Msg {
		return PauseResumeDownloadMsg{Index: index}
	}
}

func RetryDownload(index int) tea.Cmd {
	return func() tea.Msg {
		return RetryDownloadMsg{Index: index}
	}
}

func AddQueue(name, folder string, maxConcurrent, speedLimit int, timeRestriction string) tea.Cmd {
	return func() tea.Msg {
		return AddQueueMsg{Name: name, Folder: folder, MaxConcurrent: maxConcurrent, SpeedLimit: speedLimit, TimeRestriction: timeRestriction}
	}
}

func DeleteQueue(index int) tea.Cmd {
	return func() tea.Msg {
		return DeleteQueueMsg{Index: index}
	}
}

func EditQueue(index int, newName, newFolder string, newMaxConcurrent, newSpeedLimit int, newTimeRestriction string) tea.Cmd {
	return func() tea.Msg {
		return EditQueueMsg{Index: index, NewName: newName, NewFolder: newFolder, NewMaxConcurrent: newMaxConcurrent, NewSpeedLimit: newSpeedLimit, NewTimeRestriction: newTimeRestriction}
	}
}
