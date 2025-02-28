package tui

func (a *App) AddDownload(url, queue, outputName string) {
	a.downloads = append(a.downloads, Download{
		URL:      url,
		Queue:    queue,
		Status:   "Queued",
		Speed:    "0 KB/s",
		Progress: "0%",
	})
}
