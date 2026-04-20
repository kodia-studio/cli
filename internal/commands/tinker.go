package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var tinkerCmd = &cobra.Command{
	Use:   "tinker",
	Short: "Interactive REPL shell for Kodia",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🔮 Welcome to Kodia Tinker!")
		color.Cyan("Interactive Go REPL. Type 'exit' to quit.")

		i := interp.New(interp.Options{})

		// Use stdlib
		i.Use(stdlib.Symbols)

		// Basic REPL loop
		for {
			fmt.Print("kodia> ")
			var line string
			_, err := fmt.Scanln(&line)
			if err != nil {
				if err.Error() == "unexpected newline" {
					continue
				}
				break
			}

			if line == "exit" || line == "quit" {
				break
			}

			res, err := i.Eval(line)
			if err != nil {
				color.Red("error: %v", err)
				continue
			}

			if res.IsValid() {
				fmt.Println(res)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(tinkerCmd)
}
