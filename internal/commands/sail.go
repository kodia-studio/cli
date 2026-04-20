package commands

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/kodia-studio/cli/internal/scaffolding"
	"github.com/spf13/cobra"
)

var sailCmd = &cobra.Command{
	Use:   "sail",
	Short: "Docker-based development environment for Kodia",
}

var sailUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the Kodia development environment",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Ensure docker-compose.yml exists
		if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
			color.Cyan("🐳 No docker-compose.yml found. Generating one for you...")
			data := scaffolding.BuildData("Sail", "")
			if err := scaffolding.Generate("docker-compose.tmpl", "docker-compose.yml", data); err != nil {
				color.Red("❌ Error generating docker-compose.yml: %v", err)
				return
			}
		}

		// 2. Run docker compose up
		color.Cyan("🚀 Starting Kodia Sail environment...")
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = "  Orchestrating containers..."
		s.Start()

		dockerCmd := exec.Command("docker", "compose", "up", "-d")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			s.Stop()
			color.Red("❌ Failed to start Docker Compose. Is Docker running? Error: %v", err)
			return
		}

		s.Stop()
		color.Green("✅ Kodia Sail is up and running!")
		fmt.Println("\n📡 Services Status:")
		fmt.Println("  - PostgreSQL:  5432")
		fmt.Println("  - Redis:       6379")
		fmt.Println("  - Meilisearch: 7700")
		fmt.Println("  - Mailpit:     8025 (Dashboard)")
		fmt.Println()
		color.Yellow("Tip: Use 'kodia sail down' to stop the environment.")
	},
}

var sailDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the Kodia development environment",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("🛑 Stopping Kodia Sail environment...")
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = "  Tearing down containers..."
		s.Start()

		dockerCmd := exec.Command("docker", "compose", "down")
		if err := dockerCmd.Run(); err != nil {
			s.Stop()
			color.Red("❌ Failed to stop containers: %v", err)
			return
		}

		s.Stop()
		color.Green("✅ Kodia Sail stopped successfully.")
	},
}

func init() {
	rootCmd.AddCommand(sailCmd)
	sailCmd.AddCommand(sailUpCmd)
	sailCmd.AddCommand(sailDownCmd)
}
