package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m *Model) updateTempViewport() {
	var content strings.Builder
	content.WriteString(theme.TitleStyle.Render("System Temperature Sensors") + "\n\n")

	allTemps := append([]stats.TempInfo{}, m.Stats.Temperatures...)
	for _, g := range m.Stats.GPUTemps {
		found := false
		for _, t := range allTemps {
			if t.SensorKey == g.SensorKey {
				found = true
				break
			}
		}
		if !found {
			allTemps = append(allTemps, g)
		}
	}

	l := layout.New(m.Width, m.Height)
	for _, t := range allTemps {
		sensorName := t.SensorKey
		for i, g := range m.Stats.GPUTemps {
			if g.SensorKey == t.SensorKey {
				sensorName = fmt.Sprintf("GPU %d (%s)", i, g.SensorKey)
				break
			}
		}
		if sensorName == t.SensorKey && (strings.Contains(strings.ToLower(t.SensorKey), "package") || t.Temperature == m.Stats.CPUTemp) {
			sensorName = fmt.Sprintf("CPU (%s)", t.SensorKey)
		}

		p := m.getTempProg(sensorName)
		// Box Width = usableWidth (Inner Container)
		// Progress Width = usableWidth - BoxFrame - Safety(4)
		p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		label := theme.MetricLabelStyle.Render(fmt.Sprintf("%s: %.1f°C", sensorName, t.Temperature))
		content.WriteString(theme.BoxStyle.Width(l.BoxContentWidth(l.UsableWidth())).Render(label+"\n"+p.ViewAs(t.Temperature/100)) + "\n")
	}

	m.TempViewport.SetContent(content.String())
}

func (m Model) tempView() string {
	return m.TempViewport.View()
}

func (m Model) getTempProg(key string) progress.Model {
	p, ok := m.TempProgs[key]
	if !ok {
		p = progress.New(progress.WithDefaultGradient())
		p.FullColor = string(theme.AccentColor)
		// Initial width set in Update, but good to have fallback with safety margin
		l := layout.New(m.Width, m.Height)
		p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.TempProgs[key] = p
	}
	return p
}
