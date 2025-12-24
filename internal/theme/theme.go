package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	// Current active theme
	current Theme

	// Color exports (for backward compatibility)
	PrimaryColor   lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	WarningColor   lipgloss.Color
	CriticalColor  lipgloss.Color
	BgColor        lipgloss.Color
	FgColor        lipgloss.Color

	// Styles
	TitleStyle         lipgloss.Style
	BoxStyle           lipgloss.Style
	MetricLabelStyle   lipgloss.Style
	MetricValueStyle   lipgloss.Style
	ActiveTabStyle     lipgloss.Style
	InactiveTabStyle   lipgloss.Style
	ContainerStyle     lipgloss.Style
	TableHeaderStyle   lipgloss.Style
	TableSelectedStyle lipgloss.Style
	SplashTitleStyle   lipgloss.Style
	LogoStyle          lipgloss.Style
)

// Init initializes the theme system with the specified theme name
func Init(themeName string) {
	switch themeName {
	case "base16":
		current = Base16()
	case "cyberpunk":
		current = Cyberpunk()
	case "monochrome":
		current = Monochrome()
	case "phosphor":
		current = Phosphor()
	default:
		current = Base16() // Default theme
	}

	// Update color exports
	PrimaryColor = current.Primary
	SecondaryColor = current.Secondary
	AccentColor = current.Accent
	WarningColor = current.Warning
	CriticalColor = current.Critical
	BgColor = current.Background
	FgColor = current.Foreground

	// Initialize styles
	initStyles()
}

func initStyles() {
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor)

	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1)

	MetricLabelStyle = lipgloss.NewStyle().
		Foreground(FgColor).
		Bold(true)

	MetricValueStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor).
		Bold(true)

	ActiveTabStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		Background(BgColor).
		Padding(0, 2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(PrimaryColor)

	InactiveTabStyle = lipgloss.NewStyle().
		Foreground(FgColor).
		Padding(0, 2)

	ContainerStyle = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(PrimaryColor)
	// No padding - maximizes width while keeping border

	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(PrimaryColor)

	TableSelectedStyle = lipgloss.NewStyle().
		Foreground(BgColor).
		Background(PrimaryColor).
		Bold(true)

	SplashTitleStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true).
		MarginBottom(1)

	LogoStyle = lipgloss.NewStyle().
		Foreground(AccentColor).
		Bold(true)
}

func Gradient(start, end string, steps int) []string {
	c1, _ := colorful.Hex(start)
	c2, _ := colorful.Hex(end)
	var colors []string
	for i := 0; i < steps; i++ {
		t := float64(i) / float64(steps-1)
		colors = append(colors, c1.BlendHcl(c2, t).Clamped().Hex())
	}
	return colors
}
