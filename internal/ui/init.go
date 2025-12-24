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
	cp := theme.NewThemedProgress()
	mp := theme.NewThemedProgress()
	sp := theme.NewThemedProgress()
	dp := theme.NewThemedProgress()
	nup := theme.NewThemedProgress()
	ndn := theme.NewThemedProgress()

	return Model{
		CPUProg:          cp,
		CPUCoresProgs:    []progress.Model{},
		CPUViewport:      viewport.New(0, 0),
		MemProg:          mp,
		SwapProg:         sp,
		DiskProg:         dp,
		DiskProgs:        make(map[string]progress.Model),
		NetUpProg:        nup,
		NetDnProg:        ndn,
		IfaceUpProgs:     make(map[string]progress.Model),
		IfaceDnProgs:     make(map[string]progress.Model),
		NetViewport:      viewport.New(0, 0),
		MemViewport:      viewport.New(0, 0),
		DiskViewport:     viewport.New(0, 0),
		OverviewViewport: viewport.New(0, 0),
		TempViewport:     viewport.New(0, 0),
		ProcViewport:     viewport.New(0, 0),
		TempProgs:        make(map[string]progress.Model),
		Loading:          true,
		LoadingMsg:       "Scanning CPU/Memory...",
		LoadingStartTime: time.Now(),
		LoadProg:         theme.NewThemedProgress(),
		SortBy:           "cpu",
		UpdateInterval:   time.Second, // Default 1 second update interval
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickMsgCmd("Scanning CPU/Memory..."),
	)
}

func tickMsgCmd(step string) tea.Cmd {
	return func() tea.Msg {
		// Initial stats collection (no full process details needed yet)
		s, err := stats.GetStats(false)
		return tickMsg{Stats: s, Err: err, Step: step}
	}
}
