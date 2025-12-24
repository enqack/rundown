package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) cpuView() string {
	return m.CPUViewport.View()
}

func (m *Model) updateCpuViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed CPU Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)

	// Set width before viewing
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), "Global CPU Usage", m.CPUProg, m.Stats.CPUUsage/100, m.Stats.CPUUsage, "%") + "\n")

	// Per-Core CPU Usage
	if len(m.Stats.CPUCores) > 0 {
		if len(m.CPUCoresProgs) == 0 {
			for range m.Stats.CPUCores {
				p := theme.NewThemedProgress()
				m.CPUCoresProgs = append(m.CPUCoresProgs, p)
			}
		}

		// Grid Layout
		var rows []string
		leftBoxW, rightBoxW := l.SplitTwoColumns(l.UsableWidth(), 2)

		if m.Width > 120 {
			// 2 Columns
			for i := 0; i < len(m.Stats.CPUCores); i += 2 {
				// Left Box
				left := m.renderTacticalGauge(leftBoxW, fmt.Sprintf("Core %d", i), m.CPUCoresProgs[i], m.Stats.CPUCores[i]/100, m.Stats.CPUCores[i], "%")

				// Right Box
				var right string
				if i+1 < len(m.Stats.CPUCores) {
					right = m.renderTacticalGauge(rightBoxW, fmt.Sprintf("Core %d", i+1), m.CPUCoresProgs[i+1], m.Stats.CPUCores[i+1]/100, m.Stats.CPUCores[i+1], "%")
				} else {
					right = m.renderTacticalBox(rightBoxW, "")
				}
				rows = append(rows, m.renderTacticalRow(l.UsableWidth(), 2, left, right))
			}
		} else {
			// 1 Column
			for i, usage := range m.Stats.CPUCores {
				rows = append(rows, m.renderTacticalGauge(l.UsableWidth(), fmt.Sprintf("Core %d", i), m.CPUCoresProgs[i], usage/100, usage, "%"))
			}
		}
		s.WriteString(strings.Join(rows, "\n") + "\n")
	}

	headers := []string{"PID", "USER", "%CPU", "%MEM", "TIME+", "COMMAND"}
	// Fixed cols: 10, 10, 8, 8, 10. Padding: 2 per col.
	fixedWidths := []int{10, 10, 8, 8, 10}

	// Available width for the table (inside Box, so usableWidth - BoxFrame)
	availWidth := l.BoxContentWidth(l.UsableWidth())
	cmdW := l.CalculateTableDynamicWidth(availWidth, fixedWidths, 2)

	widths := []int{10, 10, 8, 8, 10, cmdW}
	var rows [][]string
	// No limit for scrolling view, show all 10 top processes (or whatever stats returns)
	for _, p := range m.Stats.TopCPU {
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

	s.WriteString(m.renderTacticalTable("Top CPU Processes", headers, widths, rows))

	m.CPUViewport.SetContent(s.String())
}
