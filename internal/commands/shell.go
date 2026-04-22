package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interactive Go REPL for Kodia (Laravel-style)",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🔮 Welcome to Kodia Shell!")
		color.Cyan("Interactive Go REPL. Pre-loaded with stdlib. Type 'exit' to quit.")
		fmt.Println()

		i := interp.New(interp.Options{})

		// Use stdlib symbols
		i.Use(stdlib.Symbols)

		// Basic REPL loop with better feedback
		for {
			fmt.Print(color.GreenString("kodia> "))
			var line string
			
			// We use a simple Scan here for now as a base implementation.
			// In a real pro environment, we'd use a line-editing library.
			fmt.Scanln(&line)

			if line == "" {
				continue
			}

			if line == "exit" || line == "quit" {
				color.Yellow("Goodbye! 👋")
				break
			}

			res, err := i.Eval(line)
			if err != nil {
				color.Red("Error: %v", err)
				continue
			}

			if res.IsValid() {
				color.White("=> %v", res)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
