package main

import (
	"fmt"
	"tui-apptrack/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateAppModel struct {
	inputs       []textinput.Model
	CurrentIndex int
	ErrMsg       error
	isConfirm    bool
}

const (
	title = iota
	company
	year
	month
	day
	status
	url
)

func NewCreateAppModel() CreateAppModel {
	inputs := make([]textinput.Model, 7)

	inputs[title] = textinput.New()
	inputs[title].Placeholder = "Title"
	inputs[title].Focus()
	inputs[title].CharLimit = 156
	inputs[title].Width = 20
	inputs[title].Prompt = ""

	inputs[company] = textinput.New()
	inputs[company].Placeholder = "Company"
	inputs[company].CharLimit = 156
	inputs[company].Width = 20
	inputs[company].Prompt = ""

	inputs[year] = textinput.New()
	inputs[year].Placeholder = "YYYY"
	inputs[year].CharLimit = 4
	inputs[year].Width = 4
	inputs[year].Prompt = ""

	inputs[month] = textinput.New()
	inputs[month].Placeholder = "MM"
	inputs[month].CharLimit = 2
	inputs[month].Width = 2
	inputs[month].Prompt = ""

	inputs[day] = textinput.New()
	inputs[day].Placeholder = "DD"
	inputs[day].CharLimit = 2
	inputs[day].Width = 2
	inputs[day].Prompt = ""

	inputs[status] = textinput.New()
	inputs[status].Placeholder = "Status"
	inputs[status].SetSuggestions([]string{"sent", "pending", "rejected"})
	inputs[status].ShowSuggestions = true
	inputs[status].CharLimit = 20
	inputs[status].Width = 20
	inputs[status].Prompt = ""

	inputs[url] = textinput.New()
	inputs[url].Placeholder = "URL"
	inputs[url].CharLimit = 200
	inputs[url].Width = 20
	inputs[url].Prompt = ""

	return CreateAppModel{
		inputs:       []textinput.Model{inputs[title], inputs[company], inputs[year], inputs[month], inputs[day], inputs[status], inputs[url]},
		CurrentIndex: 0,
		ErrMsg:       nil,
		isConfirm:    false,
	}
}

func (m Model) UpdateCreateApp(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case utils.CreateAppMsg:
		if msg.Err != nil {
			m.CreateApp.ErrMsg = msg.Err
			m.CreateApp.isConfirm = false
			return m, nil
		}
		if msg.Created {
			m.CurrentPage = AppsPage
			return m, utils.FetchAppsCmd()
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.CreateApp.isConfirm {
				m.CreateApp.isConfirm = false
				return m, nil
			}
			m.CurrentPage = AppsPage
			return m, nil

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
				return m, nil
			}

			if m.CreateApp.isConfirm {
				m.CreateApp.isConfirm = false
				app := m.parseApp()
				cmd := utils.CreateApp(app)

				return m, cmd
			}
		}

		for i := range m.CreateApp.inputs {
			m.CreateApp.inputs[i].Blur()
		}
		m.CreateApp.inputs[m.CreateApp.CurrentIndex].Focus()

		i := m.CreateApp.CurrentIndex
		var cmd tea.Cmd
		m.CreateApp.inputs[i], cmd = m.CreateApp.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) ViewCreateAppPage() string {
	confirm := ""
	if m.CreateApp.isConfirm {
		confirm = "Confirm to create app"
	}

	err := ""
	if m.CreateApp.ErrMsg != nil {
		err = errorStyle.Render(m.CreateApp.ErrMsg.Error())

	}
	content :=
		"Create new application" +
			"\n\n" +
			m.CreateApp.inputs[title].View() + "\n" +
			m.CreateApp.inputs[company].View() + "\n" +
			m.CreateApp.inputs[year].View() + "/ " + m.CreateApp.inputs[month].View() + "/ " + m.CreateApp.inputs[day].View() + "\n" +
			m.CreateApp.inputs[status].View() + "\n" +
			m.CreateApp.inputs[url].View() + "\n" +
			"\n\n\n\n" +
			GreyStyle.Render("↓(Tab) ↑(Shift + Tab) ") +
			err +
			GreyStyle.Render("Submit(enter) Cancel(esc).\n") + confirm
	return content
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
	return fmt.Sprintf(year + month + day)

}
