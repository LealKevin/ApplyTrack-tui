package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func CheckTokenCmd() tea.Cmd {
	return func() tea.Msg {

		os.ReadFile(".token")

		return true
	}

}
