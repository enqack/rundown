package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) memView() string {
	return m.MemViewport.View()
}

func (m *Model) updateMemViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Memory Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)

	memHeader := fmt.Sprintf("Memory Usage: %s / %s",
		stats.FormatBytes(m.Stats.UsedMemory),
		stats.FormatBytes(m.Stats.TotalMemory))
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), memHeader, m.MemProg, m.Stats.MemoryUsage/100, m.Stats.MemoryUsage, "%") + "\n")

	// Swap Memory
	if m.Stats.SwapTotal > 0 {
		swapHeader := "Swap Usage:"

		// Inline Detailed Swap Devices
		for _, dev := range m.Stats.SwapDevices {
			swapHeader += fmt.Sprintf(" (%s) %s / %s",
				dev.Name,
				stats.FormatBytes(dev.UsedBytes),
				stats.FormatBytes(dev.TotalBytes))
		}

		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), swapHeader, m.SwapProg, m.Stats.SwapPercent/100, m.Stats.SwapPercent, "%") + "\n")

	}

	headers := []string{"PID", "USER", "%CPU", "%MEM", "TIME+", "COMMAND"}
	fixedWidths := []int{10, 10, 8, 8, 10}

	availWidth := l.BoxContentWidth(l.UsableWidth())
	cmdW := l.CalculateTableDynamicWidth(availWidth, fixedWidths, 2)

	widths := []int{10, 10, 8, 8, 10, cmdW}
	var rows [][]string

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

	s.WriteString(m.renderTacticalTable("Top Memory Processes", headers, widths, rows))

	m.MemViewport.SetContent(s.String())
}
