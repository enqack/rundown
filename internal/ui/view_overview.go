package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) overviewView() string {
	var s strings.Builder

	l := layout.New(m.Width, m.Height)
	// Info section
	info := fmt.Sprintf("Host: %s | IP: %s | Uptime: %s",
		m.Stats.HostName, m.Stats.IPAddress, stats.FormatDuration(m.Stats.Uptime))
	s.WriteString(theme.MetricLabelStyle.Render(info) + "\n\n")

	// CPU Overview
	cpuHeader := theme.MetricLabelStyle.Render(fmt.Sprintf("CPU Usage: %.1f%%", m.Stats.CPUUsage))
	cpuBar := m.CPUProg.ViewAs(m.Stats.CPUUsage / 100)
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(cpuHeader+"\n"+cpuBar) + "\n")

	// Memory Overview
	memHeader := theme.MetricLabelStyle.Render(fmt.Sprintf("Memory: %s / %s (%.1f%%)",
		stats.FormatBytes(m.Stats.UsedMemory),
		stats.FormatBytes(m.Stats.TotalMemory),
		m.Stats.MemoryUsage))
	memBar := m.MemProg.ViewAs(m.Stats.MemoryUsage / 100)
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(memHeader+"\n"+memBar) + "\n")

	// Total Disk Overview (Grand Total)
	totalDiskHeader := theme.MetricLabelStyle.Render(fmt.Sprintf("Disk Usage: %s / %s (%.1f%%)",
		stats.FormatBytes(m.Stats.DiskUsed),
		stats.FormatBytes(m.Stats.DiskTotal),
		float64(m.Stats.DiskUsed)/float64(m.Stats.DiskTotal)*100))
	totalDiskBar := m.DiskProg.ViewAs(float64(m.Stats.DiskUsed) / float64(m.Stats.DiskTotal))
	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(totalDiskHeader+"\n"+totalDiskBar) + "\n")

	// Network Overview - Grouped
	scaleUp := getStepScale(m.Stats.NetSentDelta*8, m.Stats.LinkSpeed)
	scaleDn := getStepScale(m.Stats.NetRecvDelta*8, m.Stats.LinkSpeed)

	// Single line of info above graph
	netUpInfo := fmt.Sprintf("Egress: %s/s (Scale: %.0f%%) | Total: %s | Link: %s/s",
		stats.FormatBytes(m.Stats.NetSentDelta),
		scaleUp*100,
		stats.FormatBytes(m.Stats.NetSent),
		stats.FormatBytes(m.Stats.LinkSpeed/8))

	netDnInfo := fmt.Sprintf("Ingress: %s/s (Scale: %.0f%%) | Total: %s | Link: %s/s",
		stats.FormatBytes(m.Stats.NetRecvDelta),
		scaleDn*100,
		stats.FormatBytes(m.Stats.NetRecv),
		stats.FormatBytes(m.Stats.LinkSpeed/8))

	if m.Width > 120 {
		// Split View
		netUpBar := m.NetUpProg.ViewAs(float64(m.Stats.NetSentDelta*8) / (float64(m.Stats.LinkSpeed) * scaleUp))
		netDnBar := m.NetDnProg.ViewAs(float64(m.Stats.NetRecvDelta*8) / (float64(m.Stats.LinkSpeed) * scaleDn))

		leftW, rightW := l.SplitTwoColumns(l.UsableWidth(), 0)
		upBox := theme.BoxStyle.Width(l.BoxContentWidth(leftW)).Render(theme.MetricLabelStyle.Render(netUpInfo) + "\n" + netUpBar)
		dnBox := theme.BoxStyle.Width(l.BoxContentWidth(rightW)).Render(theme.MetricLabelStyle.Render(netDnInfo) + "\n" + netDnBar)

		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, upBox, "  ", dnBox) + "\n")

	} else {
		// Stacked View
		netUpBar := m.NetUpProg.ViewAs(float64(m.Stats.NetSentDelta*8) / (float64(m.Stats.LinkSpeed) * scaleUp))
		netDnBar := m.NetDnProg.ViewAs(float64(m.Stats.NetRecvDelta*8) / (float64(m.Stats.LinkSpeed) * scaleDn))

		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(
			theme.MetricLabelStyle.Render(netUpInfo)+"\n"+netUpBar+"\n\n"+
				theme.MetricLabelStyle.Render(netDnInfo)+"\n"+netDnBar) + "\n")
	}

	thermal := theme.MetricLabelStyle.Render("Thermal Sensors") + "\n\n"
	thermalCpuBar := m.getTempProg("CPU Package").ViewAs(m.Stats.CPUTemp / 100)
	thermal += fmt.Sprintf("CPU Package: %.1f°C\n%s\n", m.Stats.CPUTemp, thermalCpuBar)

	for i, g := range m.Stats.GPUTemps {
		gname := fmt.Sprintf("GPU %d (%s)", i, g.SensorKey)
		gBar := m.getTempProg(gname).ViewAs(g.Temperature / 100)
		thermal += fmt.Sprintf("\n%s: %.1f°C\n%s\n", gname, g.Temperature, gBar)
	}

	s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(lipgloss.NewStyle().Width(l.BoxContentWidth(l.BoxContentWidth(l.UsableWidth()))).Render(thermal)) + "\n")

	return s.String()
}
