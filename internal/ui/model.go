package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/enqack/rundown/internal/stats"
)

type Tab int

const (
	TabOverview Tab = iota
	TabCPU
	TabMem
	TabDisk
	TabNet
	TabTemp
)

type tickMsg struct {
	Stats stats.SystemStats
	Err   error
	Step  string
}

type Model struct {
	Stats          stats.SystemStats
	CPUProg        progress.Model
	CPUCoresProgs  []progress.Model
	CPUViewport    viewport.Model
	MemProg        progress.Model
	SwapProg       progress.Model
	DiskProg       progress.Model
	DiskProgs      map[string]progress.Model
	NetUpProg      progress.Model
	NetDnProg      progress.Model
	TempViewport   viewport.Model
	TempProgs      map[string]progress.Model
	SortBy         string
	Loading        bool
	LoadingMsg     string
	LoadVal        float64
	LoadProg       progress.Model
	Tab            Tab
	Width          int
	Height         int
	LastSync       time.Time
	UpdateInterval time.Duration
}
