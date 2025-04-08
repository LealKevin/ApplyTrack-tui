package handlers

import (
	"fmt"
	"tui-apptrack/models"
	"tui-apptrack/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func UpdateLogin(m models.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case utils.LoginMsg:
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil
		}
		return m, utils.SaveTokenCmd(msg.Token)

	case utils.SaveTokenMsg:
		fmt.Println("Token saved")
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			user := utils.UserInputParams{
				Email:    m.EmailInput.Value(),
				Password: m.PasswordInput.Value(),
			}
			return m, utils.LoginCmd(user)

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

	return m, tea.Batch(cmds...)
}
