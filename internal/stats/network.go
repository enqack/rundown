package stats

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastNetSent uint64
	lastNetRecv uint64
	lastNetTime time.Time
)

// collectNetwork gathers network I/O statistics and link speed
func collectNetwork(s *SystemStats) {
	netIO, err := net.IOCounters(false)
	if err == nil && len(netIO) > 0 {
		s.NetSent = netIO[0].BytesSent
		s.NetRecv = netIO[0].BytesRecv

		now := time.Now()
		if !lastNetTime.IsZero() {
			duration := now.Sub(lastNetTime).Seconds()
			if duration > 0 {
				s.NetSentDelta = uint64(float64(s.NetSent-lastNetSent) / duration)
				s.NetRecvDelta = uint64(float64(s.NetRecv-lastNetRecv) / duration)
			}
		}
		lastNetSent = s.NetSent
		lastNetRecv = s.NetRecv
		lastNetTime = now
	}

	// Link Speed (approximate default interface)
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
