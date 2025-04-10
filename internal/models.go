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
	return tea.Batch(
		CheckTokenCmd(),
		textinput.Blink,
	)
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
		return m.UpdateLogin(msg)
	case AppsPage:
		return m.UpdateAppsPage(msg)
	case LogoutPage:
		return m.UpdateLogoutPage(msg)
	}
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
