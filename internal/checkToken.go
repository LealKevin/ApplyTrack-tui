package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CheckTokenMsg struct {
	User User
	Err  error
}

func CheckTokenCmd() tea.Cmd {
	return func() tea.Msg {
		tokenBytes, err := os.ReadFile(".token")
		if err != nil {
			return CheckTokenMsg{User: User{}, Err: fmt.Errorf("Unable to read file token: %v", err)}
		}

		token := string(tokenBytes)

		req, err := http.NewRequest("GET", "http://localhost:8080/api/me", nil)
		if err != nil {
			return CheckTokenMsg{User: User{}, Err: fmt.Errorf("Unable to create request: %v", err)}
		}
		req.Header.Set("Cookie", "jwt="+token)

		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			return CheckTokenMsg{User: User{}, Err: fmt.Errorf("Unable to send request: %v", err)}
		}
		defer resp.Body.Close()

		var userResp User
		if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
			return CheckTokenMsg{User: User{}, Err: fmt.Errorf("Unable to decode json: %v", err)}
		}

		return CheckTokenMsg{User: userResp, Err: nil}
	}
}
