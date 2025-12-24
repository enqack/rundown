package theme

import "github.com/charmbracelet/lipgloss"

// Monochrome returns a simple black and white theme
func Monochrome() Theme {
	return Theme{
		Name:       "monochrome",
		Primary:    lipgloss.Color("#FFFFFF"), // White
		Secondary:  lipgloss.Color("#CCCCCC"), // Light gray
		Accent:     lipgloss.Color("#AAAAAA"), // Medium gray
		Warning:    lipgloss.Color("#FFFFFF"), // White
		Critical:   lipgloss.Color("#FFFFFF"), // White
		Background: lipgloss.Color("#000000"), // Black
		Foreground: lipgloss.Color("#FFFFFF"), // White
	}
}
