package commands

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kodia-studio/cli/internal/scaffolding"
	"github.com/spf13/cobra"
)

var initDockerCmd = &cobra.Command{
	Use:   "init:docker",
	Short: "Initialize Docker and Docker Compose for Kodia Sail",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("Initializing Kodia Sail (Docker Infrastructure) 🐨⚓")

		// Get project name from current directory
		cwd, _ := os.Getwd()
		projectName := filepath.Base(cwd)

		data := scaffolding.TemplateData{
			Name:        projectName,
			ProjectName: projectName,
		}

		// 1. Generate Backend Dockerfile
		backendDockerfile := filepath.Join("backend", "Dockerfile")
		color.Cyan("Generating %s...", backendDockerfile)
		if err := scaffolding.Generate("Dockerfile.tmpl", backendDockerfile, data); err != nil {
			color.Red("Error generating Dockerfile: %v", err)
			return
		}

		// 2. Generate Docker Compose
		composeFile := "docker-compose.yml"
		color.Cyan("Generating %s...", composeFile)
		if err := scaffolding.Generate("docker-compose.tmpl", composeFile, data); err != nil {
			color.Red("Error generating docker-compose.yml: %v", err)
			return
		}

		// 3. Generate .dockerignore
		ignoreFile := ".dockerignore"
		color.Cyan("Generating %s...", ignoreFile)
		if err := scaffolding.Generate("dockerignore.tmpl", ignoreFile, data); err != nil {
			color.Red("Error generating .dockerignore: %v", err)
			return
		}

		color.Green("\n✨ Kodia Sail initialized successfully! ✅")
		color.Yellow("\nTo start your environment, run:")
		color.Cyan("docker compose up -d")
		color.White("\nServices included: Postgres (5432), Redis (6379), MinIO (9000/9001)")
	},
}

func init() {
	rootCmd.AddCommand(initDockerCmd)
}
