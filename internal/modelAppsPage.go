package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

func NewAppsModel() AppsModel {

	leftAlign := lipgloss.NewStyle().Align(lipgloss.Left)

	columns := []table.Column{
		table.NewColumn("title", "Title", 20).WithStyle(leftAlign),
		table.NewColumn("company", "Company", 20).WithStyle(leftAlign),
		table.NewColumn("status", "Status", 10).WithStyle(leftAlign),
		table.NewColumn("sent", "Sent", 12).WithStyle(leftAlign),
	}

	t := table.New(columns).
		WithTargetWidth(80).
		WithPageSize(10).
		Focused(true)

	t = t.BorderRounded()

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("5"))

	t = t.HeaderStyle(headerStyle)

	return AppsModel{
		table: t,
	}
}
