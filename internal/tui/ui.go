package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	app       *tview.Application
	pages     *tview.Pages
	tabBar    *tview.Flex
	addForm   *tview.Form
	downloads *tview.Table
	queues    *tview.Table
	helpBar   *tview.TextView
}

func NewApp() *App {
	a := &App{
		app:       tview.NewApplication(),
		pages:     tview.NewPages(),
		addForm:   tview.NewForm(),
		downloads: tview.NewTable(),
		queues:    tview.NewTable(),
		helpBar:   tview.NewTextView(),
	}

	a.setupUI()
	return a
}

func (a *App) setupUI() {
	a.helpBar.SetText("F1:Add F2:Downloads F3:Queues A:AddDownload D:Delete E:Edit")
	a.helpBar.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)

	a.setupAddTab()
	a.setupDownloadsTab()
	a.setupQueuesTab()

	a.tabBar = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(a.pages, 0, 1, true).
		AddItem(a.helpBar, 1, 1, false)

	a.pages.AddPage("Add Download", a.addForm, true, true)
	a.pages.AddPage("Downloads List", a.downloads, true, false)
	a.pages.AddPage("Queues List", a.queues, true, false)

	a.setupKeyBindings()
}

func (a *App) setupAddTab() {
	a.addForm.
		AddInputField("URL", "", 50, nil, nil).
		AddInputField("Output Name", "", 50, nil, nil).
		AddButton("OK", func() {
			// Handle adding a download (To be implemented)
		}).
		AddButton("Cancel", func() {
			a.switchToTab("Downloads List")
		})
}

func (a *App) setupDownloadsTab() {
	a.downloads.SetBorders(true).SetTitle("Downloads").SetTitleAlign(tview.AlignLeft)

	headers := []string{"URL", "Queue", "Status", "Speed", "Progress"}
	for i, h := range headers {
		a.downloads.SetCell(0, i, tview.NewTableCell(h).SetTextColor(tcell.ColorGhostWhite).SetSelectable(false))
	}
}

func (a *App) setupQueuesTab() {
	a.queues.SetBorders(true).SetTitle("Queues").SetTitleAlign(tview.AlignLeft)

	headers := []string{"Queue Name", "Folder", "Max Downloads", "Speed Limit", "Schedule"}
	for i, h := range headers {
		a.queues.SetCell(0, i, tview.NewTableCell(h).SetTextColor(tcell.ColorGhostWhite).SetSelectable(false))
	}
}

func (a *App) setupKeyBindings() {
	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			a.switchToTab("Add Download")
		case tcell.KeyF2:
			a.switchToTab("Downloads List")
		case tcell.KeyF3:
			a.switchToTab("Queues List")
		case tcell.KeyRune:
			switch event.Rune() {
			case 'a':
				a.switchToTab("Add Download")
			case 'd':
				// Handle delete action
			case 'e':
				// Handle edit action
			}
		}
		return event
	})
}

func (a *App) switchToTab(tabName string) {
	a.pages.SwitchToPage(tabName)
}

func (a *App) Run() error {
	return a.app.SetRoot(a.tabBar, true).Run()
}
