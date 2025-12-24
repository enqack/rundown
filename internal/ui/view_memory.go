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

func (m Model) memView() string {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Memory Statistics") + "\n\n")

	memUsage := theme.MetricLabelStyle.Render(fmt.Sprintf("Used: %s / %s (%.1f%%)",
		stats.FormatBytes(m.Stats.UsedMemory),
		stats.FormatBytes(m.Stats.TotalMemory),
		m.Stats.MemoryUsage)) + "\n" + m.MemProg.ViewAs(m.Stats.MemoryUsage/100)
	l := layout.New(m.Width, m.Height)
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(memUsage) + "\n\n")

	// Swap Memory
	if m.Stats.SwapTotal > 0 {
		swapUsage := theme.MetricLabelStyle.Render(fmt.Sprintf("Swap: %s / %s (%.1f%%)",
			stats.FormatBytes(m.Stats.SwapUsed),
			stats.FormatBytes(m.Stats.SwapTotal),
			m.Stats.SwapPercent)) + "\n" + m.SwapProg.ViewAs(m.Stats.SwapPercent/100)
		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(swapUsage) + "\n\n")
	}

	headers := []string{"PID", "USER", "%CPU", "%MEM", "TIME+", "COMMAND"}
	// Fixed cols: 10, 10, 8, 8, 10. Padding: 2 per col.
	fixedWidths := []int{10, 10, 8, 8, 10}

	availWidth := l.BoxContentWidth(l.UsableWidth())
	cmdW := l.CalculateTableDynamicWidth(availWidth, fixedWidths, 2)

	widths := []int{10, 10, 8, 8, 10, cmdW}
	var rows [][]string
	// Calculate height of top content to determine remaining space for table
	topContent := s.String()
	// renderTacticalTable adds title (1) + blank (2) + header (1+border?)

	// Actually, let's just measure what we have so far.
	usedHeight := lipgloss.Height(topContent)

	// Table Overhead: Title(1) + Blank(2) + Header(3 includes border/padding) + FooterGap?
	tableOverhead := 6

	limit := l.UsableHeight() - usedHeight - tableOverhead
	if limit < 1 {
		limit = 1
	}
	for _, p := range m.Stats.TopMem {
		rows = append(rows, []string{
			fmt.Sprintf("%d", p.PID),
			p.User,
			fmt.Sprintf("%.1f", p.CPU),
			fmt.Sprintf("%.1f", p.Memory),
			p.Time,
			p.Cmdline,
		})
	}

	sort.SliceStable(rows, func(i, j int) bool {
		switch m.SortBy {
		case "pid":
			return rows[i][0] < rows[j][0]
		case "user":
			return rows[i][1] < rows[j][1]
		case "cpu":
			return strings.Compare(rows[i][2], rows[j][2]) > 0
		case "mem":
			return strings.Compare(rows[i][3], rows[j][3]) > 0
		case "time":
			return rows[i][4] > rows[j][4]
		case "name":
			return rows[i][5] < rows[j][5]
		default:
			return false
		}
	})

	if len(rows) > limit {
		rows = rows[:limit]
	}

	s.WriteString(m.renderTacticalTable("Top Memory Processes", headers, widths, rows))

	return s.String()
}
