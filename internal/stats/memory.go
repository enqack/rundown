package stats

import "github.com/shirou/gopsutil/v3/mem"

// collectMemory gathers RAM and swap memory statistics
func collectMemory(s *SystemStats) {
	// Virtual Memory (RAM)
	v, err := mem.VirtualMemory()
	if err == nil {
		s.MemoryUsage = v.UsedPercent
		s.TotalMemory = v.Total
		s.UsedMemory = v.Used
	}

	// Swap Memory
	swap, err := mem.SwapMemory()
	if err == nil {
		s.SwapTotal = swap.Total
		s.SwapUsed = swap.Used
		s.SwapFree = swap.Free
		s.SwapPercent = swap.UsedPercent
	}
}
