package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
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
			m.Tab = (m.Tab + 1) % 6
		case "shift+tab":
			m.Tab = (m.Tab + 5) % 6 // Go backward (equivalent to -1 with wrapping)
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
			if m.Tab == TabTemp {
				m.TempViewport.LineUp(1)
			} else if m.Tab == TabCPU {
				m.CPUViewport.LineUp(1)
			}
		case "down", "j":
			if m.Tab == TabTemp {
				m.TempViewport.LineDown(1)
			} else if m.Tab == TabCPU {
				m.CPUViewport.LineDown(1)
			}
		case "pgup":
			if m.Tab == TabTemp {
				m.TempViewport.PageUp()
			} else if m.Tab == TabCPU {
				m.CPUViewport.PageUp()
			}
		case "pgdn":
			if m.Tab == TabTemp {
				m.TempViewport.PageDown()
			} else if m.Tab == TabCPU {
				m.CPUViewport.PageDown()
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

		m.CPUViewport.Width = l.UsableWidth()
		m.CPUViewport.Height = l.UsableHeight()

		// Full Box sizing
		fullBoxContentW := l.BoxContentWidth(l.UsableWidth())
		graphSafeW := l.GraphWidth(fullBoxContentW)

		m.MemProg.Width = graphSafeW
		m.SwapProg.Width = graphSafeW
		m.DiskProg.Width = graphSafeW
		m.LoadProg.Width = graphSafeW
		m.CPUProg.Width = graphSafeW

		// Disk Progs
		for k := range m.DiskProgs {
			p := m.DiskProgs[k]
			p.Width = graphSafeW
			m.DiskProgs[k] = p
		}

		// Temp Progs
		for k := range m.TempProgs {
			p := m.TempProgs[k]
			p.Width = graphSafeW
			m.TempProgs[k] = p
		}

		// Split View (Net/Cores)
		leftBoxW, _ := l.SplitTwoColumns(l.UsableWidth(), 0)

		splitBoxContentW := l.BoxContentWidth(leftBoxW)
		splitGraphSafeW := l.GraphWidth(splitBoxContentW)

		m.NetUpProg.Width = splitGraphSafeW
		if m.Width <= 120 {
			m.NetUpProg.Width = graphSafeW
		}
		m.NetDnProg.Width = m.NetUpProg.Width

		// CPU Cores Progs
		for i := range m.CPUCoresProgs {
			m.CPUCoresProgs[i].Width = splitGraphSafeW
			if m.Width <= 120 {
				m.CPUCoresProgs[i].Width = graphSafeW
			}
		}

		m.TempViewport.Width = l.UsableWidth()
		m.TempViewport.Height = l.UsableHeight()
		m.updateCpuViewport()
		m.updateTempViewport()

	case tickMsg:
		m.Stats = msg.Stats
		m.LastSync = time.Now()

		if m.Loading {
			m.LoadingMsg = msg.Step

			// Check what data we actually have to determine progress
			stepsCompleted := 0
			totalSteps := 5.0

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

			m.LoadVal = float64(stepsCompleted) / totalSteps

			if m.LoadVal >= 1.0 {
				m.Loading = false
				m.updateCpuViewport()
				m.updateTempViewport()
			} else {
				// Determine next message based on what's missing
				var nextStep string
				if len(m.Stats.CPUCores) == 0 {
					nextStep = "Scanning CPU..."
				} else if m.Stats.TotalMemory == 0 {
					nextStep = "Reading Memory..."
				} else if len(m.Stats.Disks) == 0 {
					nextStep = "Mounting Disks..."
				} else if len(m.Stats.TopCPU) == 0 {
					nextStep = "Listing Processes..."
				} else if m.Stats.NetSent == 0 {
					nextStep = "Detecting Network..."
				} else {
					nextStep = "Finalizing..."
				}
				return m, tickMsgCmd(nextStep)
			}
		}

		if msg.Err == nil {
			m.Stats = msg.Stats
			m.LastSync = time.Now()
			if m.Tab == TabTemp {
				m.updateTempViewport()
			} else if m.Tab == TabCPU {
				m.updateCpuViewport()
			}
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

	var cmd tea.Cmd
	if m.Tab == TabTemp {
		m.TempViewport, cmd = m.TempViewport.Update(msg)
	} else if m.Tab == TabCPU {
		m.CPUViewport, cmd = m.CPUViewport.Update(msg)
	}
	return m, cmd
}
