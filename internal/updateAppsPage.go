package main

import (
	"fmt"
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/pkg/browser"
)

func (a AppsModel) filterRows(status string) []table.Row {
	var rows []table.Row

	for _, app := range a.Apps {
		if status == "all" || app.Status == status {
			rows = append(rows, MakeAppRow(app))
		}
	}
	return rows
}

func (m Model) UpdateAppsPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.Alerts = "All Applications"

	if m.CreateApp.focused {
		return m.UpdateCreateApp(msg)
	}

	switch msg := msg.(type) {
	case utils.FetchAppsMsg:
		if msg.Err != nil {
			m.Apps.Err = msg.Err
			return m, nil
		}
		m.Apps.Apps = msg.Apps
		filtered := m.Apps.filterRows("all")
		m.Apps.table = m.Apps.table.WithRows(filtered)
		return m, nil

	case utils.DeleteMsg:
		if msg.Err != nil {
			m.Alerts = fmt.Sprintf("%v", msg.Err)
			return m, nil
		}

		m.Alerts = "Sucessfull deleted"
		return m, utils.FetchAppsCmd()

	case tea.KeyMsg:
		m.Alerts = ""
		switch msg.String() {
		case "/":
			m.Apps.isFiltering = true
			m.Apps.table, cmd = m.Apps.table.Update(msg)
			return m, cmd
		case "esc", "enter":
			if m.Apps.isFiltering {
				m.Apps.isFiltering = false
				m.Apps.table, cmd = m.Apps.table.Update(msg)
				return m, cmd
			}
		}

		if m.Apps.isFiltering {
			m.Apps.table, cmd = m.Apps.table.Update(msg)
			return m, cmd
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			m.Apps.table = m.Apps.table.PageDown()
			return m, nil
		case tea.KeyShiftTab:
			m.Apps.table = m.Apps.table.PageUp()
			return m, nil
		}

		switch msg.String() {
		case "v":
			row := m.Apps.table.HighlightedRow()
			rawURL, ok := row.Data["url"]
			if !ok {
				fmt.Println("Error: URL not found")
				return m, nil
			}

			urlStr, ok := rawURL.(string)
			if !ok || urlStr == "" {
				m.Alerts = "No valid URL to open."
				return m, nil
			}

			err := browser.OpenURL(urlStr)
			if err != nil {
				m.Alerts = "Failed to open URL."
			} else {
				m.Alerts = "Opening URL in your browser..."
			}
			return m, nil

		case "e":
			row := m.Apps.table.HighlightedRow()
			id, ok := row.Data["id"].(int32)
			if !ok {
				fmt.Printf("Error: casting ID:%v", row)
				return m, nil
			}
			for _, app := range m.Apps.Apps {
				if app.ID == id {
					m.Apps.Temp = app
					break
				}
			}
			return m, nil
		case "d":
			if !m.Delete.ConfirmDelete {
				m.Delete.ConfirmDelete = true
				m.Alerts = "Press 'd' again to confirm deletion"
				return m, nil
			}
			m.Delete.ConfirmDelete = false

			row := m.Apps.table.HighlightedRow()
			rawAppId := row.Data["id"]
			id, ok := rawAppId.(int32)
			if !ok {
				m.Alerts = "Invalid app ID"
				return m, nil
			}
			appIDstr := fmt.Sprintf("%d", id)

			return m, utils.DeleteAppCmd(appIDstr)

		case "n":
			m.CreateApp.focused = true
			m.Apps.table = m.Apps.table.WithBaseStyle(tableBlured)
			m.Apps.table.Focused(false)
			m.Alerts = "Create new application"
			return m, nil
		case "4":
			filtered := m.Apps.filterRows("all")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			m.Alerts = "All Applications"
			return m, nil
		case "1":
			filtered := m.Apps.filterRows("sent")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			m.Alerts = "Filtered by status: sent"
			return m, nil
		case "2":
			filtered := m.Apps.filterRows("pending")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			m.Alerts = "Filtered by status: pending"
			return m, nil
		case "3":
			filtered := m.Apps.filterRows("rejected")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			m.Alerts = "Filtered by status: rejected"
			return m, nil
		case "r":
			return m, utils.FetchAppsCmd()
		case "l":
			m.CurrentPage = LogoutPage
			return m, nil
		default:
			m.Apps.table, cmd = m.Apps.table.Update(msg)
			return m, cmd
		}
	}

	if len(m.Apps.Apps) == 0 && m.Apps.Err == nil {
		return m, utils.FetchAppsCmd()
	}

	return m, tea.Batch()
}
