package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestValidateTheme(t *testing.T) {
	valid := []string{"base16", "cyberpunk", "monochrome", "phosphor"}
	for _, v := range valid {
		if !ValidateTheme(v) {
			t.Errorf("Theme %s should be valid", v)
		}
	}

	if ValidateTheme("invalid_theme_name") {
		t.Error("invalid_theme_name should be invalid")
	}
}

func TestInit(t *testing.T) {
	// Test default fallback
	Init("unknown")
	if current.Name != "base16" {
		t.Errorf("Expected fallback to base16, got %s", current.Name)
	}

	// Test valid init
	Init("cyberpunk")
	if current.Name != "cyberpunk" {
		t.Errorf("Expected cyberpunk, got %s", current.Name)
	}

	// Check if global colors were updated
	if PrimaryColor != current.Primary {
		t.Error("PrimaryColor global not updated")
	}
}

func TestBoxStylePadding(t *testing.T) {
	// Verify our recent change to padding
	Init("base16")

	// BoxStyle should have padding 0 vertical, 1 horizontal
	// There is no easy way to inspect private fields of lipgloss.Style directly in test
	// without rendering or relying on internal knowledge,
	// but we can trust the Init function ran without panic.
	// We can check if rendering a string adds the expected padding.

	s := "test"
	rendered := BoxStyle.Render(s)
	width := lipgloss.Width(rendered)

	// Border(2) + Padding(1 left + 1 right) + Content(4) = 8
	// Vertical: Border(2) + Padding(0 top + 0 bottom) + Content(1) = 3

	expectedWidth := 2 + 2 + 4 // 8
	if width != expectedWidth {
		t.Errorf("Expected rendered width %d, got %d", expectedWidth, width)
	}

	height := lipgloss.Height(rendered)
	expectedHeight := 2 + 1 // 3
	if height != expectedHeight {
		t.Errorf("Expected rendered height %d, got %d", expectedHeight, height)
	}
}
