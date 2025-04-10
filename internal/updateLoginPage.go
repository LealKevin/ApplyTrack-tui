package main

import (
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) UpdateLogin(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case utils.LoginMsg:
		if msg.Err != nil {
			m.Login.ErrMsg = msg.Err
			return m, nil
		}
		return m, utils.SaveTokenCmd(msg.Token)

	case utils.SaveTokenMsg:
		if msg.Err != nil {
			m.Login.ErrMsg = msg.Err
			return m, nil
		}
		m.CurrentPage = AppsPage
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			m.Login.ErrMsg = nil
			user := utils.UserInputParams{
				Email:    m.Login.EmailInput.Value(),
				Password: m.Login.PasswordInput.Value(),
			}
			return m, utils.LoginCmd(user)

		case tea.KeyTab, tea.KeyShiftTab:
			if m.Login.CurrentIndex == 0 {
				m.Login.EmailInput.Blur()
				m.Login.PasswordInput.Focus()
				m.Login.CurrentIndex = 1
			} else {
				m.Login.EmailInput.Focus()
				m.Login.PasswordInput.Blur()
				m.Login.CurrentIndex = 0
			}
		}

		if m.Login.CurrentIndex == 0 {
			var cmd tea.Cmd
			m.Login.EmailInput, cmd = m.Login.EmailInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			var cmd tea.Cmd
			m.Login.PasswordInput, cmd = m.Login.PasswordInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}
