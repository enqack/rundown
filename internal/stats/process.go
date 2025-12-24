package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// collectProcesses gathers process statistics
// fullDetails controls whether to collect detailed info for all top processes (Process tab)
// or just lightweight stats for counts and top 10 (other tabs)
func collectProcesses(s *SystemStats, fullDetails bool) []*process.Process {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	// Reset counters
	s.ProcTotal = 0
	s.ProcRunning = 0
	s.ProcSleeping = 0
	s.ProcStopped = 0
	s.ProcZombie = 0

	// First pass: lightweight counting only (no expensive fields)
	type lightProc struct {
		p      *process.Process
		cpuPct float64
	}
	var candidates []lightProc

	for _, p := range procs {
		s.ProcTotal++

		// Get status for counting (relatively cheap)
		statusStrSlice, _ := p.Status()
		state := " "
		if len(statusStrSlice) > 0 && len(statusStrSlice[0]) > 0 {
			state = statusStrSlice[0][0:1]
		}

		// Count states
		switch state {
		case "R":
			s.ProcRunning++
		case "S", "I":
			s.ProcSleeping++
		case "T":
			s.ProcStopped++
		case "Z":
			s.ProcZombie++
		}

		// Get CPU for sorting (expensive but necessary)
		cpuPct, err := p.CPUPercent()
		if err == nil {
			// For full details (Process tab), include ALL processes
			// For lightweight mode, only include processes with active CPU
			if fullDetails || cpuPct > 0.01 {
				candidates = append(candidates, lightProc{p: p, cpuPct: cpuPct})
			}
		}
	}

	// Sort by CPU and take top 10 for lightweight mode, or all for Process tab
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].cpuPct > candidates[j].cpuPct
	})

	// Limit based on detail level
	if !fullDetails {
		// Lightweight mode: only top 10 for display on other tabs
		limit := 10
		if len(candidates) < limit {
			limit = len(candidates)
		}
		candidates = candidates[:limit]
	}
	// else: Process tab gets ALL processes (no limit)

	// Second pass: detailed info only for top processes
	for _, lp := range candidates {
		p := lp.p
		name, _ := p.Name()
		user, _ := p.Username()
		statusStrSlice, _ := p.Status()
		state := " "
		if len(statusStrSlice) > 0 && len(statusStrSlice[0]) > 0 {
			state = statusStrSlice[0][0:1]
		}

		nice, _ := p.Nice()

		memInfo, _ := p.MemoryInfo()
		var vms, rss, shared uint64
		if memInfo != nil {
			vms = memInfo.VMS
			rss = memInfo.RSS
		}

		memPct, _ := p.MemoryPercent()

		createTime, _ := p.CreateTime()
		runningTime := ""
		if createTime > 0 {
			dt := time.Since(time.Unix(createTime/1000, 0))
			runningTime = fmt.Sprintf("%02d:%02d.%02d", int(dt.Minutes()), int(dt.Seconds())%60, (dt.Milliseconds()/10)%100)
		}

		cmdline, _ := p.Cmdline()
		if cmdline == "" {
			cmdline = name
		}

		pinfo := ProcessInfo{
			PID:      p.Pid,
			Name:     name,
			User:     user,
			State:    state,
			Priority: 20 + nice,
			Nice:     nice,
			Virtual:  vms,
			Resident: rss,
			Shared:   shared,
			CPU:      lp.cpuPct,
			Memory:   float64(memPct),
			Time:     runningTime,
			Cmdline:  cmdline,
		}

		s.Processes = append(s.Processes, pinfo)
	}

	// Sort and limit for top views (use full process list)
	s.TopCPU = sortAndLimit(s.Processes, "cpu", 10)
	s.TopMem = sortAndLimit(s.Processes, "mem", 10)

	return procs
}

// collectConnections gathers network connections
func collectConnections(s *SystemStats, procs []*process.Process) {
	conns, err := net.Connections("all")
	if err != nil {
		return
	}

	pnames := make(map[int32]string)
	for _, p := range procs {
		name, _ := p.Name()
		pnames[p.Pid] = name
	}

	for _, c := range conns {
		if c.Status == "LISTEN" || c.Status == "ESTABLISHED" || c.Status == "TIME_WAIT" {
			proto := "TCP"
			if c.Type == 2 {
				proto = "UDP"
			}

			dir := "OUT"
			if c.Status == "LISTEN" {
				dir = "LISTEN"
			} else if c.Laddr.Port < 32768 {
				dir = "IN"
			}

			s.Connections = append(s.Connections, ConnectionInfo{
				Proto:     proto,
				Local:     fmt.Sprintf("%s:%d", c.Laddr.IP, c.Laddr.Port),
				LocalPort: c.Laddr.Port,
				Remote:    fmt.Sprintf("%s:%d", c.Raddr.IP, c.Raddr.Port),
				State:     c.Status,
				Direction: dir,
				PID:       c.Pid,
				PName:     pnames[c.Pid],
			})
		}
	}

	// Sort by State (LISTEN first), then Local Addr
	sort.SliceStable(s.Connections, func(i, j int) bool {
		if s.Connections[i].State != s.Connections[j].State {
			return s.Connections[i].State < s.Connections[j].State
		}
		return s.Connections[i].Local < s.Connections[j].Local
	})

	// Limit to prevent overflow
	if len(s.Connections) > 15 {
		s.Connections = s.Connections[:15]
	}
}

// sortAndLimit sorts processes by CPU or memory and limits results
func sortAndLimit(procs []ProcessInfo, sortBy string, limit int) []ProcessInfo {
	sort.SliceStable(procs, func(i, j int) bool {
		if sortBy == "cpu" {
			return procs[i].CPU > procs[j].CPU
		}
		return procs[i].Memory > procs[j].Memory
	})

	if len(procs) > limit {
		return procs[:limit]
	}
	return procs
}
