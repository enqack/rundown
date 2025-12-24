package stats

import (
	"net"

	"github.com/shirou/gopsutil/v3/host"
)

// collectHostInfo gathers host information including IP address
func collectHostInfo(s *SystemStats) {
	h, err := host.Info()
	if err == nil {
		s.Uptime = h.Uptime
		s.HostName = h.Hostname
	}

	// Get primary IP address
	s.IPAddress = getIPAddress()
}

// getIPAddress attempts to get the primary IP address
func getIPAddress() string {
	// Try common methods to get local IP
	// Method 1: Dial out to get preferred outbound IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	// Method 2: Fallback to first non-local address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "N/A"
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return "N/A"
}
