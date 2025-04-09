package main

import (
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func UpdateLogin(m LoginModel, msg tea.Msg) (LoginModel, tea.Cmd, Page) {
	var cmds []tea.Cmd
	var page Page = LoginPage

	switch msg := msg.(type) {
	case utils.LoginMsg:
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil, page
		}
		return m, utils.SaveTokenCmd(msg.Token), page

	case utils.SaveTokenMsg:
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil, page
		}
		page = AppsPage
		return m, nil, page

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit, page

		case tea.KeyEnter:
			user := utils.UserInputParams{
				Email:    m.EmailInput.Value(),
				Password: m.PasswordInput.Value(),
			}
			return m, utils.LoginCmd(user), page

		case tea.KeyTab, tea.KeyShiftTab:
			if m.CurrentIndex == 0 {
				m.EmailInput.Blur()
				m.PasswordInput.Focus()
				m.CurrentIndex = 1
			} else {
				m.EmailInput.Focus()
				m.PasswordInput.Blur()
				m.CurrentIndex = 0
			}
		}

		if m.CurrentIndex == 0 {
			var cmd tea.Cmd
			m.EmailInput, cmd = m.EmailInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.PasswordInput, cmd = m.PasswordInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...), page
}
