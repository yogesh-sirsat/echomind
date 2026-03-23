package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"echomind/internal/ui"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Add EchoMind to your system PATH (Windows only)",
	Run: func(cmd *cobra.Command, args []string) {
		exePath, err := os.Executable()
		if err != nil {
			fmt.Printf("Error getting executable path: %v\n", err)
			return
		}
		exeDir := filepath.Dir(exePath)

		// PowerShell command to add to user PATH if not already there
		psCmd := fmt.Sprintf(`
			$newPath = "%s"
			$oldPath = [Environment]::GetEnvironmentVariable("Path", "User")
			if ($oldPath -notlike "*$newPath*") {
				[Environment]::SetEnvironmentVariable("Path", "$oldPath;$newPath", "User")
				Write-Output "Successfully added to PATH."
			} else {
				Write-Output "Already in PATH."
			}
		`, exeDir)

		out, err := exec.Command("powershell", "-Command", psCmd).CombinedOutput()
		if err != nil {
			fmt.Printf("Error updating PATH: %v\n%s\n", err, string(out))
			return
		}

		fmt.Println(ui.InfoStyle.Render(string(out)))
		fmt.Println(ui.StatusStyle.Render("You may need to restart your terminal for changes to take effect."))
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
