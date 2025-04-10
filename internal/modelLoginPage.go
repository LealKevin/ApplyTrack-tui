package main

import "github.com/charmbracelet/bubbles/textinput"

func NewLoginModel() LoginModel {
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

	return LoginModel{
		EmailInput:    emailInput,
		PasswordInput: passwordInput,
		CurrentIndex:  0,
		ErrMsg:        nil,
	}
}
