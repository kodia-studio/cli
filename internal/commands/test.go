package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run all backend and frontend tests",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🧪 Running Kodia Test Suite...")

		// 1. Backend Tests
		color.Cyan("\n📦 Testing Backend (Go)...")
		goTest := exec.Command("go", "test", "./...")
		goTest.Dir = "backend"
		goTest.Stdout = os.Stdout
		goTest.Stderr = os.Stderr
		if err := goTest.Run(); err != nil {
			color.Red("❌ Backend tests failed!")
		} else {
			color.Green("✅ Backend tests passed!")
		}

		// 2. Frontend Tests (Optional)
		if _, err := os.Stat("frontend/package.json"); err == nil {
			color.Cyan("\n🎨 Testing Frontend (Vitest/Svelte)...")
			npmTest := exec.Command("npm", "test", "--", "--run")
			npmTest.Dir = "frontend"
			npmTest.Stdout = os.Stdout
			npmTest.Stderr = os.Stderr
			if err := npmTest.Run(); err != nil {
				color.Yellow("⚠️  Frontend tests failed or not configured.")
			} else {
				color.Green("✅ Frontend tests passed!")
			}
		}
	},
}

var testCoverageCmd = &cobra.Command{
	Use:   "test:coverage",
	Short: "Run tests and generate coverage report",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("📊 Generating Coverage Report...")

		// 1. Run tests with coverage
		color.Cyan("\n📦 Running coverage for Backend...")
		coverageFile := "coverage.out"
		goTest := exec.Command("go", "test", "-coverprofile="+coverageFile, "./...")
		goTest.Dir = "backend"
		goTest.Stdout = os.Stdout
		goTest.Stderr = os.Stderr
		if err := goTest.Run(); err != nil {
			color.Red("❌ Coverage run failed!")
			return
		}

		// 2. Display summary
		color.Cyan("\n📈 Coverage Summary:")
		goTool := exec.Command("go", "tool", "cover", "-func="+filepath.Join("backend", coverageFile))
		goTool.Stdout = os.Stdout
		goTool.Stderr = os.Stderr
		goTool.Run()

		// 3. Optional HTML report
		html, _ := cmd.Flags().GetBool("html")
		if html {
			color.Cyan("\n🌐 Opening HTML Coverage Report...")
			htmlTool := exec.Command("go", "tool", "cover", "-html="+filepath.Join("backend", coverageFile))
			htmlTool.Run()
		}
	},
}

func init() {
	testCoverageCmd.Flags().Bool("html", false, "Open the coverage report in the browser")
	
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(testCoverageCmd)
}
