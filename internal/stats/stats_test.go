package stats

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero", 0, "0 B"},
		{"Bytes", 500, "500 B"},
		{"Kilobytes", 1024, "1.0 KiB"},
		{"Megabytes", 1024 * 1024, "1.0 MiB"},
		{"Gigabytes", 10 * 1024 * 1024 * 1024, "10.0 GiB"},
		{"Fractional", 1536, "1.5 KiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBytes(tt.input); got != tt.expected {
				t.Errorf("FormatBytes() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  uint64
		expected string
	}{
		{"LessThanMinute", 30, "0m"},
		{"OneMinute", 60, "1m"},
		{"HourAndMinute", 3661, "1h 1m"},
		{"DayHourMinute", 90061, "1d 1h 1m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.seconds); got != tt.expected {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}
