package main

import (
	"tui-apptrack/internal/misc"

	"github.com/charmbracelet/lipgloss"
)

var (
	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205"))

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
	content := m.User.Name +
		"\nApps Page:\n\n"

	return content

}
