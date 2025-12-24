package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enqack/rundown/internal/stats"
	"github.com/enqack/rundown/internal/theme"
)

func NewModel() Model {
	cp := progress.New(progress.WithDefaultGradient())
	mp := progress.New(progress.WithDefaultGradient())
	sp := progress.New(progress.WithDefaultGradient())
	dp := progress.New(progress.WithDefaultGradient())
	nup := progress.New(progress.WithDefaultGradient())
	ndn := progress.New(progress.WithDefaultGradient())

	// Customizing progress bars with theme colors
	cp.FullColor = string(theme.PrimaryColor)
	mp.FullColor = string(theme.SecondaryColor)
	sp.FullColor = string(theme.WarningColor)
	dp.FullColor = string(theme.AccentColor)
	nup.FullColor = string(theme.WarningColor)
	ndn.FullColor = string(theme.SecondaryColor)

	return Model{
		CPUProg:        cp,
		CPUCoresProgs:  []progress.Model{},
		CPUViewport:    viewport.New(0, 0),
		MemProg:        mp,
		SwapProg:       sp,
		DiskProg:       dp,
		DiskProgs:      make(map[string]progress.Model),
		NetUpProg:      nup,
		NetDnProg:      ndn,
		TempViewport:   viewport.New(0, 0),
		TempProgs:      make(map[string]progress.Model),
		Loading:        true,
		LoadingMsg:     "Scanning CPU/Memory...",
		LoadProg:       progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
		SortBy:         "cpu",
		UpdateInterval: time.Second, // Default 1 second update interval
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickMsgCmd("Scanning CPU/Memory..."),
	)
}

func tickMsgCmd(step string) tea.Cmd {
	return func() tea.Msg {
		s, err := stats.GetStats()
		return tickMsg{Stats: s, Err: err, Step: step}
	}
}

func tick() tea.Cmd {
	return tickMsgCmd("")
}
