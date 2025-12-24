package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) View() string {
	if m.Loading {
		l := layout.New(m.Width, m.Height)
		return theme.ContainerStyle.
			Width(l.UsableWidth()).
			Height(l.ContainerHeight()).
			Render(m.splashView())
	}

	var s strings.Builder

	// Tabs
	tabs := []string{"1:Overview", "2:CPU", "3:Memory", "4:Disk", "5:Network", "6:Thermal"}
	var renderedTabs []string
	for i, t := range tabs {
		if Tab(i) == m.Tab {
			renderedTabs = append(renderedTabs, theme.ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, theme.InactiveTabStyle.Render(t))
		}
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...) + "\n\n")

	var body string
	switch m.Tab {
	case TabOverview:
		body = m.overviewView()
	case TabCPU:
		body = m.cpuView()
	case TabMem:
		body = m.memView()
	case TabDisk:
		body = m.diskView()
	case TabNet:
		body = m.netView()
	case TabTemp:
		body = m.tempView()
	}

	s.WriteString(body)

	// Combine Content and Footer
	footer := m.renderFooter()
	content := s.String()

	// Create a filled block with footer pinned to bottom

	// Inner Height available in the container
	l := layout.New(m.Width, m.Height)
	innerHeight := l.InnerHeight()

	// Total content height so far (Tabs + Body)
	contentHeight := lipgloss.Height(content)
	footerHeight := lipgloss.Height(footer) // Should be 1 typically

	gapHeight := innerHeight - contentHeight - footerHeight
	if gapHeight < 0 {
		gapHeight = 0
	}

	fullBlock := lipgloss.JoinVertical(lipgloss.Left,
		content,
		strings.Repeat("\n", gapHeight),
		footer,
	)

	// Final container wrap
	return theme.ContainerStyle.
		Width(l.UsableWidth()).
		Height(l.ContainerHeight()).
		Render(fullBlock)
}

func (m Model) renderFooter() string {
	var keys []string
	keys = append(keys, "Tab: Switch", "1-6: Select")

	switch m.Tab {
	case TabMem:
		keys = append(keys, "Sort: c/m/p/n/u/t/v/r/s")
	case TabCPU:
		keys = append(keys, "Scroll: j/k", "Sort: c/m/p/n/u/t/v/r/s")
	case TabNet:
		keys = append(keys, "Sort: l/f/p/n/s")
	case TabTemp:
		keys = append(keys, "Scroll: j/k, PgUp/PgDn")
	}

	keys = append(keys, "+/-: Interval", "q: Quit")

	// Format interval
	intervalMs := m.UpdateInterval.Milliseconds()
	var intervalText string
	if intervalMs >= 1000 {
		intervalText = fmt.Sprintf("%.1fs", float64(intervalMs)/1000.0)
	} else {
		intervalText = fmt.Sprintf("%dms", intervalMs)
	}

	footer := fmt.Sprintf("Last Check: %s | Interval: %s | %s", m.LastSync.Format("15:04:05.000"), intervalText, strings.Join(keys, " | "))
	return lipgloss.NewStyle().Faint(true).Render(footer)
}

func truncate(s string, l int) string {
	if l <= 0 {
		return ""
	}
	if len(s) > l {
		return s[:l-3] + "..."
	}
	return s
}

func (m Model) renderTacticalTable(title string, headers []string, widths []int, data [][]string) string {
	var s strings.Builder
	s.WriteString(theme.MetricLabelStyle.Render(title) + "\n\n")

	// Build Header
	var headerCells []string
	for i, h := range headers {
		style := theme.TableHeaderStyle.Width(widths[i]).Padding(0, 1).MaxHeight(1)
		headerCells = append(headerCells, style.Render(truncate(h, widths[i]-2)))
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, headerCells...) + "\n")

	// Build Rows
	for i, row := range data {
		var rowCells []string
		for j, cell := range row {
			cellStr := truncate(cell, widths[j]-2)
			style := lipgloss.NewStyle().Padding(0, 1).Width(widths[j]).MaxHeight(1)
			rowCells = append(rowCells, style.Render(cellStr))
		}
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, rowCells...)
		if i%2 == 1 {
			s.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("#24283b")).Render(rowStr) + "\n")
		} else {
			s.WriteString(rowStr + "\n")
		}
	}

	return theme.BoxStyle.Width(layout.New(m.Width, m.Height).BoxContentWidth(layout.New(m.Width, m.Height).UsableWidth())).Render(s.String())
}

func getStepScale(usageBps, linkBps uint64) float64 {
	ratio := float64(usageBps) / float64(linkBps)
	if ratio <= 0.25 {
		return 0.25
	} else if ratio <= 0.50 {
		return 0.50
	} else if ratio <= 0.75 {
		return 0.75
	}
	return 1.0
}
