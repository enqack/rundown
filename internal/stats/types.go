package stats

// DiskInfo represents information about a disk mount point
type DiskInfo struct {
	MountPoint string
	Used       uint64
	Total      uint64
	Usage      float64
}

// ProcessInfo represents information about a running process
type ProcessInfo struct {
	PID      int32
	Name     string
	User     string
	State    string
	Priority int32
	Nice     int32
	Virtual  uint64
	Resident uint64
	Shared   uint64
	CPU      float64
	Memory   float64
	Time     string
	Cmdline  string
}

// ConnectionInfo represents a network connection
type ConnectionInfo struct {
	Proto     string
	Local     string
	LocalPort uint32
	Remote    string
	State     string
	Direction string
	PID       int32
	PName     string
}

// TempInfo represents temperature sensor data
type TempInfo struct {
	SensorKey   string
	Temperature float64
}

// SystemStats contains all collected system statistics
type SystemStats struct {
	CPUUsage     float64
	CPUCores     []float64
	MemoryUsage  float64
	TotalMemory  uint64
	UsedMemory   uint64
	SwapTotal    uint64
	SwapUsed     uint64
	SwapFree     uint64
	SwapPercent  float64
	Uptime       uint64
	IPAddress    string
	HostName     string
	Disks        []DiskInfo
	DiskTotal    uint64
	DiskUsed     uint64
	NetSent      uint64
	NetRecv      uint64
	NetSentDelta uint64
	NetRecvDelta uint64
	LinkSpeed    uint64 // in bits per second
	TopCPU       []ProcessInfo
	TopMem       []ProcessInfo
	Connections  []ConnectionInfo
	Temperatures []TempInfo
	CPUTemp      float64
	GPUTemps     []TempInfo
}
