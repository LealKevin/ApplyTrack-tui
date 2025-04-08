package main

import (
	"fmt"
	"log"
	"tui-apptrack/handlers"
	"tui-apptrack/models"
	m "tui-apptrack/models"
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func initModel() m.Model {
	emailInput := textinput.New()
	emailInput.Placeholder = "Email"
	emailInput.Focus()
	emailInput.CharLimit = 100
	emailInput.Width = 30

	passwordInput := textinput.New()
	passwordInput.Placeholder = "********"
	passwordInput.CharLimit = 100
	passwordInput.Width = 30
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '*'

	return m.Model{
		CurrentPage:   m.LoginPage,
		EmailInput:    emailInput,
		PasswordInput: passwordInput,
		ErrMsg:        nil,
	}
}

func Update(m m.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.CurrentPage {
	case models.LoginPage:
		return handlers.UpdateLogin(m, msg)
	}
	switch msg := msg.(type) {
	case utils.LoginMsg:
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil
		}
		m.CurrentPage = m.AppsPage
		return m, utils.SaveTokenCmd(msg.Token)

	case utils.SaveTokenMsg:
		fmt.Printf("here")
		if msg.Err != nil {
			m.ErrMsg = msg.Err
			return m, nil
		}

		fmt.Printf("Token saved")
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

func View(m m.Model) string {
	switch m.CurrentPage {
	case m.LoginPage:
		return fmt.Sprintf("\nEnter your Email and password:\n%s\n%s\n",
			m.EmailInput.View(),
			m.PasswordInput.View(),
		) + "\n"
	case m.AppsPage:
		return fmt.Sprintf("Iam here applications page")
	}
	return fmt.Sprint("")
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

}
