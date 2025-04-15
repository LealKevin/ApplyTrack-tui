package main

import (
	"fmt"
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	borderStyle = lipgloss.NewStyle().
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
	CreateAppPage
)

type LoginModel struct {
	EmailInput    textinput.Model
	PasswordInput textinput.Model
	CurrentIndex  int
	ErrMsg        error
}

type AppsModel struct {
	Apps        []utils.App
	table       table.Model
	Err         error
	isFiltering bool
	Temp        utils.App
}

type EditModel struct {
	err       error
	inEdit    bool
	editingID int32
}

func NewEditModel() EditModel {
	return EditModel{
		err:    nil,
		inEdit: false,
	}
}

type Model struct {
	User         User
	WindowWidth  int
	WindowHeight int

	Alerts      string
	CurrentPage Page
	Login       LoginModel
	Apps        AppsModel
	CreateApp   CreateAppModel
	Delete      DeleteModel
	Edit        EditModel
}

type DeleteModel struct {
	ConfirmDelete bool
}

func NewDeleteModel() DeleteModel {
	return DeleteModel{
		ConfirmDelete: false,
	}
}

func NewModel() Model {

	return Model{
		Alerts:      "",
		CurrentPage: LoginPage,
		Login:       NewLoginModel(),
		Apps:        NewAppsModel(),
		CreateApp:   NewCreateAppModel(),
		Delete:      NewDeleteModel(),
		Edit:        NewEditModel(),
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
			m.User.Name = msg.User.Name
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
	case CreateAppPage:
		return m.UpdateCreateApp(msg)
	}
	return m, cmd
}

func (m Model) viewLeftArea() string {
	return m.ViewStats()
}

func (m Model) viewRightArea() string {
	return m.ViewAppsPage()
}

func (m Model) viewBottomArea() string {
	return m.ViewCreateAppPage()
}

func (m Model) ViewStats() string {
	statsStyle := lipgloss.NewStyle().
		Width(10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#117173")).
		Padding(0)

	return statsStyle.Render(fmt.Sprintf(
		"Apps: %d",
		len(m.Apps.Apps),
	))
}

func (m Model) viewTopArea() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(88).
		AlignHorizontal(lipgloss.Center).
		BorderForeground(lipgloss.Color("#dcc394")).
		Padding(0)
	return border.Render(m.viewAlerts())
}

func (m Model) viewAlerts() string {
	return m.Alerts
}

var leftSize int = 25

func (m Model) viewNotesArea() string {
	style := lipgloss.NewStyle().
		Width(leftSize).
		Height(20).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#117173")).
		Padding(1)

	return style.Render(" Notes\n\n(coming soon...)")
}

func (m Model) spaceHolder() string {
	style := lipgloss.NewStyle().
		Width(leftSize + 2).
		Height(0)

	return style.Render(" ")
}

func (m Model) View() string {
	centered := centerStyle.Width(m.WindowWidth).Height(m.WindowHeight)

	var content string
	switch m.CurrentPage {
	case LoginPage:
		content = m.ViewLoginPage()
	case AppsPage:
		content =
			lipgloss.JoinVertical(

				lipgloss.Right,
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.viewTopArea(),
					m.spaceHolder(),
				),

				lipgloss.JoinHorizontal(lipgloss.Left,
					m.viewLeftArea(),
					m.viewRightArea(),
					m.viewNotesArea(),
				),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.viewBottomArea(),
					m.spaceHolder(),
				),
			)

	case LogoutPage:
		content = m.ViewLogoutPage()
	}

	return centered.Render(borderStyle.Render(content))
}
