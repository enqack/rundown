package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/theme"
)

func (m *Model) updateTempViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("System Temperature Sensors") + "\n\n")

	l := layout.New(m.Width, m.Height)

	// CPU Package Temperature (0-100 scale)
	pCPU := m.getTempProg("CPU Package")
	s.WriteString(m.renderTacticalGauge(l.UsableWidth(), "CPU Package", pCPU, m.Stats.CPUTemp/110.0, m.Stats.CPUTemp, "°C") + "\n")

	// Individual Sensor View
	for _, t := range m.Stats.Temperatures {
		if t.SensorKey == "CPU Package" {
			continue
		}
		p := m.getTempProg(t.SensorKey)
		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), t.SensorKey, p, t.Temperature/110.0, t.Temperature, "°C") + "\n")
	}

	// GPU Temperatures
	for i, g := range m.Stats.GPUTemps {
		sensorName := fmt.Sprintf("GPU %d (%s)", i, g.SensorKey)
		p := m.getTempProg(sensorName)
		s.WriteString(m.renderTacticalGauge(l.UsableWidth(), sensorName, p, g.Temperature/110.0, g.Temperature, "°C") + "\n")
	}

	m.TempViewport.SetContent(s.String())
}

func (m Model) tempView() string {
	return m.TempViewport.View()
}

func (m Model) getTempProg(key string) progress.Model {
	p, ok := m.TempProgs[key]
	if !ok {
		p = theme.NewThemedProgress()
		// Initial width set in Update, but good to have fallback with safety margin
		l := layout.New(m.Width, m.Height)
		p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.TempProgs[key] = p
	}
	return p
}
