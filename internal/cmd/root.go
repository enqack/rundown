package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/enqack/rundown/internal/theme"
	"github.com/enqack/rundown/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	themeName string
	RootCmd   = &cobra.Command{
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

	// Bind flags to viper
	viper.BindPFlag("theme", RootCmd.PersistentFlags().Lookup("theme"))
}

func initConfig() {
	// Config file support (optional)
	viper.SetConfigName("rundown")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/rundown")
	viper.AddConfigPath(".")

	// Read config if it exists (don't error if it doesn't)
	viper.ReadInConfig()
}

func runApp(cmd *cobra.Command, args []string) {
	// Get theme from viper (respects config file and flags)
	selectedTheme := viper.GetString("theme")

	// Initialize theme
	theme.Init(selectedTheme)

	// Run application
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
