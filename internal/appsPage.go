package main

import tea "github.com/charmbracelet/bubbletea"

func (m Model) UpdateAppsPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.CurrentPage = LoginPage
			return m, nil
		}
		switch msg.String() {
		case "l":
			m.CurrentPage = LogoutPage
		}
	}
	return m, tea.Batch(cmds...)
}
