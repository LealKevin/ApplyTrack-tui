package main

import (
	"fmt"
	"strconv"
	"time"
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CreateAppModel struct {
	inputs       []textinput.Model
	CurrentIndex int
	ErrMsg       error
	isConfirm    bool
	focused      bool
	handled      bool
}

const (
	title = iota
	company
	status
	year
	month
	day
	url
)

var (
	statusColor = lipgloss.NewStyle()
)

func today(timeType string) string {

	year, month, day := time.Now().Date()

	var result string

	switch timeType {
	case "year":
		result = strconv.Itoa(year)
	case "month":
		result = fmt.Sprintf("%02d", month)
	case "day":
		result = fmt.Sprintf("%02d", day)
	}

	return result

}

func NewCreateAppModel() CreateAppModel {
	inputs := make([]textinput.Model, 7)

	inputs[title] = textinput.New()
	inputs[title].Placeholder = "Title"
	inputs[title].CharLimit = 156
	inputs[title].Width = 30
	inputs[title].Prompt = ""

	inputs[company] = textinput.New()
	inputs[company].Placeholder = "Company"
	inputs[company].CharLimit = 156
	inputs[company].Width = 20
	inputs[company].Prompt = ""

	inputs[year] = textinput.New()
	inputs[year].Placeholder = today("year")
	inputs[year].SetValue(today("year"))
	inputs[year].CharLimit = 4
	inputs[year].Width = 4
	inputs[year].Prompt = ""

	inputs[month] = textinput.New()
	inputs[month].Placeholder = today("month")
	inputs[month].SetValue(today("month"))
	inputs[month].CharLimit = 2
	inputs[month].Width = 2
	inputs[month].Prompt = ""

	inputs[day] = textinput.New()
	inputs[day].Placeholder = today("day")
	inputs[day].SetValue(today("day"))
	inputs[day].CharLimit = 2
	inputs[day].Width = 2
	inputs[day].Prompt = ""

	inputs[status] = textinput.New()
	inputs[status].Placeholder = "Status"
	inputs[status].CharLimit = 20
	inputs[status].Width = 20
	inputs[status].Prompt = ""

	inputs[url] = textinput.New()
	inputs[url].Placeholder = "URL"
	inputs[url].CharLimit = 200
	inputs[url].Width = 20
	inputs[url].Prompt = ""

	return CreateAppModel{
		inputs:       []textinput.Model{inputs[title], inputs[company], inputs[status], inputs[year], inputs[month], inputs[day], inputs[url]},
		CurrentIndex: 0,
		ErrMsg:       nil,
		isConfirm:    false,
		focused:      false,
	}
}

func (m Model) UpdateCreateApp(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.CreateApp.CurrentIndex == len(m.CreateApp.inputs)-1 {
				m.CreateApp.CurrentIndex = 0
			} else {
				m.CreateApp.CurrentIndex++
			}
		case tea.KeyShiftTab:
			if m.CreateApp.CurrentIndex == 0 {
				m.CreateApp.CurrentIndex = len(m.CreateApp.inputs) - 1
			} else {
				m.CreateApp.CurrentIndex--
			}
		case tea.KeyEnter:
			if !m.CreateApp.isConfirm {
				m.CreateApp.isConfirm = true
				m.Alerts = "Press 'enter' to confirm " + ternary(m.Edit.inEdit, "edition", "creation")
				return m, nil
			}

			m.CreateApp.isConfirm = false
			m.CreateApp.focused = false
			m.Apps.table = m.Apps.table.WithBaseStyle(tableFocus).Focused(true)

			for i := range m.CreateApp.inputs {
				m.CreateApp.inputs[i].Blur()
			}

			if m.Edit.inEdit {
				m.Edit.inEdit = false
				app := m.parseApp()
				cmd := utils.UpdateAppCmd(m.Edit.editingID, app)
				m.ResetCreateAppInputs()
				m.Alerts = "Updating application..."
				return m, cmd
			} else {
				m.Alerts = "Creating application..."
				app := m.parseApp()
				cmd := utils.CreateApp(app)
				m.ResetCreateAppInputs()
				return m, cmd
			}

		case tea.KeyEsc, tea.KeyCtrlC:
			m.CreateApp.isConfirm = false
			m.CreateApp.focused = false
			m.Edit.inEdit = false
			m.Apps.table = m.Apps.table.WithBaseStyle(tableFocus).Focused(true)
			m.ResetCreateAppInputs()
			for i := range m.CreateApp.inputs {
				m.CreateApp.inputs[i].Blur()
			}
			m.Alerts = "Cancelled"
			return m, nil
		}

		switch msg.String() {
		case "1":
			if m.CreateApp.inputs[status].Focused() {
				m.CreateApp.inputs[status].SetValue("sent")
				m.CreateApp.inputs[status].CursorEnd()
				statusColor = statusColor.Foreground(lipgloss.Color("#8a943e"))
				m.CreateApp.handled = true
			}
		case "2":
			if m.CreateApp.inputs[status].Focused() {
				m.CreateApp.inputs[status].SetValue("pending")
				m.CreateApp.inputs[status].CursorEnd()
				statusColor = statusColor.Foreground(lipgloss.Color("#de935f"))
				m.CreateApp.handled = true
			}
		case "3":
			if m.CreateApp.inputs[status].Focused() {
				m.CreateApp.inputs[status].SetValue("rejected")
				m.CreateApp.inputs[status].CursorEnd()
				statusColor = statusColor.Foreground(lipgloss.Color("#a54241"))
				m.CreateApp.handled = true
			}
		}

		for i := range m.CreateApp.inputs {
			m.CreateApp.inputs[i].Blur()
		}
		m.CreateApp.inputs[m.CreateApp.CurrentIndex].Focus()

		if !m.CreateApp.handled {
			i := m.CreateApp.CurrentIndex
			var cmd tea.Cmd
			m.CreateApp.inputs[i], cmd = m.CreateApp.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
		m.CreateApp.handled = false

	case utils.CreateAppMsg:
		if msg.Err != nil {
			m.CreateApp.ErrMsg = msg.Err
			return m, nil
		}
		if msg.Created {
			m.Alerts = "Application created with success"
			m.ResetCreateAppInputs()
			return m, utils.FetchAppsCmd()
		}

	case utils.UpdateAppMsg:
		if msg.Err != nil {
			m.Alerts = fmt.Sprintf("Error while updating: %v", msg.Err)
			return m, nil
		}
		if msg.IsUpdated {
			m.Alerts = "Application updated with success"
			m.ResetCreateAppInputs()
			return m, utils.FetchAppsCmd()
		}
	}

	return m, tea.Batch(cmds...)
}
func ternary(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}

