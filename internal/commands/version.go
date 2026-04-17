package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Kodia",
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("Kodia Framework CLI")
		fmt.Printf("Version: %s\n", "v0.1.0-alpha")
		fmt.Printf("Go version: %s\n", "1.26.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
