package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	ID               int32
	TitleApplication string
	Company          string
	SentDate         int
	Status           string
	Notes            string
	UrlApplication   string
	UserID           int32
	CreatedAt        int
	UpdatedAt        int
}

type Apps []App

func FetchApps(token string) (Apps, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/applications", nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create a request, error: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %s", resp.Status)
	}

	var apps Apps
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("Unable to fetch apps, error: %v", err)
	}

	return apps, nil
}

type userLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInputParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginMsg struct {
	Token string
	Err   error
}

func LoginCmd(user UserInputParams) tea.Cmd {
	return func() tea.Msg {
		var body UserInputParams
		body.Email = user.Email
		body.Password = user.Password

		jsonBody, err := json.Marshal(body)
		if err != nil {
			return LoginMsg{Token: "", Err: fmt.Errorf("Unable to encode user to json: %v", err)}
		}
		req, err := http.NewRequest("POST", "http://localhost:8080/api/login", bytes.NewBuffer(jsonBody))
		if err != nil {
			return LoginMsg{Token: "", Err: fmt.Errorf("Unable to create request error: %v", err)}
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			return LoginMsg{Token: "", Err: fmt.Errorf("Unable to request error: %v", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return LoginMsg{Token: "", Err: fmt.Errorf("server returned status %s", resp.Status)}
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "jwt" {
				return LoginMsg{Token: cookie.Value, Err: nil}
			}
		}
		return LoginMsg{Token: "", Err: fmt.Errorf("token cookie 'jwt' not found")}
	}
}

type SaveTokenMsg struct {
	Err error
}

func SaveTokenCmd(token string) tea.Cmd {
	return func() tea.Msg {
		if token == "" {
			return SaveTokenMsg{Err: fmt.Errorf("Unable to find token")}
		}

		err := os.WriteFile(".token", []byte(token), 0600)
		if err != nil {
			return SaveTokenMsg{Err: fmt.Errorf("Unable to save token on file: %v", err)}
		}

		return SaveTokenMsg{Err: nil}
	}
}

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
