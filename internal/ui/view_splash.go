package ui

import (
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
		theme.LogoStyle.Render(m.LoadingMsg) + "\n"

	// Center horizontally and vertically
	l := layout.New(m.Width, m.Height)

	return lipgloss.Place(l.UsableWidth(), l.InnerHeight(),
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
