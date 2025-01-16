package view

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	dotChar = " • "
)

// Общие стили
var (
	// Заголовки
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#e66100")).Bold(true).AlignVertical(lipgloss.Center).BorderBottomForeground(lipgloss.Color("#e66100"))
	bodyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#c0bfbc")).AlignHorizontal(lipgloss.Left)
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#57e389"))
	dotStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle         = lipgloss.NewStyle().MarginLeft(2)
	responseTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#418ce6"))
)

// чекбокс
func renderCheckbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[->] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

func renderTitle(label string) string {
	return titleStyle.Render("\n" + label + "\n")
}
