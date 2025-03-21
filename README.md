# gdm

A terminal-based **Download Manager** built with **Golang**, designed to help users manage file downloads efficiently using a **Text-based User Interface (TUI)**. This project demonstrates various aspects of Golang, including concurrency, error handling, and working with external libraries like `Bubble Tea` for creating interactive terminal UIs.

---

## Features

### Core Functionality
- **Download Management**
  - Add, pause, resume, retry, and cancel downloads.
  - View download progress, speed, and status (initializing, pending, downloading, paused, downloaded, failed).
  - Support for parallel downloads using goroutines.
  - Capable of multi-part downloads for large files, leveraging server support for `Accept-Ranges` headers.

- **Queue-Based Downloading**
  - Organize downloads into multiple queues.
  - Customize queue settings:
    - **Save folder** for downloaded files.
    - **Maximum concurrent downloads** per queue.
    - **Speed limit** for downloads (e.g., 500 KB/s).
    - **Active time range** for scheduling downloads (e.g., 10:10 to 20:30).
    - **Retry attempts** for failed downloads.

### Text-Based User Interface (TUI)
- Built using the [`Bubble Tea`](https://github.com/charmbracelet/bubbletea) library.
- **Three main tabs:**
  1. **Add Download**: A form to add new downloads (URL, queue selection, file name).
  2. **Downloads List**: Displays all downloads, including their status, progress, and speed.
  3. **Queues Management**: Manage download queues, edit settings, or create/delete queues.
- Keyboard shortcuts for navigation and actions.
- Footer bar displaying helpful key bindings.

### Persistence
- Save and restore the state of downloads and queues when the application is closed and reopened.
- Store configuration and download information in a json based database.

---

## Getting Started

### Prerequisites
- **Golang** (version 1.18 or later) installed on your system.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/sajjad-mobe/gdm.git
   cd gdm
   go run cmd/main.go

By running the program and navigating to the help tab, you can see how the program works at each step.


## Contributors

This project is maintained by:

- **Sajjad Mohammadbeigi**  [GitHub Profile](https://github.com/sajjad-MoBe)
- **Mohammad Mohammadbeigi** [GitHub Profile](https://github.com/mbmohammad)

Sajjad implemented manager folder parts and Mohammad implemented TUI folder parts.
At the end Sajjad connected these parts to each other.

