package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dbMigrateCmd = &cobra.Command{
	Use:   "db:migrate",
	Short: "Run all pending database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Running UP database migrations...")
		runDbCommand("make", "migrate-up")
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "db:rollback",
	Short: "Rollback database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Running DOWN database migrations...")
		runDbCommand("make", "migrate-down")
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start both backend and frontend development servers",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("Starting Kodia Framework Development Mode 🚀")
		color.Yellow("Note: You must have 'make' and 'docker' installed.")
		
		// Typically, this might run the root 'make dev'
		execCmd := exec.Command("make", "dev")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		
		if err := execCmd.Run(); err != nil {
			color.Red("Failed to start dev servers: %v", err)
		}
	},
}

func runDbCommand(command string, args ...string) {
	// Let's assume we are in the root directory and 'backend' is a subfolder
	backendDir := "backend"
	if _, err := os.Stat(backendDir); os.IsNotExist(err) {
		// Possibly we are inside 'kodia-cli' developing the tool? 
		// If so, path should be '../backend'
		pwd, _ := os.Getwd()
		if filepath.Base(pwd) == "kodia-cli" {
			backendDir = "../backend"
		}
	}

	execCmd := exec.Command(command, args...)
	execCmd.Dir = backendDir
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		color.Red("Database command failed: %v", err)
	} else {
		color.Green("Database command completed successfully! ✅")
	}
}

func init() {
	rootCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(dbRollbackCmd)
	rootCmd.AddCommand(devCmd)
}
