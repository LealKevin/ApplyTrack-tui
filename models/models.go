package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Page int

const (
	LoginPage Page = iota
	AppsPage
)

type Model struct {
	CurrentPage   Page
	EmailInput    textinput.Model
	PasswordInput textinput.Model
	CurrentIndex  int

	ErrMsg error
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}
