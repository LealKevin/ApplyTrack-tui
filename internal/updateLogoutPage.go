package main

import (
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) UpdateLogoutPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			m.User.Name = ""
			m.User.Email = ""
			m.Login.EmailInput.Placeholder = "Email"
			m.Login.PasswordInput.Placeholder = "********"

			m.CurrentPage = LoginPage
			cmd := utils.Logout()
			return m, cmd
		case "n":
			m.CurrentPage = AppsPage
		}
	}
	return m, nil
}
