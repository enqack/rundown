package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) netView() string {
	return m.NetViewport.View()
}

func (m *Model) updateNetViewport() {
	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Network Statistics") + "\n\n")

	l := layout.New(m.Width, m.Height)

	// Global Aggregate Graphs
	scaleUp := getStepScale(m.Stats.NetSentDelta*8, m.Stats.LinkSpeed)
	scaleDn := getStepScale(m.Stats.NetRecvDelta*8, m.Stats.LinkSpeed)

	netUpTitle := fmt.Sprintf("Global Egress (Scale: %s/s) | Total: %s",
		stats.FormatBytes(uint64(float64(m.Stats.LinkSpeed/8)*scaleUp)), stats.FormatBytes(m.Stats.NetSent))
	netDnTitle := fmt.Sprintf("Global Ingress (Scale: %s/s) | Total: %s",
		stats.FormatBytes(uint64(float64(m.Stats.LinkSpeed/8)*scaleDn)), stats.FormatBytes(m.Stats.NetRecv))

	upRatio := float64(m.Stats.NetSentDelta*8) / (float64(m.Stats.LinkSpeed) * scaleUp)
	dnRatio := float64(m.Stats.NetRecvDelta*8) / (float64(m.Stats.LinkSpeed) * scaleDn)

	leftBoxW, rightBoxW := l.SplitTwoColumns(l.UsableWidth(), 2)

	if m.Width > 120 {
		upBox := m.renderTacticalGauge(leftBoxW, netUpTitle, m.NetUpProg, upRatio, upRatio*100, "%")
		dnBox := m.renderTacticalGauge(rightBoxW, netDnTitle, m.NetDnProg, dnRatio, dnRatio*100, "%")
		s.WriteString(m.renderTacticalRow(l.UsableWidth(), 2, upBox, dnBox) + "\n")
	} else {
		upBar := m.renderTacticalGauge(l.UsableWidth(), netUpTitle, m.NetUpProg, upRatio, upRatio*100, "%")
		dnBar := m.renderTacticalGauge(l.UsableWidth(), netDnTitle, m.NetDnProg, dnRatio, dnRatio*100, "%")
		s.WriteString(upBar + "\n" + dnBar + "\n")
	}

	// Per-Interface Graphs
	for _, iface := range m.Stats.Interfaces {
		if iface.Sent == 0 && iface.Recv == 0 && (iface.Name == "lo" || iface.Name == "lo0") {
			continue
		}

		if _, ok := m.IfaceUpProgs[iface.Name]; !ok {
			up := theme.NewThemedProgress()
			dn := theme.NewThemedProgress()
			m.IfaceUpProgs[iface.Name] = up
			m.IfaceDnProgs[iface.Name] = dn
		}

		upProg := m.IfaceUpProgs[iface.Name]
		dnProg := m.IfaceDnProgs[iface.Name]

		iScaleUp := getStepScale(iface.SentDelta*8, m.Stats.LinkSpeed)
		iScaleDn := getStepScale(iface.RecvDelta*8, m.Stats.LinkSpeed)

		iUpTitle := fmt.Sprintf("%s Egress: %s/s", iface.Name, stats.FormatBytes(iface.SentDelta))
		iDnTitle := fmt.Sprintf("%s Ingress: %s/s", iface.Name, stats.FormatBytes(iface.RecvDelta))

		upRatio := float64(iface.SentDelta*8) / (float64(m.Stats.LinkSpeed) * iScaleUp)
		dnRatio := float64(iface.RecvDelta*8) / (float64(m.Stats.LinkSpeed) * iScaleDn)

		if m.Width > 120 {
			upBox := m.renderTacticalGauge(leftBoxW, iUpTitle, upProg, upRatio, upRatio*100, "%")
			dnBox := m.renderTacticalGauge(rightBoxW, iDnTitle, dnProg, dnRatio, dnRatio*100, "%")
			s.WriteString(m.renderTacticalRow(l.UsableWidth(), 2, upBox, dnBox) + "\n")
		} else {
			upBar := m.renderTacticalGauge(l.UsableWidth(), iUpTitle, upProg, upRatio, upRatio*100, "%")
			dnBar := m.renderTacticalGauge(l.UsableWidth(), iDnTitle, dnProg, dnRatio, dnRatio*100, "%")
			s.WriteString(upBar + "\n" + dnBar + "\n")
		}

		m.IfaceUpProgs[iface.Name] = upProg
		m.IfaceDnProgs[iface.Name] = dnProg
	}

	s.WriteString("\n")

	// Active Connections Table (Bottom)
	headers := []string{"PROTO", "DIR", "STATE", "LOCAL", "FOREIGN", "PROGRAM"}
	usedW := 6 + 7 + 15
	paddingOverhead := 6 * 2
	rem := l.UsableWidth() - usedW - paddingOverhead
	if rem < 10 {
		rem = 10
	}
	w := rem / 3
	widths := []int{6, 7, 15, w, w, w}

	var rows [][]string
	for _, c := range m.Stats.Connections {
		rows = append(rows, []string{c.Proto, c.Direction, c.State, c.Local, c.Remote, fmt.Sprintf("%d/%s", c.PID, c.PName)})
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i][3] < rows[j][3] // Sort by local addr by default
	})

	s.WriteString(m.renderTacticalTable("Active Network Connections", headers, widths, rows))

	m.NetViewport.SetContent(s.String())
}
