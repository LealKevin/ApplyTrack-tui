package main

import (
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Lipgloss style
var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(1, 2).
			Align(lipgloss.Center).
			BorderForeground(lipgloss.Color("23"))

	centerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center)
)

type Page int

const (
	LoginPage Page = iota
	LogoutPage
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
	User         User
	WindowWidth  int
	WindowHeight int

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
	var cmds []tea.Cmd
	cmds = append(cmds, CheckTokenCmd())
	cmds = append(cmds, textinput.Blink)
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case CheckTokenMsg:
		if msg.Err == nil {
			m.CurrentPage = AppsPage
			m.User = msg.User
		}
	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
	}

	switch m.CurrentPage {
	case LoginPage:
		m.Login, cmd, m.CurrentPage = UpdateLogin(m.Login, msg)
		return m, cmd

	case AppsPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyTab:
				m.CurrentPage = LoginPage
				return m, cmd
			}
			switch msg.String() {
			case "l":
				m.CurrentPage = LogoutPage
				return m, cmd
			}
		}

	case LogoutPage:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y":
				m.User.Name = ""
				m.User.Email = ""
				m.Login.EmailInput.Placeholder = "Email"
				m.Login.PasswordInput.Placeholder = "********"

				m.CurrentPage = LoginPage
				cmd = utils.Logout()
				return m, cmd
			case "n":
				m.CurrentPage = AppsPage
				return m, cmd
			}
		}
	}

	m.Login.EmailInput, cmd = m.Login.EmailInput.Update(msg)
	m.Login.PasswordInput, _ = m.Login.PasswordInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	centered := centerStyle.Width(m.WindowWidth).Height(m.WindowHeight)

	var content string
	switch m.CurrentPage {

	case LoginPage:
		content = m.ViewLoginPage()

	case AppsPage:
		content = m.ViewAppsPage()

	case LogoutPage:
		content = m.ViewLogoutPage()
	}

	return centered.Render(borderStyle.Render(content))
}
