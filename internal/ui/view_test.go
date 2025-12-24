package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestGetAdaptiveScale(t *testing.T) {
	// Link Speed: 100 Mbps (100,000,000)
	linkSpeed := uint64(100_000_000)

	tests := []struct {
		name     string
		usage    uint64
		expected float64
	}{
		{"Tiny Usage (1kbps)", 1000, float64(linkSpeed) * 0.001},      // Should snap to 0.1%
		{"Low Usage (2Mbps)", 2_000_000, float64(linkSpeed) * 0.05},   // 2% usage -> snaps to 5%
		{"Mid Usage (40Mbps)", 40_000_000, float64(linkSpeed) * 0.50}, // 40% usage -> snaps to 50%
		{"High Usage (90Mbps)", 90_000_000, float64(linkSpeed) * 1.0}, // 90% usage -> snaps to 100%
		{"Burst Usage (150Mbps)", 150_000_000, 150_000_000.0},         // >100% -> actual usage
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAdaptiveScale(tt.usage, linkSpeed)
			if got != tt.expected {
				t.Errorf("getAdaptiveScale(%d, %d) = %f; want %f", tt.usage, linkSpeed, got, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	s := "Requesting"

	// Test no truncate needed
	if got := truncate(s, 20); got != s {
		t.Errorf("Expected %s, got %s", s, got)
	}

	// Test truncate
	limit := 5
	got := truncate(s, limit)
	if lipgloss.Width(got) > limit {
		t.Errorf("Truncated string width %d exceeds limit %d", lipgloss.Width(got), limit)
	}
}