func (m Model) ViewCreateAppPage() string {

	notFocused := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(88).
		BorderForeground(lipgloss.Color("240")).
		Padding(0)

	focused := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(88).
		BorderForeground(lipgloss.Color("#dcc394")).
		Padding(0)

	alignCenter := lipgloss.NewStyle()
	alignCenter.AlignHorizontal(lipgloss.Center)
	confirm := ""

	err := ""
	if m.CreateApp.ErrMsg != nil {
		err = errorStyle.Render(m.CreateApp.ErrMsg.Error())

	}
	content :=
		m.CreateApp.inputs[title].View() +
			m.CreateApp.inputs[company].View() +
			statusColor.Render(m.CreateApp.inputs[status].View()) +
			m.CreateApp.inputs[year].View() + "/" + m.CreateApp.inputs[month].View() + "/" + m.CreateApp.inputs[day].View() + "\n" +
			m.CreateApp.inputs[url].View() + "\n" +
			err

	if confirm != "" {
		content += "\n"
	}

	if m.CreateApp.focused {
		return focused.Render(content)
	} else {
		return notFocused.Render(content)
	}
}

func (m Model) parseApp() utils.CreateAppRequest {
	date := joinDate(m.CreateApp.inputs[year].Value(), m.CreateApp.inputs[month].Value(), m.CreateApp.inputs[day].Value())

	return utils.CreateAppRequest{
		TitleApplication: m.CreateApp.inputs[title].Value(),
		Company:          m.CreateApp.inputs[company].Value(),
		SentDate:         date,
		Status:           m.CreateApp.inputs[status].Value(),
		UrlApplication:   m.CreateApp.inputs[url].Value(),
	}
}

func joinDate(year, month, day string) string {
	return fmt.Sprintf(year + "-" + month + "-" + day)
}

func (m *Model) ResetCreateAppInputs() {
	todayYear := today("year")
	todayMonth := fmt.Sprintf("%02d", time.Now().Month())
	todayDay := fmt.Sprintf("%02d", time.Now().Day())

	m.CreateApp.inputs[title].SetValue("")
	m.CreateApp.inputs[company].SetValue("")
	m.CreateApp.inputs[status].SetValue("")
	m.CreateApp.inputs[year].SetValue(todayYear)
	m.CreateApp.inputs[month].SetValue(todayMonth)
	m.CreateApp.inputs[day].SetValue(todayDay)
	m.CreateApp.inputs[url].SetValue("")

	m.CreateApp.CurrentIndex = title
	m.CreateApp.inputs[title].Focus()
}

func (m *Model) SetEditInputs() {
	row := m.Apps.table.HighlightedRow()
	id, ok := row.Data["id"].(int32)
	if !ok {
		m.Alerts = "Unable to get application ID"
		return
	}

	var app utils.App
	found := false
	for _, a := range m.Apps.Apps {
		if a.ID == id {
			app = a
			found = true
			break
		}
	}

	if !found {
		m.Alerts = "Application not found"
		return
	}

	m.Edit.editingID = app.ID

	m.CreateApp.inputs[title].SetValue(app.TitleApplication)
	m.CreateApp.inputs[company].SetValue(app.Company)
	m.CreateApp.inputs[status].SetValue(app.Status)

	date, err := time.Parse("2006-01-02", app.SentDate)
	if err == nil {
		m.CreateApp.inputs[year].SetValue(fmt.Sprintf("%04d", date.Year()))
		m.CreateApp.inputs[month].SetValue(fmt.Sprintf("%02d", int(date.Month())))
		m.CreateApp.inputs[day].SetValue(fmt.Sprintf("%02d", date.Day()))
	} else {
		m.Alerts = "Invalid date format"
	}
	m.CreateApp.inputs[url].SetValue(app.UrlApplication)

	m.CreateApp.CurrentIndex = title
	m.CreateApp.inputs[title].Focus()

	m.Alerts = "Editing: " + app.TitleApplication
}
