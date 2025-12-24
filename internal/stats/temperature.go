package stats

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

// collectTemperatures gathers temperature sensor data
func collectTemperatures(s *SystemStats) {
	temps, _ := host.SensorsTemperatures()
	if len(temps) > 0 {
		for _, t := range temps {
			tinfo := TempInfo{
				SensorKey:   t.SensorKey,
				Temperature: t.Temperature,
			}
			s.Temperatures = append(s.Temperatures, tinfo)

			// Heuristic for CPU/GPU
			key := strings.ToLower(t.SensorKey)
			if s.CPUTemp == 0 && (strings.Contains(key, "coretemp") || strings.Contains(key, "cpu") || strings.Contains(key, "k10temp") || strings.Contains(key, "tctl") || strings.Contains(key, "tdie") || strings.Contains(key, "package")) {
				s.CPUTemp = t.Temperature
			}
			if strings.Contains(key, "amdgpu") || strings.Contains(key, "edge") || strings.Contains(key, "junction") || strings.Contains(key, "gpu") {
				if !strings.Contains(key, "nvidia") {
					s.GPUTemps = append(s.GPUTemps, tinfo)
				}
			}
		}
	}

	// Nvidia Fallback (nvidia-smi)
	nvTemps := getNvidiaTemps()
	if len(nvTemps) > 0 {
		// Remove any generic nvidia sensors from gopsutil
		newTemps := []TempInfo{}
		for _, t := range s.Temperatures {
			if !strings.Contains(strings.ToLower(t.SensorKey), "nvidia") {
				newTemps = append(newTemps, t)
			}
		}
		s.Temperatures = append(newTemps, nvTemps...)
		s.GPUTemps = append(s.GPUTemps, nvTemps...)
	}

	sort.SliceStable(s.Temperatures, func(i, j int) bool {
		return s.Temperatures[i].SensorKey < s.Temperatures[j].SensorKey
	})
}

// getNvidiaTemps attempts to get Nvidia GPU temperatures via nvidia-smi
func getNvidiaTemps() []TempInfo {
	var results []TempInfo
	out, err := exec.Command("nvidia-smi", "--query-gpu=temperature.gpu", "--format=csv,noheader,nounits").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for i, line := range lines {
			var temp float64
			fmt.Sscanf(line, "%f", &temp)
			results = append(results, TempInfo{
				SensorKey:   fmt.Sprintf("Nvidia GPU %d", i),
				Temperature: temp,
			})
		}
	}
	return results
}
