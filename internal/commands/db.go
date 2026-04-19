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

var dbSeedCmd = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with dummy data",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Seeding database with realistic dummy data... 🌱")
		runDbCommand("go", "run", "cmd/seeder/main.go")
	},
}

var dbFreshCmd = &cobra.Command{
	Use:   "db:fresh",
	Short: "Drop all tables and re-run all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("⚠️  This will drop all tables. Are you sure? (no confirmation, use with caution)")
		runDbCommand("make", "db-fresh")
	},
}

var dbStatusCmd = &cobra.Command{
	Use:   "db:status",
	Short: "Show pending and applied migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Checking migration status...")
		color.Yellow("Tip: Run 'kodia db:migrate' to apply pending migrations")
		runDbCommand("make", "db-status")
	},
}

var dbResetCmd = &cobra.Command{
	Use:   "db:reset",
	Short: "Rollback all, re-run migrations, and seed the database",
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("⚠️  Resetting database: rollback → migrate → seed")
		runDbCommand("make", "db-reset")
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start both backend and frontend development servers",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("Starting Kodia Framework Development Mode 🚀")
		color.Yellow("Note: You must have 'make' and 'docker' installed.")

		execCmd := exec.Command("make", "dev")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			color.Red("Failed to start dev servers: %v", err)
		}
	},
}

func findBackendDir() string {
	// Try multiple locations for the backend directory
	possiblePaths := []string{
		"backend",
		"./backend",
		"../backend",
		"../../backend",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			if _, err := os.Stat(filepath.Join(path, "Makefile")); err == nil {
				return path
			}
		}
	}

	// Fallback to "backend" if nothing found
	return "backend"
}

func runDbCommand(command string, args ...string) {
	backendDir := findBackendDir()

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
	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(dbFreshCmd)
	rootCmd.AddCommand(dbStatusCmd)
	rootCmd.AddCommand(dbResetCmd)
	rootCmd.AddCommand(devCmd)
}
