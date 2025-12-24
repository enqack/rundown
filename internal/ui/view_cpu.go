package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
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

	cpuUsage := theme.MetricLabelStyle.Render(fmt.Sprintf("Global CPU Usage: %.2f%%", m.Stats.CPUUsage)) + "\n" + m.CPUProg.ViewAs(m.Stats.CPUUsage/100)
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(cpuUsage) + "\n\n")

	// Per-Core CPU Usage
	if len(m.Stats.CPUCores) > 0 {
		// Ensure we have enough progress bars
		for len(m.CPUCoresProgs) < len(m.Stats.CPUCores) {
			p := progress.New(progress.WithDefaultGradient())
			p.FullColor = string(theme.PrimaryColor)
			// Initial width - should be dynamic in Update, but good fallback
			leftBoxW, _ := l.SplitTwoColumns(l.UsableWidth(), 0)
			p.Width = l.GraphWidth(l.BoxContentWidth(leftBoxW))

			if m.Width <= 120 {
				p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
			}
			m.CPUCoresProgs = append(m.CPUCoresProgs, p)
		}

		var coreBoxes []string
		for i, coreUsage := range m.Stats.CPUCores {
			label := theme.MetricLabelStyle.Render(fmt.Sprintf("Core %d: %.1f%%", i, coreUsage))
			bar := m.CPUCoresProgs[i].ViewAs(coreUsage / 100)
			coreBoxes = append(coreBoxes, label+"\n"+bar)
		}

		// Grid Layout
		var rows []string
		leftBoxW, rightBoxW := l.SplitTwoColumns(l.UsableWidth(), 0)

		// Calculate the standard box width once to ensure consistency
		standardBoxWidth := l.BoxContentWidth(l.UsableWidth())

		if m.Width > 120 {
			// 2 Columns - use split widths directly
			for i := 0; i < len(coreBoxes); i += 2 {
				left := theme.BoxStyle.Width(l.BoxContentWidth(leftBoxW)).Render(coreBoxes[i])
				right := ""
				if i+1 < len(coreBoxes) {
					right = theme.BoxStyle.Width(l.BoxContentWidth(rightBoxW)).Render(coreBoxes[i+1])
				} else {
					right = theme.BoxStyle.Width(l.BoxContentWidth(rightBoxW)).Render("")
				}
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right))
			}
		} else {
			// 1 Column - use the same width as Global CPU box
			for _, box := range coreBoxes {
				rows = append(rows, theme.BoxStyle.Width(standardBoxWidth).Render(box))
			}
		}
		s.WriteString(strings.Join(rows, "\n") + "\n\n")
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
