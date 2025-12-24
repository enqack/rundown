package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enqack/rundown/internal/theme"
	"github.com/enqack/rundown/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	themeName     string
	enableProfile bool
	profilePort   int
	RootCmd       = &cobra.Command{
		Use:   "rundown",
		Short: "Rundown - A beautiful terminal system monitor",
		Long: `Rundown is a fast, beautiful system monitor for your terminal.
It displays CPU, memory, disk, network, and temperature information
in a clean, easy-to-read interface.`,
		Run: runApp,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	RootCmd.PersistentFlags().StringVarP(&themeName, "theme", "t", "base16",
		"Color theme: base16, cyberpunk, monochrome, phosphor")
	RootCmd.PersistentFlags().BoolVar(&enableProfile, "profile", false,
		"Enable pprof profiling server")
	RootCmd.PersistentFlags().IntVar(&profilePort, "profile-port", 6060,
		"Port for pprof profiling server")

	// Bind flags to viper
	_ = viper.BindPFlag("theme", RootCmd.PersistentFlags().Lookup("theme"))
}

func initConfig() {
	// Config file support (optional)
	viper.SetConfigName("rundown")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/rundown")
	viper.AddConfigPath(".")

	// Read config if it exists (don't error if it doesn't)
	_ = viper.ReadInConfig()
}

func runApp(cmd *cobra.Command, args []string) {
	// Get theme from viper (respects config file and flags)
	selectedTheme := viper.GetString("theme")

	// Validate theme name
	if !theme.ValidateTheme(selectedTheme) {
		fmt.Fprintf(os.Stderr, "Error: Invalid theme '%s'\n", selectedTheme)
		fmt.Fprintf(os.Stderr, "Available themes: %v\n", theme.AvailableThemes())
		os.Exit(1)
	}

	// Initialize theme
	theme.Init(selectedTheme)

	// Start profiling server if enabled
	if enableProfile {
		go func() {
			addr := fmt.Sprintf("localhost:%d", profilePort)
			fmt.Fprintf(os.Stderr, "Starting pprof server on http://%s/debug/pprof/\n", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				fmt.Fprintf(os.Stderr, "pprof server error: %v\n", err)
			}
		}()
	}

	// Run application
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
