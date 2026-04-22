package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run all pending database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Running UP database migrations...")
		runDbCommand("make", "migrate-up")
	},
}

var migrateRollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
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

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show pending and applied migrations",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Checking migration status...")
		color.Yellow("Tip: Run 'kodia migrate' to apply pending migrations")
		runDbCommand("make", "db-status")
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
	// Register commands to the root command directly
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(migrateRollbackCmd)
	rootCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(dbSeedCmd)
	rootCmd.AddCommand(dbFreshCmd)
	rootCmd.AddCommand(dbResetCmd)
	rootCmd.AddCommand(devCmd)
}
