package layout

import (
	"testing"

	"github.com/enqack/rundown/internal/theme"
)

func TestNew(t *testing.T) {
	l := New(100, 50)
	if l.TermW != 100 {
		t.Errorf("Expected TermW 100, got %d", l.TermW)
	}
	if l.TermH != 50 {
		t.Errorf("Expected TermH 50, got %d", l.TermH)
	}
}

func TestUsableWidth(t *testing.T) {
	// Setup theme to predictable state if necessary,
	// but assuming default container style which has Border + No Padding (from view.go override, wait)
	// Actually theme.ContainerStyle in theme.go is: Border(ThickBorder) -> Frame size is usually 2 (1 left, 1 right).

	// Let's rely on the actual calculation: TermW - HorizontalFrameSize
	// ThickBorder is 1 char wide on each side => 2 chars total.

	l := New(100, 50)
	expected := 100 - theme.ContainerStyle.GetHorizontalFrameSize()
	got := l.UsableWidth()

	if got != expected {
		t.Errorf("Expected UsableWidth %d, got %d", expected, got)
	}
}

func TestSplitTwoColumns(t *testing.T) {
	l := New(100, 50)
	left, right := l.SplitTwoColumns(100, 2)
	// 100 - 2 = 98. 98 / 2 = 49. Left=49, Right=49.
	if left != 49 || right != 49 {
		t.Errorf("Expected 49, 49. Got %d, %d", left, right)
	}

	left, right = l.SplitTwoColumns(101, 1)
	// 101 - 1 = 100. 100 / 2 = 50. Left=50, Right=50.
	if left != 50 || right != 50 {
		t.Errorf("Expected 50, 50. Got %d, %d", left, right)
	}
}

func TestGraphWidth(t *testing.T) {
	l := New(100, 50)
	// GraphWidth just returns content width
	if l.GraphWidth(50) != 50 {
		t.Errorf("Expected 50, got %d", l.GraphWidth(50))
	}
	if l.GraphWidth(-10) != 0 {
		t.Errorf("Expected 0 for negative width, got %d", l.GraphWidth(-10))
	}
}
