package theme

import "github.com/charmbracelet/lipgloss"

// Phosphor returns a retro green phosphor monochrome monitor theme
func Phosphor() Theme {
	return Theme{
		Name:       "phosphor",
		Primary:    lipgloss.Color("#33FF33"), // Bright green
		Secondary:  lipgloss.Color("#00FF00"), // Pure green
		Accent:     lipgloss.Color("#66FF66"), // Lighter green
		Warning:    lipgloss.Color("#AAFF00"), // Yellow-green
		Critical:   lipgloss.Color("#FFFF00"), // Yellow (amber warning)
		Background: lipgloss.Color("#001100"), // Very dark green
		Foreground: lipgloss.Color("#00DD00"), // Medium green
	}
}
