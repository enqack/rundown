package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) View() string {
	l := layout.New(m.Width, m.Height)
	if m.Loading {
		return theme.ContainerStyle.
			Width(l.UsableWidth()).
			Height(l.ContainerHeight()).
			Render(m.splashView())
	}

	var s strings.Builder

	// Tabs
	tabs := []string{"1:Over", "2:CPU", "3:Mem", "4:Disk", "5:Net", "6:Temp", "7:Proc"}
	var renderedTabs []string
	for i, t := range tabs {
		if Tab(i) == m.Tab {
			renderedTabs = append(renderedTabs, theme.ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, theme.InactiveTabStyle.Render(t))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	s.WriteString(lipgloss.NewStyle().MaxWidth(l.UsableWidth()).Render(tabRow) + "\n\n")

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
	case TabProc:
		body = m.procView()
	case TabHelp:
		body = m.helpView()
	}
	s.WriteString(body)

	// Combine Content and Footer
	footer := m.renderFooter()
	content := s.String()

	// Height available in the container
	innerHeight := l.InnerHeight()
	contentHeight := lipgloss.Height(content)
	footerHeight := lipgloss.Height(footer)

	gapHeight := innerHeight - contentHeight - footerHeight
	if gapHeight < 0 {
		gapHeight = 0
	}

	fullBlock := lipgloss.JoinVertical(lipgloss.Left,
		content,
		strings.Repeat("\n", gapHeight),
		footer,
	)

	// Final container wrap - ensure the outer border fits within terminal
	return theme.ContainerStyle.
		Width(l.UsableWidth()).
		Height(l.ContainerHeight()).
		MaxWidth(m.Width).
		MaxHeight(m.Height).
		Render(fullBlock)
}

func (m Model) getActiveViewport() viewport.Model {
	switch m.Tab {
	case TabOverview:
		return m.OverviewViewport
	case TabCPU:
		return m.CPUViewport
	case TabMem:
		return m.MemViewport
	case TabDisk:
		return m.DiskViewport
	case TabNet:
		return m.NetViewport
	case TabTemp:
		return m.TempViewport
	case TabProc:
		return m.ProcViewport
	case TabHelp:
		return m.HelpViewport
	default:
		return m.OverviewViewport
	}
}

func (m Model) renderFooter() string {
	l := layout.New(m.Width, m.Height)
	footerW := l.UsableWidth()

	// Format interval
	intervalMs := m.UpdateInterval.Milliseconds()
	var intervalText string
	if intervalMs >= 1000 {
		intervalText = fmt.Sprintf("%.1fs", float64(intervalMs)/1000.0)
	} else {
		intervalText = fmt.Sprintf("%dms", intervalMs)
	}

	leftStr := fmt.Sprintf("Last updated: %s | Interval: %s", m.LastSync.Format("15:04:05"), intervalText)

	// Scroll/Sector Indicator
	v := m.getActiveViewport()
	sectorInfo := ""
	if v.TotalLineCount() > v.Height && v.Height > 0 {
		total := v.TotalLineCount()
		// Calculate sector based on the bottom line of the viewport
		// This ensures that when we are at the bottom, we are in the last sector.
		bottomLine := v.YOffset + v.Height - 1
		sector := (bottomLine / v.Height) + 1
		maxSector := (total + v.Height - 1) / v.Height
		sectorInfo = fmt.Sprintf("[Sector %02d/%02d] ", sector, maxSector)
	}

	rightStr := fmt.Sprintf("%s?: Help", sectorInfo)

	left := lipgloss.NewStyle().Faint(true).Render(leftStr)
	right := lipgloss.NewStyle().Foreground(theme.SecondaryColor).Bold(true).Render(rightStr)

	// Simple Layout: Left -------- Right
	return lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		lipgloss.NewStyle().Width(footerW-lipgloss.Width(left)-lipgloss.Width(right)).Render(""),
		right,
	)
}

func truncate(s string, l int) string {
	if l <= 0 {
		return ""
	}
	s = strings.ReplaceAll(s, "\n", " ")
	visualW := lipgloss.Width(s)
	if visualW <= l {
		return s
	}

	// ANSI-safe truncation using lipgloss's internal ability
	return lipgloss.NewStyle().MaxWidth(l).Render(s)
}

