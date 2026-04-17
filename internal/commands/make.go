package commands

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kodia-studio/kodia/internal/scaffolding"
	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "Generate boilerplate code for your project",
}

var makeHandlerCmd = &cobra.Command{
	Use:   "make:handler [Name]",
	Short: "Create a new Gin HTTP handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "adapters", "http", "handlers", data.LowerName+"_handler.go")
		
		color.Cyan("Generating handler for %s...", data.Name)
		if err := scaffolding.Generate("handler.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makeServiceCmd = &cobra.Command{
	Use:   "make:service [Name]",
	Short: "Create a new business logic service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "core", "services", data.LowerName+"_service.go")
		
		color.Cyan("Generating service for %s...", data.Name)
		if err := scaffolding.Generate("service.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makeRepositoryCmd = &cobra.Command{
	Use:   "make:repository [Name]",
	Short: "Create a new database repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "adapters", "repository", "postgres", data.LowerName+"_repository.go")
		
		color.Cyan("Generating repository for %s...", data.Name)
		if err := scaffolding.Generate("repository.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makePageCmd = &cobra.Command{
	Use:   "make:page [route]",
	Short: "Create a new SvelteKit page",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		route := args[0]
		// For pages, the route is usually the lower name or plural name
		data := scaffolding.BuildData(route)
		dest := filepath.Join("frontend", "src", "routes", "(app)", route, "+page.svelte")
		
		color.Cyan("Generating Svelte page for %s...", route)
		if err := scaffolding.Generate("svelte-page.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration [table_name]",
	Short: "Create up/down SQL migration files",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		
		baseDest := filepath.Join("backend", "internal", "infrastructure", "database", "migrations", "sql")
		upDest := filepath.Join(baseDest, data.Timestamp+"_create_"+data.LowerPlural+"_table.up.sql")
		downDest := filepath.Join(baseDest, data.Timestamp+"_create_"+data.LowerPlural+"_table.down.sql")
		
		color.Cyan("Generating migrations for %s...", data.LowerPlural)
		scaffolding.Generate("migration_up.tmpl", upDest, data)
		scaffolding.Generate("migration_down.tmpl", downDest, data)
	},
}

var makeFeatureCmd = &cobra.Command{
	Use:   "make:feature [Name]",
	Short: "Scaffold a complete vertical slice feature (Handler, Service, Repo, DB, Frontend)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		color.Magenta("🔥 Scaffolding full feature: %s", name)
		
		// Run all generators
		makeHandlerCmd.Run(cmd, args)
		makeServiceCmd.Run(cmd, args)
		makeRepositoryCmd.Run(cmd, args)
		makeMigrationCmd.Run(cmd, args)
		
		// Map the frontend route to the lower plural form typically
		data := scaffolding.BuildData(name)
		makePageCmd.Run(cmd, []string{data.LowerPlural})
		
		color.Magenta("✨ Feature %s fully scaffolded!", name)
		color.Yellow("Don't forget to:")
		color.Yellow("1. Add the domain entity to internal/core/domain")
		color.Yellow("2. Add interface definitions to internal/core/ports")
		color.Yellow("3. Register the handler in router.go")
	},
}

func init() {
	// Register commands to the root command directly so users can just do `kodia make:handler` instead of `kodia make handler`
	rootCmd.AddCommand(makeHandlerCmd)
	rootCmd.AddCommand(makeServiceCmd)
	rootCmd.AddCommand(makeRepositoryCmd)
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(makePageCmd)
	rootCmd.AddCommand(makeFeatureCmd)
}
