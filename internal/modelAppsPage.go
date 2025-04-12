package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

func NewAppsModel() AppsModel {

	leftAlign := lipgloss.NewStyle().Align(lipgloss.Left)

	baseStyle := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#dcc394")).
		Foreground(lipgloss.Color("#c6c8c5"))

	columns := []table.Column{
		table.NewColumn("title", "Title", 20).WithStyle(leftAlign).WithFiltered(true),
		table.NewColumn("company", "Company", 20).WithStyle(leftAlign).WithFiltered(true),
		table.NewColumn("status", "Status", 10).WithStyle(leftAlign),
		table.NewColumn("sent", "Sent", 12).WithStyle(leftAlign),
	}

	t := table.New(columns).
		WithTargetWidth(80).
		WithPageSize(10).
		Focused(true).
		Filtered(true)

	t = t.BorderRounded()
	t = t.WithBaseStyle(baseStyle)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#b293bb"))

	t = t.HeaderStyle(headerStyle)

	return AppsModel{
		table: t,
	}
}
