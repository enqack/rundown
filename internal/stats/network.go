package stats

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastIfaceStats map[string]net.IOCountersStat
	lastNetTime    time.Time
)

func init() {
	lastIfaceStats = make(map[string]net.IOCountersStat)
}

// collectNetwork gathers network I/O statistics and link speed
func collectNetwork(s *SystemStats) {
	// Collect per-interface counters
	ifaces, err := net.IOCounters(true)
	if err != nil || len(ifaces) == 0 {
		s.LinkSpeed = detectLinkSpeed()
		return
	}

	now := time.Now()
	var totalSent, totalRecv uint64
	var duration float64
	if !lastNetTime.IsZero() {
		duration = now.Sub(lastNetTime).Seconds()
	}

	for _, iface := range ifaces {
		// Aggregate global totals
		totalSent += iface.BytesSent
		totalRecv += iface.BytesRecv

		stat := InterfaceStat{
			Name: iface.Name,
			Sent: iface.BytesSent,
			Recv: iface.BytesRecv,
		}

		// Calculate deltas if we have previous state
		if duration > 0 {
			if last, ok := lastIfaceStats[iface.Name]; ok {
				stat.SentDelta = uint64(float64(iface.BytesSent-last.BytesSent) / duration)
				stat.RecvDelta = uint64(float64(iface.BytesRecv-last.BytesRecv) / duration)
			}
		}

		s.Interfaces = append(s.Interfaces, stat)
		lastIfaceStats[iface.Name] = iface
	}

	// Update global aggregate metrics for HUD consistency
	s.NetSent = totalSent
	s.NetRecv = totalRecv

	if duration > 0 {
		var aggSentDelta, aggRecvDelta uint64
		for _, iface := range s.Interfaces {
			aggSentDelta += iface.SentDelta
			aggRecvDelta += iface.RecvDelta
		}
		s.NetSentDelta = aggSentDelta
		s.NetRecvDelta = aggRecvDelta
	}

	lastNetTime = now
	s.LinkSpeed = detectLinkSpeed()
}

// detectLinkSpeed attempts to detect network link speed
func detectLinkSpeed() uint64 {
	iface := getDefaultInterface()
	if iface == "" {
		return 1000000000 // 1 Gbps default
	}

	// Try to read link speed from sysctl
	// On Linux: /sys/class/net/<iface>/speed (in Mbps)
	// This requires root or appropriate permissions
	// Fallback to 1 Gbps
	return 1000000000
}

// getDefaultInterface attempts to find the default network interface
func getDefaultInterface() string {
	// Simplified: return first non-loopback interface
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		// Skip loopback
		if iface.Name == "lo" || iface.Name == "lo0" {
			continue
		}
		// Return first viable interface
		if len(iface.Addrs) > 0 {
			return iface.Name
		}
	}
	return ""
}
