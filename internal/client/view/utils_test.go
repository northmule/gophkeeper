package view

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestRenderCheckbox(t *testing.T) {
	checked := renderCheckbox("Checked", true)
	expectedChecked := lipgloss.NewStyle().Foreground(lipgloss.Color("#57e389")).Render("[->] Checked")
	assert.Equal(t, expectedChecked, checked)

	unchecked := renderCheckbox("Unchecked", false)
	expectedUnchecked := "[ ] Unchecked"
	assert.Equal(t, expectedUnchecked, unchecked)
}

func TestRenderTitle(t *testing.T) {
	title := renderTitle("Test Title")
	expectedTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#e66100")).Bold(true).AlignVertical(lipgloss.Center).BorderBottomForeground(lipgloss.Color("#e66100")).Render("\nTest Title\n")
	assert.Equal(t, expectedTitle, title)

	emptyTitle := renderTitle("")
	expectedEmptyTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#e66100")).Bold(true).AlignVertical(lipgloss.Center).BorderBottomForeground(lipgloss.Color("#e66100")).Render("\n\n")
	assert.Equal(t, expectedEmptyTitle, emptyTitle)
}
