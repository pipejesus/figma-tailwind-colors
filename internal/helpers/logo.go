package helpers

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

const tamagoBlack string = "#212529"
const tamagoYellow string = "#faba0f"

func CreateStyles() (yellowStyle lipgloss.Style, blueStyle lipgloss.Style, darkStyle lipgloss.Style) {
	lipgloss.NewStyle()
	yellowStyle = lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color(tamagoBlack)).
		Background(lipgloss.Color(tamagoYellow)).
		PaddingLeft(1).
		PaddingRight(1)
	blueStyle = lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#567DF4")).
		PaddingLeft(1).
		PaddingRight(1)
	darkStyle = lipgloss.NewStyle().
		Bold(false).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color(tamagoBlack)).
		PaddingLeft(1).
		PaddingRight(1)
	return
}

func DisplayLogoNewProject() {
	yellowStyle, blueStyle, _ := CreateStyles()

	fmt.Println()
	fmt.Println(yellowStyle.Render("Tamago::Fitaco()") + blueStyle.Render("🐣 1.0 "))
	fmt.Println()
}
