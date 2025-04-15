package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	ID               int32
	TitleApplication string
	Company          string
	SentDate         string
	Status           string
	Notes            string
	UrlApplication   string
	UserID           int32
	CreatedAt        string
	UpdatedAt        string
}

type FetchAppsMsg struct {
	Apps []App
	Err  error
}

const url string = "http://localhost:8080/api/"

func FetchAppsCmd() tea.Cmd {
	return func() tea.Msg {

		tokenBytes, err := os.ReadFile(".token")
		if err != nil {
			return FetchAppsMsg{Apps: nil, Err: fmt.Errorf("unable to read token: %v", err)}
		}
		token := string(tokenBytes)

		req, err := http.NewRequest("GET", url+"applications?status=all", nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return FetchAppsMsg{Apps: nil, Err: fmt.Errorf("Unable to create a request, error: %v", err)}
		}

		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
		})

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error executing request: %v\n", err)
			return FetchAppsMsg{Apps: nil, Err: fmt.Errorf("Unable to execute request, error: %v", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: server returned non-OK status: %d %s\n",
				resp.StatusCode, http.StatusText(resp.StatusCode))
			return FetchAppsMsg{
				Apps: nil,
				Err:  fmt.Errorf("server returned status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)),
			}
		}
		bodyBytes, err := io.ReadAll(resp.Body)

		var apps []App
		if err := json.Unmarshal(bodyBytes, &apps); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			return FetchAppsMsg{Apps: nil, Err: fmt.Errorf("Unable to decode response, error: %v", err)}
		}
		return FetchAppsMsg{Apps: apps, Err: nil}
	}
}

type CreateAppMsg struct {
	App     App
	Err     error
	Created bool
}

type CreateAppRequest struct {
	TitleApplication string `json:"title"`
	Company          string `json:"company"`
	SentDate         string `json:"sent_date"`
	Status           string `json:"status"`
	UrlApplication   string `json:"url_application"`
}

func CreateApp(app CreateAppRequest) tea.Cmd {
	fmt.Printf("Sent date: %v, Status: %v ", app.SentDate, app.Status)
	return func() tea.Msg {

		tokenBytes, err := os.ReadFile(".token")
		if err != nil {
			return CreateAppMsg{Err: fmt.Errorf("To read token file: %v", err)}
		}
		token := string(tokenBytes)
		jsonApp, err := json.Marshal(app)
		if err != nil {
			return CreateAppMsg{Err: fmt.Errorf("Unable to encode app to json: %v", err)}
		}

		req, err := http.NewRequest("POST", url+"application", bytes.NewBuffer(jsonApp))
		if err != nil {
			return CreateAppMsg{Err: fmt.Errorf("Unable to create request: %v", err)}
		}

		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
		})

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return CreateAppMsg{Err: fmt.Errorf("Unable to send request: %v", err)}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return CreateAppMsg{Err: fmt.Errorf("Server returned status %s", resp.Status)}
		}

		var respApp App
		if err := json.NewDecoder(resp.Body).Decode(&respApp); err != nil {
			return CreateAppMsg{Err: fmt.Errorf("Unable to decode json: %v", err)}
		}
		return CreateAppMsg{App: respApp, Err: nil, Created: true}
	}
}

type userLoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInputParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOutputParams struct {
	Name string `json:"name"`
}

type LoginMsg struct {
	Name  string
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
			return LoginMsg{Token: "", Err: fmt.Errorf("Server returned status %s", resp.Status)}
		}
		var user UserOutputParams
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return LoginMsg{Name: "", Token: "", Err: fmt.Errorf("Unable to decode json: %v", err)}
		}

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "jwt" {
				return LoginMsg{Name: user.Name, Token: cookie.Value, Err: nil}
			}
		}

		return LoginMsg{Name: user.Name, Token: "", Err: fmt.Errorf("token cookie 'jwt' not found")}
	}
}

type LogoutMsg struct {
	Err error
}

func Logout() tea.Cmd {
	return func() tea.Msg {

		err := os.Remove(".token")
		if err != nil {
			return LogoutMsg{Err: err}
		}

		return nil
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

type DeleteMsg struct {
	Err error
}

func DeleteAppCmd(appId string) tea.Cmd {
	return func() tea.Msg {
		tokenBytes, err := os.ReadFile(".token")
		if err != nil {
			return CreateAppMsg{Err: fmt.Errorf("To read token file: %v", err)}
		}
		token := string(tokenBytes)

		req, err := http.NewRequest("DELETE", url+"applications/"+appId, nil)
		if err != nil {
			return DeleteMsg{Err: fmt.Errorf("Unable to create request error: %v", err)}
		}

		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{
			Name:  "jwt",
			Value: token,
		})

		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			return DeleteMsg{Err: fmt.Errorf("Unable to request error: %v", err)}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return DeleteMsg{Err: fmt.Errorf("Status error: %v", resp.Status)}
		}

		return DeleteMsg{Err: nil}
	}
}
