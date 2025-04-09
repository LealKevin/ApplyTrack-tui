package main

import "tui-apptrack/internal/misc"

func (m Model) ViewLoginPage() string {
	content := misc.Logo + "\n\nEnter your Email and password:\n\n" +
		m.Login.EmailInput.View() + "\n" +
		m.Login.PasswordInput.View() + "\n\nPress 'esc' to quit.\n"
	return content
}

func (m Model) ViewLogoutPage() string {
	content := m.User.Name +
		"\nAre you sure to logout?:\n\n" +
		"Confirm(y) Cancel(n)"
	return content
}

func (m Model) ViewAppsPage() string {
	content := m.User.Name +
		"\nApps Page:\n\n"

	return content

}
