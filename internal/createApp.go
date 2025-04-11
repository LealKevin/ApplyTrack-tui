package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateAppModel struct {
	inputs       []textinput.Model
	CurrentIndex int
	ErrMsg       error
}

const (
	title = iota
	company
	sentDate
	status
	url
)

func NewCreateAppModel() CreateAppModel {
	inputs := make([]textinput.Model, 5)

	inputs[title] = textinput.New()
	inputs[title].Placeholder = "Title"
	inputs[title].Focus()
	inputs[title].CharLimit = 156
	inputs[title].Width = 20

	inputs[company] = textinput.New()
	inputs[company].Placeholder = "Company"
	inputs[company].CharLimit = 156
	inputs[company].Width = 20

	inputs[status] = textinput.New()
	inputs[status].Placeholder = "Status"
	inputs[status].CharLimit = 20
	inputs[status].Width = 20

	inputs[sentDate] = textinput.New()
	inputs[sentDate].Placeholder = "Sent Date"
	inputs[sentDate].CharLimit = 20
	inputs[sentDate].Width = 20

	inputs[url] = textinput.New()
	inputs[url].Placeholder = "URL"
	inputs[url].CharLimit = 200
	inputs[url].Width = 20

	return CreateAppModel{
		inputs:       []textinput.Model{inputs[title], inputs[company], inputs[status], inputs[sentDate], inputs[url]},
		CurrentIndex: 0,
		ErrMsg:       nil,
	}
}

func (m Model) UpdateCreateApp(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
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
	content := "\n\n" +
		m.CreateApp.inputs[title].View() + "\n" +
		m.CreateApp.inputs[company].View() + "\n" +
		m.CreateApp.inputs[sentDate].View() + "\n" +
		m.CreateApp.inputs[status].View() + "\n" +
		m.CreateApp.inputs[url].View() + "\n" +
		"\n\n" +
		"Press 'tab' to switch fields, 'enter' to submit, 'esc' to quit.\n"
	return content
}
