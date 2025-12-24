package stats

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/mem"
)

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

	// Detailed Swap Devices
	devices, err := mem.SwapDevices()
	if err == nil {
		s.SwapDevices = make([]SwapDevice, len(devices))
		for i, d := range devices {
			finalName := d.Name

			// 1. Resolve Symlinks (e.g. /dev/mapper/swap -> /dev/dm-0)
			resolved, err := filepath.EvalSymlinks(d.Name)
			if err == nil {
				finalName = resolved
			}

			// 2. Resolve DM to Physical (e.g. /dev/dm-0 -> /dev/nvme0n1p4)
			if strings.HasPrefix(finalName, "/dev/dm-") {
				dmName := filepath.Base(finalName)
				slavesDir := fmt.Sprintf("/sys/class/block/%s/slaves", dmName)
				entries, err := os.ReadDir(slavesDir)
				if err == nil && len(entries) > 0 {
					// Use the first slave found (usually only one for linear swap)
					finalName = "/dev/" + entries[0].Name()
				}
			}

			s.SwapDevices[i] = SwapDevice{
				Name:       finalName,
				UsedBytes:  d.UsedBytes,
				FreeBytes:  d.FreeBytes,
				TotalBytes: d.UsedBytes + d.FreeBytes,
			}
		}
	}
}
