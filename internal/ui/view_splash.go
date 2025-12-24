package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) splashView() string {
	title := `
██████╗ ██╗   ██╗███╗   ██╗██████╗  ██████╗ ██╗    ██╗███╗   ██╗
██╔══██╗██║   ██║████╗  ██║██╔══██╗██╔═══██╗██║    ██║████╗  ██║
██████╔╝██║   ██║██╔██╗ ██║██║  ██║██║   ██║██║ █╗ ██║██╔██╗ ██║
██╔══██╗██║   ██║██║╚██╗██║██║  ██║██║   ██║██║███╗██║██║╚██╗██║
██║  ██║╚██████╔╝██║ ╚████║██████╔╝╚██████╔╝╚███╔███╔╝██║ ╚████║
╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═════╝  ╚═════╝  ╚══╝╚══╝ ╚═╝  ╚═══╝`

	content := theme.SplashTitleStyle.Render(title) + "\n\n" +
		m.LoadProg.ViewAs(m.LoadVal) + "\n" +
		theme.LogoStyle.Render(m.LoadingMsg)

	l := layout.New(m.Width, m.Height)

	// Center content horizontally
	centeredContent := lipgloss.Place(l.UsableWidth(), lipgloss.Height(content),
		lipgloss.Center, lipgloss.Top,
		content,
	)

	// Calculate vertical centering by filling to InnerHeight
	// Account for footer height (1 line) even though splash doesn't show footer
	innerHeight := l.InnerHeight()
	footerHeight := 1 // Footer is always 1 line tall
	availableHeight := innerHeight - footerHeight
	contentHeight := lipgloss.Height(centeredContent)

	// Add padding above and below to center vertically within available space
	topPadding := (availableHeight - contentHeight) / 2
	bottomPadding := availableHeight - contentHeight - topPadding

	if topPadding < 0 {
		topPadding = 0
	}
	if bottomPadding < 0 {
		bottomPadding = 0
	}

	// Build the full-height content using JoinVertical like other views

	return lipgloss.JoinVertical(lipgloss.Left,
		strings.Repeat("\n", topPadding),
		centeredContent,
		strings.Repeat("\n", bottomPadding+footerHeight),
	)
}
