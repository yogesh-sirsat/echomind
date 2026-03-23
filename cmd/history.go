package cmd

import (
	"fmt"
	"os"
	"strings"

	"echomind/internal/config"
	"echomind/internal/ui"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View recording history",
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

		fmt.Println(ui.TitleStyle.Render("📊 Recording History"))
		fmt.Println()

		for _, entry := range history {
			timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
			fmt.Printf("%s %s\n", ui.PromptStyle.Render("Date:"), timestamp)
			fmt.Printf("%s %s\n", ui.StatusStyle.Render("File:"), entry.FileName)
			fmt.Printf("%s %s\n", ui.StatusStyle.Render("Path:"), entry.FilePath)
			fmt.Println(lipgloss.NewStyle().Foreground(ui.MutedColor).Render(strings.Repeat("-", 40)))
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
