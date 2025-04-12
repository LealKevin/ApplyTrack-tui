#  AppliTrack - TUI (Work In Progress)

**AppliTrack** is a lightweight Terminal UI (TUI) to **manage job applications** super fast, with just a few key presses.

Built with [Bubbletea](https://github.com/charmbracelet/bubbletea), [Bubble-table](https://github.com/Evertras/bubble-table), and [Lipgloss](https://github.com/charmbracelet/lipgloss).

---

##  Goal

Quickly create, filter, and manage your job applications **without ever leaving your terminal**.

Few clicks. Maximum efficiency.

---

##  Features

-  Login / Logout
-  List all applications
-  Filter by status (`all`, `sent`, `pending`, `rejected`)
-  Create a new application
- ⌨ Navigate everything with the keyboard
-  Clean and readable styling

---

## ️ Tech

- **Go**
- **Bubbletea** (TUI framework)
- **Bubble-table** (smart tables)
- **Lipgloss** (terminal styling)

---

##  Controls

| Key          | Action                 |
|:-------------|:------------------------|
| `Tab`        | Next input field         |
| `Shift + Tab`| Previous input field     |
| `Enter`      | Submit / Validate        |
| `Esc`        | Cancel or Exit           |
| `/`          | Start typing to filter   |
| `e`          | Edit selected application|
| `n`          | New application          |
| `r`          | Refresh list             |
| `l`          | Logout                   |

---

##  Running locally

```bash
git clone https://github.com/LealKevin/tui-applitrack.git
cd tui-applitrack
go run ./internal/.

