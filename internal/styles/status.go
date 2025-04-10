package styles

import "github.com/charmbracelet/lipgloss"

func StyleStatus(status string) string {
	switch status {
	case "rejected":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Render(status)
	case "pending":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Render(status)
	case "sent":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Render(status)
	default:
		return status
	}
}
