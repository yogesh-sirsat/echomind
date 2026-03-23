package cmd

import (
	"fmt"
	"os"

	"echomind/internal/config"
	"echomind/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set default values for format and directory",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		p := tea.NewProgram(ui.InitialConfigModel(cfg))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running config UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
