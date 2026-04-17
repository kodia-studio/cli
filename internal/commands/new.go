package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new Kodia project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		projectPath := projectName

		color.Cyan("🚀 Creating new Kodia project: %s", projectName)
		
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		s.Suffix = "  Cloning Kodia template from GitHub..."
		s.Start()
		
		// 1. Clone repository from GitHub
		// Menggunakan branch main dari repo Anda
		cloneCmd := exec.Command("git", "clone", "https://github.com/andiaryatno/framework-kodia.git", projectPath)
		if err := cloneCmd.Run(); err != nil {
			s.Stop()
			color.Red("Failed to clone repository. Is git installed? Error: %v", err)
			return
		}
		
		s.Stop()
		color.Green("✅ Template downloaded successfully!")

		s.Suffix = "  Cleaning up template files..."
		s.Restart()
		
		// 2. Remove .git to start fresh, and remove the CLI folder from the clone
		os.RemoveAll(filepath.Join(projectPath, ".git"))
		os.RemoveAll(filepath.Join(projectPath, "kodia-cli")) // Not needed for the end user project
		
		s.Stop()
		
		s.Suffix = "  Initializing new Git repository..."
		s.Restart()
		
		// 3. Init new git
		exec.Command("git", "-C", projectPath, "init").Run()
		
		time.Sleep(500 * time.Millisecond)
		s.Stop()
		color.Green("✅ Fresh Git repository initialized!")

		fmt.Println()
		color.Yellow("Next steps:")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Printf("  2. Perbarui nama module di backend/go.mod dengan nama project Anda\n")
		fmt.Printf("  3. kodia dev\n")
		fmt.Println()
		color.Cyan("Happy coding with Kodia! 🐨")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
