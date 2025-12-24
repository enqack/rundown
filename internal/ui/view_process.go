package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func (m Model) procView() string {
	return m.ProcViewport.View()
}

func (m *Model) updateProcViewport() {
	l := layout.New(m.Width, m.Height)

	var s strings.Builder
	s.WriteString(theme.TitleStyle.Render("Detailed Process Monitor") + "\n\n")

	// Process Summary (5 Columns)
	// [Total] [Running] [Sleeping] [Stopped] [Zombie]
	procWidths := l.SplitColumns(l.UsableWidth(), 5, 1) // 5 columns, 1 char gap

	if len(procWidths) == 5 {
		bTotal := m.renderTacticalBox(procWidths[0], fmt.Sprintf("Total: %d", m.Stats.ProcTotal))
		bRun := m.renderTacticalBox(procWidths[1], fmt.Sprintf("Running: %d", m.Stats.ProcRunning))
		bSleep := m.renderTacticalBox(procWidths[2], fmt.Sprintf("Sleeping: %d", m.Stats.ProcSleeping))
		bStop := m.renderTacticalBox(procWidths[3], fmt.Sprintf("Stopped: %d", m.Stats.ProcStopped))
		bZom := m.renderTacticalBox(procWidths[4], fmt.Sprintf("Zombie: %d", m.Stats.ProcZombie))

		s.WriteString(m.renderTacticalRow(l.UsableWidth(), 1, bTotal, bRun, bSleep, bStop, bZom) + "\n")
	} else {
		// Fallback
		s.WriteString(m.renderTacticalBox(l.UsableWidth(), fmt.Sprintf("Process Summary: %d", m.Stats.ProcTotal)) + "\n")
	}

	// Prepare Headers
	// Default htop columns: PID, USER, PRI, NI, VIRT, RES, SHR, S, CPU%, MEM%, TIME+, Command
	headers := []string{"PID", "USER", "PRI", "NI", "VIRT", "RES", "SHR", "S", "CPU%", "MEM%", "TIME+", "COMMAND"}

	// Column Widths (total UsableWidth - padding)
	// PID(7) USER(9) PRI(3) NI(3) VIRT(7) RES(7) SHR(7) S(1) CPU%(5) MEM%(5) TIME+(9)
	widths := []int{7, 9, 3, 3, 7, 7, 7, 1, 5, 5, 9, 0} // 0 means dynamic

	// Calculate dynamic width for COMMAND column
	fixedW := 0
	for _, w := range widths[:len(widths)-1] {
		fixedW += w
	}
	widths[len(widths)-1] = l.CalculateTableDynamicWidth(l.UsableWidth(), widths[:len(widths)-1], 2)

	// Sort Processes (Current SortBy or Default to CPU)
	procs := make([]stats.ProcessInfo, len(m.Stats.Processes))
	copy(procs, m.Stats.Processes)

	sort.SliceStable(procs, func(i, j int) bool {
		switch m.SortBy {
		case "mem":
			return procs[i].Memory > procs[j].Memory
		case "pid":
			return procs[i].PID < procs[j].PID
		case "name":
			return procs[i].Name < procs[j].Name
		case "user":
			return procs[i].User < procs[j].User
		case "time":
			return procs[i].Time > procs[j].Time
		default: // cpu
			return procs[i].CPU > procs[j].CPU
		}
	})

	// Prepare Row Data
	var data [][]string
	for _, p := range procs {
		virt := formatBytes(p.Virtual)
		res := formatBytes(p.Resident)
		shr := formatBytes(p.Shared)

		data = append(data, []string{
			fmt.Sprintf("%d", p.PID),
			p.User,
			fmt.Sprintf("%d", p.Priority),
			fmt.Sprintf("%d", p.Nice),
			virt,
			res,
			shr,
			p.State,
			fmt.Sprintf("%.1f", p.CPU),
			fmt.Sprintf("%.1f", p.Memory),
			p.Time,
			p.Cmdline,
		})
	}

	table := m.renderTacticalTable("Processes", headers, widths, data)
	s.WriteString(table)

	m.ProcViewport.Width = l.UsableWidth()
	m.ProcViewport.Height = l.UsableHeight()
	m.ProcViewport.SetContent(s.String())
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cI", float64(b)/float64(div), "KMGTPE"[exp])
}
