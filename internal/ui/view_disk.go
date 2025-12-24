package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) diskView() string {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Disk Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)
	for _, d := range m.Stats.Disks {
		diskHeader := theme.MetricLabelStyle.Render(fmt.Sprintf("Disk (%s): %s / %s (%.1f%%)",
			d.MountPoint,
			stats.FormatBytes(d.Used),
			stats.FormatBytes(d.Total),
			d.Usage))
		dp, ok := m.DiskProgs[d.MountPoint]
		if !ok {
			dp = progress.New(progress.WithDefaultGradient())
			dp.FullColor = string(theme.AccentColor)
			// Ensure it spans the page with safety margin
			dp.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
			m.DiskProgs[d.MountPoint] = dp
		}
		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(diskHeader+"\n"+dp.ViewAs(d.Usage/100)) + "\n")
	}

	// Swap (displayed as virtual disk)
	if m.Stats.SwapTotal > 0 {
		swapHeader := theme.MetricLabelStyle.Render(fmt.Sprintf("Swap Memory: %s / %s (%.1f%%)",
			stats.FormatBytes(m.Stats.SwapUsed),
			stats.FormatBytes(m.Stats.SwapTotal),
			m.Stats.SwapPercent))
		s.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(swapHeader+"\n"+m.SwapProg.ViewAs(m.Stats.SwapPercent/100)) + "\n")
	}

	return s.String()
}
