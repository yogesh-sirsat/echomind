package cmd

import (
	"fmt"
	"os"

	"echomind/internal/config"
	"echomind/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "Start a new recording session",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		p := tea.NewProgram(ui.InitialRecordModel(cfg))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running record UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
