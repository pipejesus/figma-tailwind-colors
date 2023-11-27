package helpers

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func Pretty(prefix string, str string) {
	blueStyle := lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#567DF4")).
		PaddingLeft(1).
		PaddingRight(1)
	plainStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.NoColor{}).
		Background(lipgloss.NoColor{}).
		PaddingLeft(1).
		PaddingRight(1)
	fmt.Println(blueStyle.Render(prefix) + plainStyle.Render(str))
}
