package main

import (
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Page int

const (
	LoginPage Page = iota
	AppsPage
)

type LoginModel struct {
	EmailInput    textinput.Model
	PasswordInput textinput.Model
	CurrentIndex  int
	ErrMsg        error
}

type AppsModel struct {
	Apps []utils.App
}

type Model struct {
	CurrentPage Page
	Login       LoginModel
	Apps        AppsModel
}

func NewModel() Model {
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

	return Model{
		CurrentPage: LoginPage,
		Login: LoginModel{
			EmailInput:    emailInput,
			PasswordInput: passwordInput,
			CurrentIndex:  0,
			ErrMsg:        nil,
		},
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.CurrentPage {

	case LoginPage:
		var loginCmd tea.Cmd
		m.Login, loginCmd, m.CurrentPage = UpdateLogin(m.Login, msg)
		return m, loginCmd

	case AppsPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "tab", "enter":
				m.CurrentPage = LoginPage
				return m, cmd
			}
		}
	}

	m.Login.EmailInput, cmd = m.Login.EmailInput.Update(msg)
	m.Login.PasswordInput, _ = m.Login.PasswordInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.CurrentPage {
	case LoginPage:
		return "\nEnter your Email and password:\n" +
			m.Login.EmailInput.View() + "\n" +
			m.Login.PasswordInput.View() + "\n\nPress 'q' to quit.\n"
	case AppsPage:
		return "Welcome to the Apps Page!"
	}
	return ""
}
