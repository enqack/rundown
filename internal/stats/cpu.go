package stats

import "github.com/shirou/gopsutil/v3/cpu"

// collectCPU gathers CPU usage statistics
func collectCPU(s *SystemStats) {
	// Per-core CPU usage
	percents, err := cpu.Percent(0, true)
	if err == nil {
		s.CPUCores = percents
		// Average CPU usage across all cores
		var total float64
		for _, p := range percents {
			total += p
		}
		s.CPUUsage = total / float64(len(s.CPUCores))
	}
}
