package cmd

import (
	"fmt"
	"os"

	"echomind/internal/config"
	"echomind/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View and browse recording history",
	Run: func(cmd *cobra.Command, args []string) {
		history, err := config.LoadHistory()
		if err != nil {
			fmt.Printf("Error loading history: %v\n", err)
			os.Exit(1)
		}

		if len(history) == 0 {
			fmt.Println(ui.InfoStyle.Render("No recordings found in history."))
			return
		}

		p := tea.NewProgram(ui.InitialHistoryModel(history), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running history UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
