package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enqack/rundown/internal/layout"
	"github.com/enqack/rundown/internal/stats"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.Tab = (m.Tab + 1) % 7 // Cycle through 7 tabs (0-6)
		case "shift+tab":
			m.Tab = (m.Tab + 6) % 7 // Go backward (equivalent to -1 with wrapping for 7 tabs)
		case "1":
			m.Tab = TabOverview
		case "2":
			m.Tab = TabCPU
		case "3":
			m.Tab = TabMem
		case "4":
			m.Tab = TabDisk
		case "5":
			m.Tab = TabNet
		case "6":
			m.Tab = TabTemp
		case "7":
			m.Tab = TabProc
		case "c":
			m.SortBy = "cpu"
		case "m":
			m.SortBy = "mem"
		case "p":
			m.SortBy = "pid"
		case "n":
			m.SortBy = "name"
		case "u":
			m.SortBy = "user"
		case "t":
			m.SortBy = "time"
		case "v":
			m.SortBy = "virt"
		case "f":
			m.SortBy = "foreign"
		case "up", "k":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.ScrollUp(1)
			case TabCPU:
				m.CPUViewport.ScrollUp(1)
			case TabMem:
				m.MemViewport.ScrollUp(1)
			case TabDisk:
				m.DiskViewport.ScrollUp(1)
			case TabNet:
				m.NetViewport.ScrollUp(1)
			case TabTemp:
				m.TempViewport.ScrollUp(1)
			case TabProc:
				m.ProcViewport.ScrollUp(1)
			}
		case "down", "j":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.ScrollDown(1)
			case TabCPU:
				m.CPUViewport.ScrollDown(1)
			case TabMem:
				m.MemViewport.ScrollDown(1)
			case TabDisk:
				m.DiskViewport.ScrollDown(1)
			case TabNet:
				m.NetViewport.ScrollDown(1)
			case TabTemp:
				m.TempViewport.ScrollDown(1)
			case TabProc:
				m.ProcViewport.ScrollDown(1)
			}
		case "pgup":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.PageUp()
			case TabCPU:
				m.CPUViewport.PageUp()
			case TabMem:
				m.MemViewport.PageUp()
			case TabDisk:
				m.DiskViewport.PageUp()
			case TabNet:
				m.NetViewport.PageUp()
			case TabTemp:
				m.TempViewport.PageUp()
			case TabProc:
				m.ProcViewport.PageUp()
			}
		case "pgdn":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.PageDown()
			case TabCPU:
				m.CPUViewport.PageDown()
			case TabMem:
				m.MemViewport.PageDown()
			case TabDisk:
				m.DiskViewport.PageDown()
			case TabNet:
				m.NetViewport.PageDown()
			case TabTemp:
				m.TempViewport.PageDown()
			case TabProc:
				m.ProcViewport.PageDown()
			}
		case "+", "=":
			// Increase update interval (max 10 seconds)
			if m.UpdateInterval < 10*time.Second {
				m.UpdateInterval += 500 * time.Millisecond
			}
		case "-", "_":
			// Decrease update interval (min 250ms)
			if m.UpdateInterval > 250*time.Millisecond {
				m.UpdateInterval -= 250 * time.Millisecond
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		// Unified Layout Calculator
		l := layout.New(m.Width, m.Height)

		// Global Progress Bars
		m.CPUProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.MemProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.DiskProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.SwapProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.NetUpProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
		m.NetDnProg.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))

		// Disk Progs
		for k := range m.DiskProgs {
			p := m.DiskProgs[k]
			p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
			m.DiskProgs[k] = p
		}

		// Temp Progs
		for k := range m.TempProgs {
			p := m.TempProgs[k]
			p.Width = l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))
			m.TempProgs[k] = p
		}

		// Viewport Resizing
		m.CPUViewport.Width = l.UsableWidth()
		m.CPUViewport.Height = l.UsableHeight()
		m.MemViewport.Width = l.UsableWidth()
		m.MemViewport.Height = l.UsableHeight()
		m.DiskViewport.Width = l.UsableWidth()
		m.DiskViewport.Height = l.UsableHeight()
		m.NetViewport.Width = l.UsableWidth()
		m.NetViewport.Height = l.UsableHeight()
		m.OverviewViewport.Width = l.UsableWidth()
		m.OverviewViewport.Height = l.UsableHeight()
		m.TempViewport.Width = l.UsableWidth()
		m.TempViewport.Height = l.UsableHeight()
		m.ProcViewport.Width = l.UsableWidth()
		m.ProcViewport.Height = l.UsableHeight()

		// Split View (Net/Cores)
		leftBoxW, _ := l.SplitTwoColumns(l.UsableWidth(), 2)
		splitBoxContentW := l.BoxContentWidth(leftBoxW)
		splitGraphSafeW := l.GraphWidth(splitBoxContentW)
		fullGraphSafeW := l.GraphWidth(l.BoxContentWidth(l.UsableWidth()))

		for k := range m.IfaceUpProgs {
			p := m.IfaceUpProgs[k]
			p.Width = splitGraphSafeW
			if m.Width <= 120 {
				p.Width = fullGraphSafeW
			}
			m.IfaceUpProgs[k] = p
		}
		for k := range m.IfaceDnProgs {
			p := m.IfaceDnProgs[k]
			p.Width = splitGraphSafeW
			if m.Width <= 120 {
				p.Width = fullGraphSafeW
			}
			m.IfaceDnProgs[k] = p
		}
		for i := range m.CPUCoresProgs {
			m.CPUCoresProgs[i].Width = splitGraphSafeW
			if m.Width <= 120 {
				m.CPUCoresProgs[i].Width = fullGraphSafeW
			}
		}

		m.updateCpuViewport()
		m.updateOverviewViewport()
		m.updateMemViewport()
		m.updateDiskViewport()
		m.updateNetViewport()
		m.updateTempViewport()
		m.updateProcViewport()

	case tickMsg:
		m.Stats = msg.Stats
		m.LastSync = time.Now()

		if m.Loading {
			m.LoadingMsg = msg.Step

			// Check what data we actually have to determine progress
			stepsCompleted := 0
			totalSteps := 6.0 // Increased to 6 for processes

			if len(m.Stats.CPUCores) > 0 {
				stepsCompleted++
			}
			if m.Stats.TotalMemory > 0 {
				stepsCompleted++
			}
			if len(m.Stats.Disks) > 0 {
				stepsCompleted++
			}
			if len(m.Stats.TopCPU) > 0 {
				stepsCompleted++
			}
			if m.Stats.NetSent > 0 {
				stepsCompleted++
			}
			if len(m.Stats.Processes) > 0 { // Check for processes
				stepsCompleted++
			}

			m.LoadVal = float64(stepsCompleted) / totalSteps

			if m.LoadVal >= 1.0 {
				m.Loading = false
				m.updateOverviewViewport()
				m.updateCpuViewport()
				m.updateMemViewport()
				m.updateDiskViewport()
				m.updateNetViewport()
				m.updateTempViewport()
				m.updateProcViewport()
			} else {
				// Determine next message based on what's missing
				var nextStep string
				if len(m.Stats.CPUCores) == 0 {
					nextStep = "Scanning CPU."
				} else if m.Stats.TotalMemory == 0 {
					nextStep = "Reading Memory."
				} else if len(m.Stats.Disks) == 0 {
					nextStep = "Mounting Disks."
				} else if len(m.Stats.TopCPU) == 0 {
					nextStep = "Listing Processes."
				} else if m.Stats.NetSent == 0 {
					nextStep = "Detecting Network."
				} else if len(m.Stats.Processes) == 0 { // New check for processes
					nextStep = "Gathering Process Info."
				} else {
					nextStep = "Finalizing."
				}
				return m, tickMsgCmd(nextStep)
			}
		}

		if msg.Err == nil {
			m.Stats = msg.Stats
			m.LastSync = time.Now()
			// Update all viewports to ensure data is ready when switching tabs
			m.updateOverviewViewport()
			m.updateCpuViewport()
			m.updateMemViewport()
			m.updateDiskViewport()
			m.updateNetViewport()
			m.updateTempViewport()
			m.updateProcViewport()
		}
		// Use configurable interval
		return m, tea.Tick(m.UpdateInterval, func(t time.Time) tea.Msg {
			s, err := stats.GetStats()
			return tickMsg{Stats: s, Err: err}
		})

	case progress.FrameMsg:
		progressModel, cmd := m.CPUProg.Update(msg)
		m.CPUProg = progressModel.(progress.Model)
		return m, cmd
	}

	var v *viewport.Model
	var cmd tea.Cmd
	switch m.Tab {
	case TabOverview:
		v = &m.OverviewViewport
	case TabCPU:
		v = &m.CPUViewport
	case TabMem:
		v = &m.MemViewport
	case TabDisk:
		v = &m.DiskViewport
	case TabNet:
		v = &m.NetViewport
	case TabTemp:
		v = &m.TempViewport
	case TabProc:
		v = &m.ProcViewport
	}

	if v != nil {
		*v, cmd = v.Update(msg)
	}

	return m, cmd
}
