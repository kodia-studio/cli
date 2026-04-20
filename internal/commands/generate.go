package commands

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kodia-studio/cli/internal/astutil"
	"github.com/kodia-studio/cli/internal/scaffolding"
	"github.com/kodia-studio/cli/internal/validation"
	"github.com/spf13/cobra"
)

var (
	fieldsStr string
	withTests bool
	withMig   bool
	withAuth  bool
)

// generateCmd represents the base command for generating code
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate various framework components",
	Long:  `Generate models, controllers, migrations, CRUD features, and more with a single command.`,
}

// generateCrudCmd represents the crud generator
var generateCrudCmd = &cobra.Command{
	Use:   "crud [Name]",
	Short: "Generate a complete vertical slice CRUD feature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := validation.ValidateName(name); err != nil {
			color.Red("Error: Invalid name - %v", err)
			return
		}

		color.Magenta("🔥 Scaffolding full CRUD feature: %s", name)
		data := scaffolding.BuildData(name, fieldsStr)

		// 1. Backend Layers
		scaffolding.Generate("model.tmpl", filepath.Join("backend", "internal", "core", "domain", data.LowerName+".go"), data)
		scaffolding.Generate("ports.tmpl", filepath.Join("backend", "internal", "core", "ports", data.LowerName+".go"), data)
		scaffolding.Generate("handler.tmpl", filepath.Join("backend", "internal", "adapters", "http", "handlers", data.LowerName+"_handler.go"), data)
		scaffolding.Generate("service.tmpl", filepath.Join("backend", "internal", "core", "services", data.LowerName+"_service.go"), data)
		scaffolding.Generate("repository.tmpl", filepath.Join("backend", "internal", "adapters", "repository", "postgres", data.LowerName+"_repository.go"), data)
		
		// 2. Specialized Layers (Requests, Resources, Policies)
		scaffolding.Generate("request.tmpl", filepath.Join("backend", "internal", "adapters", "http", "dto", data.LowerName+"_request.go"), data)
		scaffolding.Generate("resource.tmpl", filepath.Join("backend", "internal", "adapters", "http", "resources", data.LowerName+"_resource.go"), data)
		scaffolding.Generate("policy.tmpl", filepath.Join("backend", "internal", "core", "policies", data.LowerName+"_policy.go"), data)

		// 3. Migration
		if withMig {
			upDest := filepath.Join("backend", "internal", "infrastructure", "database", "migrations", "sql", data.Timestamp+"_create_"+data.LowerPlural+"_table.up.sql")
			downDest := filepath.Join("backend", "internal", "infrastructure", "database", "migrations", "sql", data.Timestamp+"_create_"+data.LowerPlural+"_table.down.sql")
			scaffolding.Generate("migration_up.tmpl", upDest, data)
			scaffolding.Generate("migration_down.tmpl", downDest, data)
		}

		// 3. Frontend
		makePageCmd.Run(cmd, []string{data.LowerPlural})

		// 4. Tests
		if withTests {
			scaffolding.Generate("service_test.tmpl", filepath.Join("backend", "internal", "core", "services", data.LowerName+"_service_test.go"), data)
			scaffolding.Generate("handler_test.tmpl", filepath.Join("backend", "internal", "adapters", "http", "handlers", data.LowerName+"_handler_test.go"), data)
		}

		// 5. Auto-wiring
		color.Cyan("🪄  Performing auto-wiring magic...")
		mainPath := filepath.Join("backend", "cmd", "server", "main.go")
		astutil.InjectDependencyInjection(mainPath, data)
		
		routerPath := filepath.Join("backend", "internal", "adapters", "http", "router.go")
		astutil.InjectRouteRegistration(routerPath, data)

		color.Magenta("✨ Feature %s fully scaffolded and wired! 🚀", name)
	},
}

// generateModelCmd represents the model generator
var generateModelCmd = &cobra.Command{
	Use:   "model [Name]",
	Short: "Generate a domain model",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name, fieldsStr)
		dest := filepath.Join("backend", "internal", "core", "domain", data.LowerName+".go")
		scaffolding.Generate("model.tmpl", dest, data)
	},
}

// generatePolicyCmd represents the policy generator
var generatePolicyCmd = &cobra.Command{
	Use:   "policy [Name]",
	Short: "Generate an authorization policy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name, "")
		dest := filepath.Join("backend", "internal", "core", "policies", data.LowerName+"_policy.go")
		scaffolding.Generate("policy.tmpl", dest, data)
	},
}

// generateRequestCmd represents the request/DTO generator
var generateRequestCmd = &cobra.Command{
	Use:   "request [Name]",
	Short: "Generate a request DTO",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name, fieldsStr)
		dest := filepath.Join("backend", "internal", "adapters", "http", "dto", data.LowerName+"_request.go")
		scaffolding.Generate("request.tmpl", dest, data)
	},
}

// generateResourceCmd represents the resource transformer generator
var generateResourceCmd = &cobra.Command{
	Use:   "resource [Name]",
	Short: "Generate an API resource transformer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name, fieldsStr)
		dest := filepath.Join("backend", "internal", "adapters", "http", "resources", data.LowerName+"_resource.go")
		scaffolding.Generate("resource.tmpl", dest, data)
	},
}

// generatePackageCmd represents the package generator
var generatePackageCmd = &cobra.Command{
	Use:   "package [Name]",
	Short: "Generate a new internal package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		desc, _ := cmd.Flags().GetString("description")
		data := scaffolding.BuildData(name, "")
		
		// Logic to create a package would go here
		fmt.Printf("Generating package %s: %s\n", data.Name, desc)
	},
}

func init() {
	// Flags
	generateCmd.PersistentFlags().StringVarP(&fieldsStr, "fields", "f", "", "Fields for the entity (e.g. name:string,email:string)")
	
	generateCrudCmd.Flags().BoolVar(&withTests, "with-tests", true, "Generate tests")
	generateCrudCmd.Flags().BoolVar(&withMig, "with-migrations", true, "Generate migrations")
	generateCrudCmd.Flags().BoolVar(&withAuth, "with-auth", false, "Add auth middleware")

	generatePackageCmd.Flags().String("description", "A new package", "Package description")

	// Add subcommands
	generateCmd.AddCommand(generateCrudCmd)
	generateCmd.AddCommand(generateModelCmd)
	generateCmd.AddCommand(generatePolicyCmd)
	generateCmd.AddCommand(generateRequestCmd)
	generateCmd.AddCommand(generateResourceCmd)
	generateCmd.AddCommand(generatePackageCmd)
	
	// Add migrate, event, listener as aliases/subcommands
	generateCmd.AddCommand(makeMigrationCmd)
	generateCmd.AddCommand(makeEventCmd)
	generateCmd.AddCommand(makeListenerCmd)

	rootCmd.AddCommand(generateCmd)
}
