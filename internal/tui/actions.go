package tui

import (
	"github.com/rivo/tview"
)

func (a *App) AddDownload(url, queue, outputName string) {
	row := a.downloads.GetRowCount()
	a.downloads.SetCell(row, 0, tview.NewTableCell(url))
	a.downloads.SetCell(row, 1, tview.NewTableCell(queue))
	a.downloads.SetCell(row, 2, tview.NewTableCell("Queued"))
	a.downloads.SetCell(row, 3, tview.NewTableCell("0 KB/s"))
	a.downloads.SetCell(row, 4, tview.NewTableCell("0%"))
}
