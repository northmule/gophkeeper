package view

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// Общие стили
var (
	// Заголовки
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#e66100")).Bold(true).AlignVertical(lipgloss.Center).BorderBottomForeground(lipgloss.Color("#e66100"))
	bodyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#c0bfbc")).AlignHorizontal(lipgloss.Left)
	keywordStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	checkboxStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#57e389"))
	dotStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle         = lipgloss.NewStyle().MarginLeft(2)
	responseTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#418ce6"))

	// Gradient colors we'll use for the progress bar
	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
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

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
