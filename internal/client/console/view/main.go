package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	infoStyle = lipgloss.NewStyle().Bold(true).Border(lipgloss.HiddenBorder()).Foreground(lipgloss.Color("#ff0000"))
)

func main() {
	fmt.Println(infoStyle.Render("GoPhkeeper"))
	m := menu{
		options: []menuItem{
			menuItem{
				text:    "new check-in",
				onPress: func() tea.Msg { return toggleCasingMsg{} },
			},
			menuItem{
				text:    "view check-ins",
				onPress: func() tea.Msg { return toggleCasingMsg{} },
			},
		},
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
