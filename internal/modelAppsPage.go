package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	tableFocus = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("#dcc394")).
			Foreground(lipgloss.Color("#c6c8c5"))

	tableBlured = lipgloss.NewStyle().
			BorderForeground(lipgloss.Color("240")).
			Foreground(lipgloss.Color("#c6c8c5"))
)

func NewAppsModel() AppsModel {

	leftAlign := lipgloss.NewStyle().Align(lipgloss.Left)

	columns := []table.Column{
		table.NewColumn("title", "Title", 30).WithStyle(leftAlign).WithFiltered(true),
		table.NewColumn("company", "Company", 20).WithStyle(leftAlign).WithFiltered(true),
		table.NewColumn("status", "Status", 20).WithStyle(leftAlign),
		table.NewColumn("sent", "Sent", 15).WithStyle(leftAlign),
	}

	t := table.New(columns).
		WithTargetWidth(100).
		WithMaxTotalWidth(90).
		WithPageSize(16).
		Focused(true).
		Filtered(true)

	t = t.BorderRounded()
	t = t.WithBaseStyle(tableFocus)
	t = t.SortByDesc("sent")
	t = t.WithMinimumHeight(22)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#dcc394")).Align(lipgloss.Center)

	t = t.HeaderStyle(headerStyle)

	return AppsModel{
		table: t,
	}
}
