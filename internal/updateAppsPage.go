package main

import (
	"fmt"
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
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
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case utils.FetchAppsMsg:
		if msg.Err != nil {
			m.Apps.Err = msg.Err
			return m, nil
		}
		m.Apps.Apps = msg.Apps
		filtered := m.Apps.filterRows("all")
		m.Apps.table = m.Apps.table.WithRows(filtered)

		m.Apps.Apps = msg.Apps
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.CurrentPage = LoginPage
			return m, nil
		}
		switch msg.String() {
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

		case "0":
			filtered := m.Apps.filterRows("all")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			return m, nil
		case "1":
			filtered := m.Apps.filterRows("sent")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			return m, nil
		case "2":
			filtered := m.Apps.filterRows("pending")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			return m, nil
		case "3":
			filtered := m.Apps.filterRows("rejected")
			m.Apps.table = m.Apps.table.WithRows(filtered)
			return m, nil
		case "r":
			return m, utils.FetchAppsCmd()
		case "esc":
			return m, tea.Quit
		case "l":
			m.CurrentPage = LogoutPage
		default:
			m.Apps.table, cmd = m.Apps.table.Update(msg)
			return m, cmd
		}
	}

	if len(m.Apps.Apps) == 0 && m.Apps.Err == nil {
		return m, utils.FetchAppsCmd()
	}

	return m, tea.Batch(cmds...)
}
