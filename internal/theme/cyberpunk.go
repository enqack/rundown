package theme

import "github.com/charmbracelet/lipgloss"

// Cyberpunk returns the original vibrant purple/green/pink theme
func Cyberpunk() Theme {
	return Theme{
		Name:       "cyberpunk",
		Primary:    lipgloss.Color("#7D56F4"), // Purple
		Secondary:  lipgloss.Color("#04B575"), // Green
		Accent:     lipgloss.Color("#EE6FF8"), // Pink
		Warning:    lipgloss.Color("#FF8C00"), // Orange
		Critical:   lipgloss.Color("#FF4444"), // Red
		Background: lipgloss.Color("#1A1B26"), // Dark blue-gray
		Foreground: lipgloss.Color("#C0CAF5"), // Light blue-white
	}
}
