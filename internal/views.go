package main

import (
	"github.com/charmbracelet/lipgloss"
	t "github.com/evertras/bubble-table/table"
	"tui-apptrack/internal/misc"
	"tui-apptrack/utils"
)

var (
	greyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#dcc394"))

	blurredInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))

	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

var ErrMsg string

func (m Model) ViewLoginPage() string {
	emailStyle := blurredInputStyle
	passwordStyle := blurredInputStyle

	if m.Login.ErrMsg != nil {
		ErrMsg = errorStyle.Render("Something went wrong: " + m.Login.ErrMsg.Error())
	}

	if m.Login.CurrentIndex == 0 {
		emailStyle = focusedInputStyle
	} else {
		passwordStyle = focusedInputStyle
	}

	content := misc.Logo + "\n\n" +
		emailStyle.Render(m.Login.EmailInput.View()) + "\n" +
		passwordStyle.Render(m.Login.PasswordInput.View()) + "\n\n" +
		ErrMsg + "\n\n" +
		"Press 'tab' to switch fields, 'enter' to submit, 'esc' to quit.\n"
	return content
}

func (m Model) ViewLogoutPage() string {
	content := m.User.Name +
		"\nAre you sure to logout?:\n\n" +
		"Confirm(y) Cancel(n)"
	return content
}

func (m Model) ViewAppsPage() string {

	content := ""

	if m.Apps.Err != nil {
		content += "Error: " + m.Apps.Err.Error() + "\n"
	}

	content =
		m.Apps.table.View() +
			"\n" + greyStyle.Render("Filter by status: all(1) sent(2) pending(3) rejected(4)\n")
	return content
}

func MakeAppRow(app utils.App) t.Row {
	var statusColor string

	switch app.Status {
	case "sent":
		statusColor = "#8a943e"
	case "pending":
		statusColor = "#de935f"
	case "rejected":
		statusColor = "#a54241"
	default:
		statusColor = "7"
	}

	return t.NewRow(t.RowData{
		"id":      app.ID,
		"title":   app.TitleApplication,
		"company": app.Company,
		"status":  t.NewStyledCell(app.Status, lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor))),
		"sent":    app.SentDate,
	})
}
