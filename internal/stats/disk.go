package stats

import (
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// collectDisks gathers disk mount point and usage statistics
func collectDisks(s *SystemStats) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}

	for _, p := range partitions {
		// Skip virtual/pseudo filesystems and common automation mounts
		if isVirtualFS(p.Fstype) || strings.HasPrefix(p.Mountpoint, "/var/lib/docker") || strings.HasPrefix(p.Mountpoint, "/run") {
			continue
		}

		usage, err := disk.Usage(p.Mountpoint)
		if err == nil {
			s.Disks = append(s.Disks, DiskInfo{
				MountPoint: p.Mountpoint,
				Used:       usage.Used,
				Total:      usage.Total,
				Usage:      usage.UsedPercent,
			})
		}
	}

	// Sort disks by mount point
	sort.SliceStable(s.Disks, func(i, j int) bool {
		return s.Disks[i].MountPoint < s.Disks[j].MountPoint
	})

	// Calculate aggregates
	for _, d := range s.Disks {
		s.DiskTotal += d.Total
		s.DiskUsed += d.Used
	}
}

// isVirtualFS checks if a filesystem type is virtual/pseudo
func isVirtualFS(fs string) bool {
	virtualFS := map[string]bool{
		"tmpfs":    true,
		"devtmpfs": true,
		"proc":     true,
		"sysfs":    true,
		"devpts":   true,
		"debugfs":  true,
		"fuse":     true,
		"cgroup":   true,
		"mqueue":   true,
		"configfs": true,
		"autofs":   true,
		"pstore":   true,
		"squashfs": true,
	}
	return virtualFS[fs]
}