func (m Model) renderTacticalBox(outerWidth int, content string) string {
	l := layout.New(m.Width, m.Height)
	innerW := l.BoxContentWidth(outerWidth)

	// Ensure content is truncated and PADDED to exactly innerW for every line.
	// This prevents the box from shrinking due to short labels (Host, IP, etc.).
	lines := strings.Split(content, "\n")
	var stableLines []string
	for _, line := range lines {
		// Use PlaceHorizontal to guarantee innerW width with padding if needed
		stableLines = append(stableLines, lipgloss.PlaceHorizontal(innerW, lipgloss.Left, truncate(line, innerW)))
	}
	stableContent := strings.Join(stableLines, "\n")

	// Render the box. In lipgloss, Style.Width() on a bordered style sets the CONTENT width.
	box := theme.BoxStyle.
		Width(innerW).
		Align(lipgloss.Left).
		Render(stableContent)

	// AGGRESSIVE WRAPPER: Force the produced string to be exactly outerWidth.
	// This accounts for any edge cases where lipgloss might miscalculate the frame.
	return lipgloss.NewStyle().Width(outerWidth).MaxWidth(outerWidth).Render(box)
}

func (m Model) renderTacticalRow(totalWidth int, gap int, boxes ...string) string {
	if len(boxes) == 0 {
		return ""
	}
	if len(boxes) == 1 {
		return lipgloss.NewStyle().Width(totalWidth).MaxWidth(totalWidth).Render(boxes[0])
	}

	// Create spacers for the gaps using the caller's gap size
	sep := strings.Repeat(" ", gap)
	var components []string
	for i, box := range boxes {
		components = append(components, box)
		if i < len(boxes)-1 {
			components = append(components, sep)
		}
	}

	joined := lipgloss.JoinHorizontal(lipgloss.Top, components...)
	return lipgloss.NewStyle().Width(totalWidth).MaxWidth(totalWidth).Render(joined)
}

// renderTacticalGauge renders a titled tactical box containing a progress bar with a same-line suffix.
// suffix: e.g. " %" or " °C"
func (m Model) renderTacticalGauge(outerW int, title string, p progress.Model, ratio float64, value float64, unit string) string {
	l := layout.New(m.Width, m.Height)
	contentW := l.BoxContentWidth(outerW)

	// Force internal percentage off
	p.ShowPercentage = false

	// Suffix formatting
	suffix := fmt.Sprintf(" %3.0f%s", value, unit)
	sW := lipgloss.Width(suffix)

	// Available space for the bar (leave 4 char safety margin to prevent wrapping)
	barW := contentW - sW - 4
	if barW < 0 {
		barW = 0
	}
	p.Width = barW

	// Directly build the row.
	// TrimSpace ensures we don't carry any accidental graph component newlines.
	barLine := strings.TrimSpace(p.ViewAs(ratio)) + suffix

	// Wrap in a MaxHeight(1) style to prevent any accidental multiline output
	rowStr := lipgloss.NewStyle().
		MaxHeight(1).
		Render(barLine)

	content := title + "\n" + rowStr
	return m.renderTacticalBox(outerW, content)
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
			style := lipgloss.NewStyle().Padding(0, 1).Width(widths[j]).MaxHeight(1)
			rowCells = append(rowCells, style.Render(truncate(cell, widths[j]-2)))
		}
		rowStr := lipgloss.JoinHorizontal(lipgloss.Top, rowCells...)
		if i%2 == 1 {
			s.WriteString(lipgloss.NewStyle().Background(lipgloss.Color("#24283b")).Render(rowStr) + "\n")
		} else {
			s.WriteString(rowStr + "\n")
		}
	}

	l := layout.New(m.Width, m.Height)
	// AGGRESSIVE WRAPPER: Force table box to be exactly UsableWidth total
	outerW := l.UsableWidth()
	box := theme.BoxStyle.Width(l.BoxContentWidth(outerW)).Render(s.String())
	return lipgloss.NewStyle().Width(outerW).MaxWidth(outerW).Render(box)
}

func getAdaptiveScale(usageBps, linkBps uint64) float64 {
	if linkBps == 0 {
		return 1024.0 // Fallback if link speed unknown
	}

	// Calculate usage ratio
	ratio := float64(usageBps) / float64(linkBps)

	// Define stable tiers as percentage of total link speed
	// 0.1%, 1%, 5%, 10%, 25%, 50%, 75%, 100%
	tiers := []float64{0.001, 0.01, 0.05, 0.10, 0.25, 0.50, 0.75, 1.0}

	for _, t := range tiers {
		if ratio <= t {
			return float64(linkBps) * t
		}
	}

	// Over 100% (burst or error), use actual usage
	return float64(usageBps)
}
