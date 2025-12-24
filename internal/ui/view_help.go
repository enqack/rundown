package ui

import (
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) helpView() string {
	return m.HelpViewport.View()
}

func (m *Model) updateHelpViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Help & Keyboard Shortcuts") + "\n\n")

	l := layout.New(m.Width, m.Height)

	// Define key categories
	headers := []string{"KEY", "ACTION", "CONTEXT"}

	// Calculate widths
	// Key: ~10, Action: ~20, Context: ~20
	keyW := 15
	actionW := 30
	contextW := l.UsableWidth() - keyW - actionW - 6 // -6 for padding
	if contextW < 10 {
		contextW = 10
	}
	widths := []int{keyW, actionW, contextW}

	var rows [][]string

	// Global
	rows = append(rows, []string{"1-7", "Switch Tabs", "Global"})
	rows = append(rows, []string{"Tab", "Next Tab", "Global"})
	rows = append(rows, []string{"Shift+Tab", "Previous Tab", "Global"})
	rows = append(rows, []string{"?", "Toggle Help", "Global"})
	rows = append(rows, []string{"q / Ctrl+C", "Quit Application", "Global"})
	rows = append(rows, []string{"+ / -", "Adjust Update Speed", "Global"})

	// Navigation
	rows = append(rows, []string{"j / Down", "Scroll Down", "Scrollable Views"})
	rows = append(rows, []string{"k / Up", "Scroll Up", "Scrollable Views"})
	rows = append(rows, []string{"PgDn / PgUp", "Page Scroll", "Scrollable Views"})
	rows = append(rows, []string{"Home / End", "Jump to Top/Bottom", "Scrollable Views"})

	// Sorting - Process
	rows = append(rows, []string{"c", "Sort by CPU", "Process/CPU Tab"})
	rows = append(rows, []string{"m", "Sort by Memory", "Process/Mem Tab"})
	rows = append(rows, []string{"p", "Sort by PID", "Process Tab"})
	rows = append(rows, []string{"n", "Sort by Name", "Process Tab"})
	rows = append(rows, []string{"u", "Sort by User", "Process Tab"})
	rows = append(rows, []string{"t", "Sort by Time", "Process Tab"})

	// Sorting - Network
	rows = append(rows, []string{"f", "Sort by Foreign Addr", "Network Tab"})
	rows = append(rows, []string{"l", "Sort by Local Addr", "Network Tab"})
	rows = append(rows, []string{"s", "Sort by State", "Network Tab"})
	rows = append(rows, []string{"p", "Sort by Protocol", "Network Tab"})

	s.WriteString(m.renderTacticalTable("Keyboard Legend", headers, widths, rows))

	m.HelpViewport.SetContent(s.String())
}
