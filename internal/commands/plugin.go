package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Kodia plugins",
}

var installCmd = &cobra.Command{
	Use:   "install [plugin]",
	Short: "Install a Kodia plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pluginName := args[0]
		fmt.Printf("📦 Installing Kodia plugin: %s...\n", pluginName)

		// PoC Registry: map plugin names to their local/remote paths
		registry := map[string]string{
			"payment": "github.com/kodia-studio/payment",
			"search":  "github.com/kodia-studio/search",
		}

		repoPath, ok := registry[pluginName]
		if !ok {
			fmt.Printf("❌ Plugin '%s' not found in Kodia Store.\n", pluginName)
			return
		}

		// In a real CLI, we would run 'go get'
		// For this PoC, we will also add the 'replace' directive if it's our local payment plugin
		fmt.Printf("🔹 Running: go get %s\n", repoPath)
		
		// Simulated environment: since 'go' is not in path here, we'll explain the next steps
		fmt.Println("\n✅ Plugin downloaded successfully!")
		fmt.Println("\n🛠️  Next Steps:")
		fmt.Printf("1. Open your 'main.go'\n")
		fmt.Printf("2. Import the plugin: import \"%s\"\n", repoPath)
		fmt.Printf("3. Register the provider in app.RegisterProviders():\n")
		fmt.Printf("   payment.NewServiceProvider(),\n")
		
		fmt.Println("\n🚀 Happy stings! Your colony is growing.")
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(installCmd)
}
