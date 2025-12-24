package stats

import (
	"fmt"
	"time"
)

// GetStats collects all system statistics by orchestrating collection from specialized modules
// fullProcessDetails controls whether to collect detailed process info (expensive) or just lightweight stats
func GetStats(fullProcessDetails bool) (SystemStats, error) {
	var s SystemStats

	// Collect from specialized modules
	collectCPU(&s)
	collectMemory(&s)
	collectHostInfo(&s)
	collectDisks(&s)
	collectNetwork(&s)

	// Process collection returns process list needed for connections
	procs := collectProcesses(&s, fullProcessDetails)
	collectConnections(&s, procs)

	// Temperature collection
	collectTemperatures(&s)

	return s, nil
}

// FormatBytes converts bytes to human-readable format
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatDuration converts seconds to human-readable duration
func FormatDuration(seconds uint64) string {
	d := time.Duration(seconds) * time.Second
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
