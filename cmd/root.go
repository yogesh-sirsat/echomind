package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "echomind",
	Short: "EchoMind is a colorful and animated CLI voice recorder",
	Long: `A visually engaging terminal interface for recording audio, 
built with Go, Bubble Tea, and Lipgloss.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Root flags if any
}
