package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var envSetCmd = &cobra.Command{
	Use:   "env:set [KEY=VALUE]",
	Short: "Set an environment variable in the .env file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pair := args[0]
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			color.Red("Error: Invalid format. Use KEY=VALUE")
			return
		}

		key := strings.ToUpper(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])

		envPath := filepath.Join("backend", ".env")
		if err := updateEnvFile(envPath, key, value); err != nil {
			color.Red("Error updating .env: %v", err)
		} else {
			color.Green("✅ Environment variable %s set successfully!", key)
		}
	},
}

var envListCmd = &cobra.Command{
	Use:   "env:list",
	Short: "List all environment variables in the .env file",
	Run: func(cmd *cobra.Command, args []string) {
		envPath := filepath.Join("backend", ".env")
		file, err := os.Open(envPath)
		if err != nil {
			color.Red("Error: Could not find .env file at %s", envPath)
			return
		}
		defer file.Close()

		color.Magenta("🌍 Kodia Environment Variables")
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, color.CyanString("KEY")+"\t"+color.CyanString("VALUE"))

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				
				// Mask sensitive values
				if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "secret") || strings.Contains(strings.ToLower(key), "password") {
					value = "********"
				}

				fmt.Fprintln(w, key+"\t"+color.HiBlackString(value))
			}
		}
		w.Flush()
		fmt.Println()
	},
}

func updateEnvFile(path, key, value string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new if doesn't exist
			return os.WriteFile(path, []byte(fmt.Sprintf("%s=%s\n", key, value)), 0644)
		}
		return err
	}

	lines := strings.Split(string(input), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	output := strings.Join(lines, "\n")
	return os.WriteFile(path, []byte(output), 0644)
}

func init() {
	rootCmd.AddCommand(envSetCmd)
	rootCmd.AddCommand(envListCmd)
}
