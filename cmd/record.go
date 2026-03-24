package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"echomind/internal/audio"
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

		// Check for recovery file
		configDir, _ := config.GetConfigDir()
		recoveryPath := filepath.Join(configDir, "recovery.bin")
		if _, err := os.Stat(recoveryPath); err == nil {
			fmt.Println(ui.PromptStyle.Render("⚠️  Interrupted recording found!"))
			fmt.Print(ui.StatusStyle.Render("Would you like to recover it? (y/n): "))
			
			reader := bufio.NewReader(os.Stdin)
			char, _, _ := reader.ReadRune()
			if strings.ToLower(string(char)) == "y" {
				fmt.Println(ui.InfoStyle.Render("Recovering..."))
				samples, err := audio.LoadRecovery(recoveryPath)
				if err != nil {
					fmt.Printf("Error loading recovery: %v\n", err)
				} else {
					// Save it as a generic recovery file
					timestamp := time.Now().Format("2006-01-02_15-04-05")
					fileName := fmt.Sprintf("recovered_%s.wav", timestamp)
					savePath := filepath.Join(cfg.DefaultDirectory, fileName)
					
					err = audio.SaveSamples(savePath, samples, 44100)
					if err != nil {
						fmt.Printf("Error saving recovery: %v\n", err)
					} else {
						fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("Successfully recovered to: %s", savePath)))
						// Add to history
						_ = config.AddToHistory(config.HistoryEntry{
							Timestamp: time.Now(),
							FileName:  fileName,
							FilePath:  savePath,
							Format:    "wav",
						})
					}
					_ = os.Remove(recoveryPath)
				}
			} else {
				_ = os.Remove(recoveryPath)
				fmt.Println(ui.StatusStyle.Render("Interrupted session discarded."))
			}
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
