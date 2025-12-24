package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) netView() string {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Network Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)

	linkInfo := fmt.Sprintf("%s: %s  |  %s: %s  |  %s: %s",
		theme.MetricLabelStyle.Render("LINK SPEED"), stats.FormatBytes(m.Stats.LinkSpeed/8)+"/s",
		theme.MetricLabelStyle.Render("TOTAL SENT"), stats.FormatBytes(m.Stats.NetSent),
		theme.MetricLabelStyle.Render("TOTAL RECV"), stats.FormatBytes(m.Stats.NetRecv))
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(lipgloss.NewStyle().Width(l.BoxContentWidth(l.UsableWidth())).Render(linkInfo)) + "\n\n")

	scaleUp := getStepScale(m.Stats.NetSentDelta*8, m.Stats.LinkSpeed)
	scaleDn := getStepScale(m.Stats.NetRecvDelta*8, m.Stats.LinkSpeed)

	if m.Width > 120 {
		leftBoxW, rightBoxW := l.SplitTwoColumns(l.UsableWidth(), 0)

		m.NetUpProg.Width = l.GraphWidth(l.BoxContentWidth(leftBoxW))
		m.NetDnProg.Width = l.GraphWidth(l.BoxContentWidth(rightBoxW))

		netUp := theme.MetricLabelStyle.Render(fmt.Sprintf("Upload:   %s/s (Scale: %.0f%%)", stats.FormatBytes(m.Stats.NetSentDelta), scaleUp*100)) + "\n" + m.NetUpProg.ViewAs(float64(m.Stats.NetSentDelta*8)/(float64(m.Stats.LinkSpeed)*scaleUp))
		netDn := theme.MetricLabelStyle.Render(fmt.Sprintf("Download: %s/s (Scale: %.0f%%)", stats.FormatBytes(m.Stats.NetRecvDelta), scaleDn*100)) + "\n" + m.NetDnProg.ViewAs(float64(m.Stats.NetRecvDelta*8)/(float64(m.Stats.LinkSpeed)*scaleDn))

		upBox := theme.BoxStyle.Width(l.BoxContentWidth(leftBoxW)).Render(netUp)
		dnBox := theme.BoxStyle.Width(l.BoxContentWidth(rightBoxW)).Render(netDn)
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, upBox, "  ", dnBox) + "\n\n")
	} else {
		fullBoxW := l.BoxContentWidth(l.UsableWidth())
		m.NetUpProg.Width = l.GraphWidth(fullBoxW)
		m.NetDnProg.Width = l.GraphWidth(fullBoxW)

		netUp := theme.MetricLabelStyle.Render(fmt.Sprintf("Upload:   %s/s (Scale: %.0f%%)", stats.FormatBytes(m.Stats.NetSentDelta), scaleUp*100)) + "\n" + m.NetUpProg.ViewAs(float64(m.Stats.NetSentDelta*8)/(float64(m.Stats.LinkSpeed)*scaleUp))
		netDn := theme.MetricLabelStyle.Render(fmt.Sprintf("Download: %s/s (Scale: %.0f%%)", stats.FormatBytes(m.Stats.NetRecvDelta), scaleDn*100)) + "\n" + m.NetDnProg.ViewAs(float64(m.Stats.NetRecvDelta*8)/(float64(m.Stats.LinkSpeed)*scaleDn))

		// Stacked
		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(netUp) + "\n")
		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(netDn) + "\n\n")
	}

	headers := []string{"PROTO", "DIR", "STATE", "LOCAL", "FOREIGN", "PROGRAM"}
	// Fixed: 6 7 15

	availWidth := l.BoxContentWidth(l.UsableWidth())
	// Dynamic cols: 3 of them.
	// Helper supports 1 dynamic col...
	// We need 3 dynamic cols.
	// Let's optimize local calc.

	used := 6 + 7 + 15
	totalCols := 6
	paddingOverhead := totalCols * 2

	rem := availWidth - used - paddingOverhead
	if rem < 10 {
		rem = 10
	}
	w := rem / 3
	if w < 1 {
		w = 1
	}
	widths := []int{6, 7, 15, w, w, w}

	// Calculate height of top content dynamically
	topContent := s.String()
	usedHeight := lipgloss.Height(topContent)

	// Table Overhead: Title(1) + Blank(2) + Header(3)
	tableOverhead := 6

	connLimit := l.UsableHeight() - usedHeight - tableOverhead
	if connLimit < 1 {
		connLimit = 1
	}
	var rows [][]string
	for _, c := range m.Stats.Connections {
		pname := c.PName
		if pname == "" {
			pname = "-"
		}
		rows = append(rows, []string{
			c.Proto,
			c.Direction,
			c.State,
			c.Local,
			c.Remote,
			fmt.Sprintf("%d/%s", c.PID, pname),
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		switch m.SortBy {
		case "local":
			return rows[i][3] < rows[j][3]
		case "foreign":
			return rows[i][4] < rows[j][4]
		case "state":
			return rows[i][2] < rows[j][2]
		case "name":
			return rows[i][5] < rows[j][5]
		case "pid":
			return rows[i][5] < rows[j][5]
		default:
			return false
		}
	})

	if len(rows) > connLimit {
		rows = rows[:connLimit]
	}

	s.WriteString(m.renderTacticalTable("Active Network Connections", headers, widths, rows))

	return s.String()
}
