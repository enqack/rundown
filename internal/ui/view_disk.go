package ui

import (
	"fmt"
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) diskView() string {
	return m.DiskViewport.View()
}

func (m *Model) updateDiskViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Disk Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)

	// Total Disk Overview (Grand Total)
	totalDiskTag := fmt.Sprintf("Total Disk Usage: %s / %s",
		stats.FormatBytes(m.Stats.DiskUsed),
		stats.FormatBytes(m.Stats.DiskTotal))
	diskTotalRatio := float64(m.Stats.DiskUsed) / float64(m.Stats.DiskTotal)
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), totalDiskTag, m.DiskProg, diskTotalRatio, diskTotalRatio*100, "%") + "\n")

	for _, d := range m.Stats.Disks {
		// Label: Disk <Mount>
		diskTag := fmt.Sprintf("Disk %s", d.MountPoint)
		// Inline Device Node
		diskTag += fmt.Sprintf(" (%s) %s / %s",
			d.Device,
			stats.FormatBytes(d.Used),
			stats.FormatBytes(d.Total))

		dp, ok := m.DiskProgs[d.MountPoint]
		if !ok {
			dp = theme.NewThemedProgress()
			m.DiskProgs[d.MountPoint] = dp
		}
		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), diskTag, dp, d.Usage/100, d.Usage, "%") + "\n")
	}

	// Swap (displayed as virtual disk)
	if m.Stats.SwapTotal > 0 {
		swapTag := "Swap Memory"
		for _, dev := range m.Stats.SwapDevices {
			swapTag += fmt.Sprintf(" (%s) %s / %s",
				dev.Name,
				stats.FormatBytes(dev.UsedBytes),
				stats.FormatBytes(dev.TotalBytes))
		}
		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), swapTag, m.SwapProg, m.Stats.SwapPercent/100, m.Stats.SwapPercent, "%") + "\n")
	}

	m.DiskViewport.SetContent(s.String())
}
