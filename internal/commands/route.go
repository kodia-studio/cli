package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var routeListCmd = &cobra.Command{
	Use:   "route:list",
	Short: "Show all registered routes in the backend",
	Run: func(cmd *cobra.Command, args []string) {
		routerPath := filepath.Join("backend", "internal", "adapters", "http", "router.go")
		
		color.Magenta("🗺️  Kodia Route Map")
		fmt.Println()

		file, err := os.Open(routerPath)
		if err != nil {
			color.Red("Error: Could not find router.go at %s", routerPath)
			return
		}
		defer file.Close()

		// Regex to find routes like: auth.POST("/login", ...) or api.GET("/health", ...)
		// Matches: group.METHOD("/path", ...)
		routeRegex := regexp.MustCompile(`(\w+)\.(GET|POST|PUT|PATCH|DELETE|OPTIONS|HEAD)\("([^"]+)"`)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, color.CyanString("METHOD")+"\t"+color.CyanString("URI")+"\t"+color.CyanString("GROUP"))

		scanner := bufio.NewScanner(file)
		count := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			matches := routeRegex.FindStringSubmatch(line)
			if len(matches) == 4 {
				group := matches[1]
				method := matches[2]
				uri := matches[3]

				// Basic grouping logic for display
				if group == "engine" {
					group = "ROOT"
				}

				methodColor := color.New(color.Bold)
				switch method {
				case "GET":
					methodColor.Add(color.FgGreen)
				case "POST":
					methodColor.Add(color.FgYellow)
				case "PUT", "PATCH":
					methodColor.Add(color.FgCyan)
				case "DELETE":
					methodColor.Add(color.FgRed)
				default:
					methodColor.Add(color.FgWhite)
				}

				fmt.Fprintln(w, methodColor.Sprint(method)+"\t"+uri+"\t"+color.HiBlackString(group))
				count++
			}
		}

		w.Flush()
		fmt.Println()
		color.HiBlack("Total Routes Found: %d", count)
	},
}

func init() {
	rootCmd.AddCommand(routeListCmd)
}
