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
	return m.OverviewViewport.View()
}

func (m *Model) updateOverviewViewport() {
	var s strings.Builder

	l := layout.New(m.Width, m.Height)
	s.WriteString(theme.TitleStyle.Render("Overview Statistics") + "\n\n")

	// Host/IP/Uptime Info Boxes (Standardized formatting)
	hostInfo := fmt.Sprintf("%s: %s", theme.MetricLabelStyle.Render("Host"), m.Stats.HostName)
	ipInfo := fmt.Sprintf("%s: %s", theme.MetricLabelStyle.Render("IP"), m.Stats.IPAddress)
	upInfo := fmt.Sprintf("%s: %s", theme.MetricLabelStyle.Render("Uptime"), stats.FormatDuration(m.Stats.Uptime))

	if m.Width > 120 {
		lw, mw, rw := l.SplitThreeColumns(l.UsableWidth(), 2)
		hBox := m.renderTacticalBox(lw, hostInfo)
		iBox := m.renderTacticalBox(mw, ipInfo)
		uBox := m.renderTacticalBox(rw, upInfo)
		s.WriteString(m.renderTacticalRow(l.UsableWidth(), 2, hBox, iBox, uBox) + "\n")
	} else {
		fullW := l.UsableWidth()
		content := lipgloss.JoinVertical(lipgloss.Left, hostInfo, ipInfo, upInfo)
		s.WriteString(m.renderTacticalBox(fullW, content) + "\n")
	}

	// Process Status (5 Columns)
	// [Total] [Running] [Sleeping] [Stopped] [Zombie]
	procWidths := l.SplitColumns(l.UsableWidth(), 5, 1) // 5 columns, 1 char gap

	// Ensure we have correct widths even if very narrow
	if len(procWidths) == 5 {
		bTotal := m.renderTacticalBox(procWidths[0], fmt.Sprintf("Total Processes: %d", m.Stats.ProcTotal))
		bRun := m.renderTacticalBox(procWidths[1], fmt.Sprintf("Running Processes: %d", m.Stats.ProcRunning))
		bSleep := m.renderTacticalBox(procWidths[2], fmt.Sprintf("Sleeping Processes: %d", m.Stats.ProcSleeping))
		bStop := m.renderTacticalBox(procWidths[3], fmt.Sprintf("Stopped Processes: %d", m.Stats.ProcStopped))
		bZom := m.renderTacticalBox(procWidths[4], fmt.Sprintf("Zombie Processes: %d", m.Stats.ProcZombie))

		s.WriteString(m.renderTacticalRow(l.UsableWidth(), 1, bTotal, bRun, bSleep, bStop, bZom) + "\n")
	} else {
		// Fallback for extremely narrow terminals
		s.WriteString(m.renderTacticalBox(l.UsableWidth(), fmt.Sprintf("Procs: %d", m.Stats.ProcTotal)) + "\n")
	}

	// CPU Overview
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), "CPU Usage", m.CPUProg, m.Stats.CPUUsage/100, m.Stats.CPUUsage, "%") + "\n")

	// Memory Overview
	memHeader := fmt.Sprintf("Memory: %s / %s",
		stats.FormatBytes(m.Stats.UsedMemory),
		stats.FormatBytes(m.Stats.TotalMemory))
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), memHeader, m.MemProg, m.Stats.MemoryUsage/100, m.Stats.MemoryUsage, "%") + "\n")

	// Total Disk Overview (Grand Total)
	diskHeader := fmt.Sprintf("Disk Usage: %s / %s",
		stats.FormatBytes(m.Stats.DiskUsed),
		stats.FormatBytes(m.Stats.DiskTotal))
	diskRatio := float64(m.Stats.DiskUsed) / float64(m.Stats.DiskTotal)
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), diskHeader, m.DiskProg, diskRatio, diskRatio*100, "%") + "\n")

	// Network Overview - Grouped
	scaleUp := getStepScale(m.Stats.NetSentDelta*8, m.Stats.LinkSpeed)
	scaleDn := getStepScale(m.Stats.NetRecvDelta*8, m.Stats.LinkSpeed)

	netUpTitle := fmt.Sprintf("Egress (Scale: %s/s)", stats.FormatBytes(uint64(float64(m.Stats.LinkSpeed/8)*scaleUp)))
	netDnTitle := fmt.Sprintf("Ingress (Scale: %s/s)", stats.FormatBytes(uint64(float64(m.Stats.LinkSpeed/8)*scaleDn)))

	upRatio := float64(m.Stats.NetSentDelta*8) / (float64(m.Stats.LinkSpeed) * scaleUp)
	dnRatio := float64(m.Stats.NetRecvDelta*8) / (float64(m.Stats.LinkSpeed) * scaleDn)

	if m.Width > 120 {
		leftW, rightW := l.SplitTwoColumns(l.UsableWidth(), 2)
		upBox := m.renderTacticalGauge(leftW, netUpTitle, m.NetUpProg, upRatio, upRatio*100, "%")
		dnBox := m.renderTacticalGauge(rightW, netDnTitle, m.NetDnProg, dnRatio, dnRatio*100, "%")
		s.WriteString(m.renderTacticalRow(l.UsableWidth(), 2, upBox, dnBox) + "\n")
	} else {
		upBar := m.renderTacticalGauge(l.UsableWidth(), netUpTitle, m.NetUpProg, upRatio, upRatio*100, "%")
		dnBar := m.renderTacticalGauge(l.UsableWidth(), netDnTitle, m.NetDnProg, dnRatio, dnRatio*100, "%")
		s.WriteString(upBar + "\n" + dnBar + "\n")
	}

	// Thermal Sensors (Individual Gauges for standard height)
	pCPU := m.getTempProg("CPU Package")
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), "CPU Thermal Package", pCPU, m.Stats.CPUTemp/110.0, m.Stats.CPUTemp, "°C") + "\n")

	for i, g := range m.Stats.GPUTemps {
		gname := fmt.Sprintf("GPU %d (%s)", i, g.SensorKey)
		pGPU := m.getTempProg(gname)
		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), gname, pGPU, g.Temperature/110.0, g.Temperature, "°C") + "\n")
	}

	m.OverviewViewport.SetContent(s.String())
}
