package layout

import (
	"github.com/enqack/rundown/internal/theme"
)

// Calc represents the layout calculator for the application.
// It holds the raw terminal dimensions and provides methods for precise component sizing.
type Calc struct {
	TermW int
	TermH int
}

// New creates a new layout calculator with the given terminal dimensions.
func New(w, h int) *Calc {
	return &Calc{TermW: w, TermH: h}
}

// UsableWidth returns the inner width of the main container.
// Logic: TerminalWidth - ContainerFrame.
func (l *Calc) UsableWidth() int {
	// TerminalWidth - ContainerFrame - SafetyBuffer(2)
	w := l.TermW - theme.ContainerStyle.GetHorizontalFrameSize()

	if w < 20 {
		return 20
	}
	return w
}

// UsableHeight returns the inner height of the main container, safe for content.
// Logic: TerminalHeight - ContainerFrame - Header/Footer Chrome - SafetyMargin.
func (l *Calc) UsableHeight() int {
	// Full terminal - Container Frame - Header (Tabs(2)+Blank(2)=4) - Footer (1) - Safety(2)
	h := l.TermH - theme.ContainerStyle.GetVerticalFrameSize() - 7
	if h < 5 {
		return 5
	}
	return h
}

// ContainerHeight returns the strictly safe height for the container to avoid top-overrun.
// Logic: TerminalHeight - SafetyMargin(1) - ContainerFrame.
func (l *Calc) ContainerHeight() int {
	return l.TermH - 1 - theme.ContainerStyle.GetVerticalFrameSize()
}

// InnerHeight returns the raw inner height of the container structure.
// Used for gap calculations.
// Logic: TerminalHeight - SafetyMargin(1) - ContainerFrame.
func (l *Calc) InnerHeight() int {
	return l.TermH - 1 - theme.ContainerStyle.GetVerticalFrameSize()
}

// BoxContentWidth returns the width available for content INSIDE a BoxStyle component.
// Logic: OuterWidth - BoxFrame.
func (l *Calc) BoxContentWidth(outerWidth int) int {
	w := outerWidth - theme.BoxStyle.GetHorizontalFrameSize()
	if w < 0 {
		return 0
	}
	return w
}

// GraphWidth returns the safe width for a graph/progress bar starting from the CONTENT width of its parent.
// Logic: ContentWidth.
func (l *Calc) GraphWidth(contentWidth int) int {
	if contentWidth < 0 {
		return 0
	}
	return contentWidth
}

// SplitTwoColumns divides a total width into two columns with an optional gap.
// Returns the outer widths of the two columns (Left, Right).
func (l *Calc) SplitTwoColumns(totalWidth int, gap int) (int, int) {
	avail := totalWidth - gap
	if avail < 0 {
		return 0, 0
	}
	left := avail / 2
	right := avail - left
	return left, right
}

// SplitThreeColumns divides a total width into three columns with optional gaps.
// Returns the outer widths of the three columns (Left, Middle, Right).
func (l *Calc) SplitThreeColumns(totalWidth int, gap int) (int, int, int) {
	avail := totalWidth - (gap * 2)
	if avail < 0 {
		return 0, 0, 0
	}
	col := avail / 3
	left := col
	mid := col
	right := avail - left - mid
	return left, mid, right
}

func (l *Calc) CalculateTableDynamicWidth(totalWidth int, fixedColumnsWidths []int, columnPadding int) int {
	used := 0
	for _, w := range fixedColumnsWidths {
		used += w
	}

	// Add padding overhead for all columns (fixed + 1 dynamic)
	totalColumns := len(fixedColumnsWidths) + 1
	paddingOverhead := totalColumns * columnPadding

	rem := totalWidth - used - paddingOverhead
	if rem < 1 {
		return 1
	}
	return rem
}

// SplitColumns divides a total width into N columns with a uniform gap.
// Returns a slice of outer widths for each column.
func (l *Calc) SplitColumns(totalWidth int, count int, gap int) []int {
	if count <= 0 {
		return nil
	}
	if count == 1 {
		return []int{totalWidth}
	}

	// Total width available for content = Total - ((N-1) * gap)
	totalGap := (count - 1) * gap
	avail := totalWidth - totalGap
	if avail < 0 {
		// Degenerate case: return 0 widths
		widths := make([]int, count)
		return widths
	}

	// Base width for each column
	base := avail / count
	remainder := avail % count

	widths := make([]int, count)
	for i := 0; i < count; i++ {
		widths[i] = base
		// Distribute remainder pixels to the first few columns
		if i < remainder {
			widths[i]++
		}
	}
	return widths
}
