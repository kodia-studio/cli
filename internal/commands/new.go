package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		// Using the main branch from boilerplate repo
		cloneCmd := exec.Command("git", "clone", "https://github.com/kodia-studio/kodia.git", projectPath)
		if err := cloneCmd.Run(); err != nil {
			s.Stop()
			color.Red("Failed to clone repository. Is git installed? Error: %v", err)
			return
		}
		
		s.Stop()
		color.Green("✅ Template downloaded successfully!")

		s.Suffix = "  Cleaning up template files..."
		s.Restart()
		
		// 2. Remove .git to start fresh
		os.RemoveAll(filepath.Join(projectPath, ".git"))
		
		s.Stop()
		
		s.Suffix = "  Initializing new Git repository..."
		s.Restart()
		
		// 3. Init new git
		exec.Command("git", "-C", projectPath, "init").Run()
		
		time.Sleep(500 * time.Millisecond)
		s.Stop()
		color.Green("✅ Fresh Git repository initialized!")

		// 4. Update module name and imports automatically
		newModuleName := projectName
		s.Suffix = "  Updating module name and imports..."
		s.Restart()

		if err := updateModuleInProject(projectPath, newModuleName); err != nil {
			s.Stop()
			color.Red("Failed to update module name: %v", err)
			return
		}

		// Run go mod tidy
		if err := exec.Command("go", "-C", filepath.Join(projectPath, "backend"), "mod", "tidy").Run(); err != nil {
			color.Yellow("⚠️  go mod tidy failed, but continuing...")
		}

		s.Stop()
		color.Green("✅ Go module name updated successfully!")

		// 5. Setup .env for backend
		s.Suffix = "  Setting up environment variables..."
		s.Restart()
		backendEnvExamplePath := filepath.Join(projectPath, "backend", ".env.example")
		backendEnvPath := filepath.Join(projectPath, "backend", ".env")
		if _, err := os.Stat(backendEnvExamplePath); err == nil {
			input, _ := os.ReadFile(backendEnvExamplePath)
			os.WriteFile(backendEnvPath, input, 0644)
			color.Green("✅ Backend .env file created from .env.example")
		} else {
			color.Yellow("⚠️  Backend .env.example not found, skipping .env setup")
		}

		// 6. Setup .env for frontend
		s.Suffix = "  Setting up frontend environment variables..."
		s.Restart()
		frontendEnvExamplePath := filepath.Join(projectPath, "frontend", ".env.example")
		frontendEnvPath := filepath.Join(projectPath, "frontend", ".env")
		if _, err := os.Stat(frontendEnvExamplePath); err == nil {
			input, _ := os.ReadFile(frontendEnvExamplePath)
			os.WriteFile(frontendEnvPath, input, 0644)
			s.Stop()
			color.Green("✅ Frontend .env file created from .env.example")
		} else {
			s.Stop()
			color.Yellow("⚠️  Frontend .env.example not found, skipping frontend .env setup")
		}

		fmt.Println()
		color.Magenta("✨ Success! Your Kodia project is ready. 🚀")
		fmt.Println()
		color.Cyan("Quick Start:")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Print("  2. kodia key:generate   " + color.HiBlackString("(Secure your app with encryption keys)") + "\n")
		fmt.Print("  3. kodia dev            " + color.HiBlackString("(Start development servers)") + "\n")
		fmt.Println()
		color.HiBlack("✅ Backend configured at:  http://localhost:8080")
		color.HiBlack("✅ Frontend configured at: http://localhost:5173")
		color.HiBlack("✅ Database: SQLite is already configured and ready to use!")
		fmt.Println()
		color.HiBlack("📚 For detailed documentation, visit: http://localhost:5173/docs")
		fmt.Println()
	},
}

func updateModuleInProject(projectPath, newModuleName string) error {
	oldModule := "github.com/kodia-studio/kodia"
	backendPath := filepath.Join(projectPath, "backend")

	return filepath.Walk(backendPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .go files and go.mod
		if info.IsDir() || (!strings.HasSuffix(path, ".go") && filepath.Base(path) != "go.mod") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		newContent := strings.ReplaceAll(string(content), oldModule, newModuleName)
		if newContent != string(content) {
			return os.WriteFile(path, []byte(newContent), info.Mode())
		}

		return nil
	})
}

func init() {
	rootCmd.AddCommand(newCmd)
}
