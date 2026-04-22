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

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile backend and frontend into a single production-ready artifact",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🏗️  Starting Kodia Elite Build Pipeline...")
		startTime := time.Now()

		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		
		// 1. Build Frontend
		s.Suffix = "  Step 1: Compiling Frontend (SvelteKit)..."
		s.Start()
		npmBuild := exec.Command("npm", "run", "build")
		npmBuild.Dir = "frontend"
		if err := npmBuild.Run(); err != nil {
			s.Stop()
			color.Red("❌ Frontend build failed! Check your Svelte scripts.")
			return
		}
		s.Stop()
		color.Green("✅ Frontend compiled successfully!")

		// 2. Bundle Assets
		s.Suffix = "  Step 2: Bundling static assets for embedding..."
		s.Restart()
		distDir := filepath.Join("backend", "internal", "infrastructure", "static", "dist")
		os.RemoveAll(distDir)
		os.MkdirAll(distDir, 0755)

		// Copy from frontend/build to backend
		copyCmd := exec.Command("cp", "-R", "frontend/build/.", distDir)
		if err := copyCmd.Run(); err != nil {
			s.Stop()
			color.Red("❌ Failed to bundle assets: %v", err)
			return
		}
		s.Stop()
		color.Green("✅ Assets bundled for embedding!")

		// 3. Build Backend Binary
		s.Suffix = "  Step 3: Compiling Backend binary (Go)..."
		s.Restart()
		binaryName := "bin/kodia-app"
		goBuild := exec.Command("go", "build", "-ldflags", "-s -w", "-o", binaryName, "./cmd/server/main.go")
		goBuild.Dir = "backend"
		if err := goBuild.Run(); err != nil {
			s.Stop()
			color.Red("❌ Backend compilation failed!")
			return
		}
		s.Stop()
		color.Green("✅ Backend binary generated at backend/%s", binaryName)

		fmt.Println()
		color.Magenta("✨ Build Complete! 🚀")
		color.Cyan("Stats:")
		fmt.Printf("  - Total Time: %v\n", time.Since(startTime).Round(time.Second))
		fmt.Printf("  - Artifact:   backend/%s\n", binaryName)
		fmt.Printf("  - Profile:    Production\n")
		fmt.Println()
		color.HiBlack("Your application is now ready for deployment as a single binary.")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
