package main

import (
	"tui-apptrack/internal/misc"
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LipGloss styles
var (
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2).
			Align(lipgloss.Center).
			BorderForeground(lipgloss.Color("63"))

	centerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(50).
			Height(15).
			Margin(1, 2)
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
	User         utils.User
	windowWidth  int
	windowHeight int

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
	cmds = append(cmds, utils.CheckTokenCmd())
	cmds = append(cmds, textinput.Blink)
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case utils.CheckTokenMsg:
		if msg.Err == nil {
			m.CurrentPage = AppsPage
			m.User = msg.User
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	}

	switch m.CurrentPage {
	case LoginPage:
		m.Login, cmd, m.CurrentPage = UpdateLogin(m.Login, msg)
		return m, cmd

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
	centerStyle := lipgloss.NewStyle().
		Width(m.windowWidth).
		Height(m.windowHeight).
		Align(lipgloss.Center, lipgloss.Center)

	switch m.CurrentPage {
	//Login page's style
	case LoginPage:
		content := misc.Logo + "\n\nEnter your Email and password:\n\n" +
			m.Login.EmailInput.View() + "\n" +
			m.Login.PasswordInput.View() + "\n\nPress 'esc' to quit.\n"

		boxed := borderStyle.Render(content)
		return centerStyle.Render(boxed)

		//Apps page's style
	case AppsPage:
		content := m.User.Name +
			"\nApps Page:\n\n"
		boxed := borderStyle.Render(content)
		return centerStyle.Render(boxed)
	}
	return ""
}
