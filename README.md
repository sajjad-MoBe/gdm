# gdm

A terminal-based **Download Manager** built with **Golang**, designed to help users manage file downloads efficiently using a **Text-based User Interface (TUI)**. This project demonstrates various aspects of Golang, including concurrency, error handling, and working with external libraries like `tview` for creating interactive terminal UIs.

---

## Features

### Core Functionality
- **Download Management**
  - Add, pause, resume, retry, and cancel downloads.
  - View download progress, speed, and status (Pending, Downloading, Paused, Completed, Failed).
  - Support for parallel downloads using goroutines.
  - Capable of multi-part downloads for large files, leveraging server support for `Accept-Ranges` headers.

- **Queue-Based Downloading**
  - Organize downloads into multiple queues.
  - Customize queue settings:
    - **Save folder** for downloaded files.
    - **Maximum concurrent downloads** per queue.
    - **Speed limit** for downloads (e.g., 500 KB/s).
    - **Active time range** for scheduling downloads (e.g., 10 PM to 6 AM).
    - **Retry attempts** for failed downloads.

### Text-Based User Interface (TUI)
- Built using the [`tview`](https://github.com/rivo/tview) library.
- **Three main tabs:**
  1. **Add Download**: A form to add new downloads (URL, queue selection, file name).
  2. **Downloads List**: Displays all downloads, including their status, progress, and speed.
  3. **Queues Management**: Manage download queues, edit settings, or create/delete queues.
- Keyboard shortcuts for navigation and actions.
- Footer bar displaying helpful key bindings.

### Persistence
- Save and restore the state of downloads and queues when the application is closed and reopened.
- Store configuration and download information in a sqlite3 database.

---

## Getting Started

### Prerequisites
- **Golang** (version 1.18 or later) installed on your system.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/sajjad-mobe/gdm.git
   cd gdm