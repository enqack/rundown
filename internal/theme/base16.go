package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines all colors used in the application
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Warning    lipgloss.Color
	Critical   lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
}

// Base16 returns a theme using terminal's base16 ANSI colors
func Base16() Theme {
	return Theme{
		Name:       "base16",
		Primary:    lipgloss.Color("12"), // Bright Blue (base0D)
		Secondary:  lipgloss.Color("10"), // Bright Green (base0B)
		Accent:     lipgloss.Color("13"), // Bright Magenta (base0E)
		Warning:    lipgloss.Color("11"), // Bright Yellow (base0A)
		Critical:   lipgloss.Color("9"),  // Bright Red (base08)
		Background: lipgloss.Color("0"),  // Black (base00)
		Foreground: lipgloss.Color("7"),  // White (base05)
	}
}
