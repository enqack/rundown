package ui

import (
	"fmt"
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
			// Update the new tab's viewport when switching
			switch m.Tab {
			case TabOverview:
				m.updateOverviewViewport()
			case TabCPU:
				m.updateCpuViewport()
			case TabMem:
				m.updateMemViewport()
			case TabDisk:
				m.updateDiskViewport()
			case TabNet:
				m.updateNetViewport()
			case TabTemp:
				m.updateTempViewport()
			case TabProc:
				m.updateProcViewport()
			}
		case "shift+tab":
			m.Tab = (m.Tab + 6) % 7 // Go backward (equivalent to -1 with wrapping for 7 tabs)
		case "?":
			m.Tab = TabHelp
			m.updateHelpViewport()
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
			case TabHelp:
				m.HelpViewport.ScrollUp(1)
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
			case TabHelp:
				m.HelpViewport.ScrollDown(1)
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
			case TabHelp:
				m.HelpViewport.PageUp()
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
			case TabHelp:
				m.HelpViewport.PageDown()
			}
		case "home":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.GotoTop()
			case TabCPU:
				m.CPUViewport.GotoTop()
			case TabMem:
				m.MemViewport.GotoTop()
			case TabDisk:
				m.DiskViewport.GotoTop()
			case TabNet:
				m.NetViewport.GotoTop()
			case TabTemp:
				m.TempViewport.GotoTop()
			case TabProc:
				m.ProcViewport.GotoTop()
			case TabHelp:
				m.HelpViewport.GotoTop()
			}
		case "end":
			switch m.Tab {
			case TabOverview:
				m.OverviewViewport.GotoBottom()
			case TabCPU:
				m.CPUViewport.GotoBottom()
			case TabMem:
				m.MemViewport.GotoBottom()
			case TabDisk:
				m.DiskViewport.GotoBottom()
			case TabNet:
				m.NetViewport.GotoBottom()
			case TabTemp:
				m.TempViewport.GotoBottom()
			case TabProc:
				m.ProcViewport.GotoBottom()
			case TabHelp:
				m.HelpViewport.GotoBottom()
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
		m.HelpViewport.Width = l.UsableWidth()
		m.HelpViewport.Height = l.UsableHeight()

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
		m.updateHelpViewport()

	case tickMsg:
		if msg.Err != nil {
			if m.Loading {
				m.LoadingMsg = fmt.Sprintf("Error: %v", msg.Err)
			}
		} else {
			m.Stats = msg.Stats
			m.LastSync = time.Now()
		}

		if m.Loading {
			// If we have an error, we keep showing it until success or timeout (implied retry via Tick)
			if msg.Err != nil {
				return m, tea.Tick(m.UpdateInterval, func(t time.Time) tea.Msg {
					s, err := stats.GetStats(false)
					return tickMsg{Stats: s, Err: err}
				})
			}

			m.LoadingMsg = msg.Step
			m.AnimFrame++ // Increment animation frame for road effect

			// Check strictly for critical data: CPU and Memory
			stepsCompleted := 0.0
			totalSteps := 2.0

			if len(m.Stats.CPUCores) > 0 {
				stepsCompleted++
			}
			if m.Stats.TotalMemory > 0 {
				stepsCompleted++
			}

			m.LoadVal = stepsCompleted / totalSteps

			if m.LoadVal >= 1.0 {
				// Ensure minimum splash display time of 1 second
				elapsed := time.Since(m.LoadingStartTime)
				minDisplayTime := 1 * time.Second

				if elapsed < minDisplayTime {
					// Wait for remaining time before transitioning
					m.LoadingMsg = "Finished"
					remainingTime := minDisplayTime - elapsed
					return m, tea.Tick(remainingTime, func(t time.Time) tea.Msg {
						return tickMsg{Stats: m.Stats, Step: "Finished"}
					})
				}

				// Minimum time elapsed, transition immediately
				m.Loading = false
				// Update all viewports on initial load
				m.updateOverviewViewport()
				m.updateCpuViewport()
				m.updateMemViewport()
				m.updateDiskViewport()
				m.updateNetViewport()
				m.updateTempViewport()
				m.updateProcViewport()
				m.updateHelpViewport()
			} else {
				// Determine next message based on what's missing
				var nextStep string
				if len(m.Stats.CPUCores) == 0 {
					nextStep = "Scanning CPU..."
				} else if m.Stats.TotalMemory == 0 {
					nextStep = "Reading Memory..."
				} else {
					nextStep = "Finalizing..."
				}
				return m, tickMsgCmd(nextStep)
			}
		}

		if msg.Err == nil {
			// Only update the active tab's viewport to reduce CPU usage
			switch m.Tab {
			case TabOverview:
				m.updateOverviewViewport()
			case TabCPU:
				m.updateCpuViewport()
			case TabMem:
				m.updateMemViewport()
			case TabDisk:
				m.updateDiskViewport()
			case TabNet:
				m.updateNetViewport()
			case TabTemp:
				m.updateTempViewport()
			case TabProc:
				m.updateProcViewport()
			}
		}
		// Use configurable interval
		return m, tea.Tick(m.UpdateInterval, func(t time.Time) tea.Msg {
			// Only collect full process details if on Process tab
			fullDetails := m.Tab == TabProc
			s, err := stats.GetStats(fullDetails)
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
	case TabHelp:
		v = &m.HelpViewport
	}

	if v != nil {
		*v, cmd = v.Update(msg)
	}

	return m, cmd
}
