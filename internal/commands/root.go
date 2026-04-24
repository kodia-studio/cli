package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version bool
	rootCmd = &cobra.Command{
		Use:   "kodia",
		Short: "Kodia is a professional fullstack framework for Go and SvelteKit",
		Long: `Kodia Framework CLI - A powerful tool for scaffolding and managing
your professional fullstack applications built with Golang Gin and SvelteKit.`,
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println("Kodia Framework CLI")
				fmt.Println("Version: v0.1.0-alpha")
				return
			}
			cmd.Help()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "Print version information")
}
