package commands

import (
	"path/filepath"

	"github.com/fatih/color"
	"github.com/kodia-studio/cli/internal/astutil"
	"github.com/kodia-studio/cli/internal/scaffolding"
	"github.com/spf13/cobra"
)



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
		
		// Generate Tests
		color.Cyan("🧪 Generating unit tests...")
		if err := scaffolding.Generate("service_test.tmpl", filepath.Join("backend", "internal", "core", "services", data.LowerName+"_service_test.go"), data); err != nil {
			color.Red("Error generating service test: %v", err)
		}
		if err := scaffolding.Generate("handler_test.tmpl", filepath.Join("backend", "internal", "adapters", "http", "handlers", data.LowerName+"_handler_test.go"), data); err != nil {
			color.Red("Error generating handler test: %v", err)
		}

		// Auto-wiring magic
		color.Cyan("🪄  Performing auto-wiring magic...")
		
		mainPath := filepath.Join("backend", "cmd", "server", "main.go")
		if err := astutil.InjectDependencyInjection(mainPath, data); err != nil {
			color.Red("⚠️  Auto-wiring failed for main.go: %v", err)
		} else {
			color.Green("✅ Dependency injection registered in main.go")
		}

		routerPath := filepath.Join("backend", "internal", "adapters", "http", "router.go")
		if err := astutil.InjectRouteRegistration(routerPath, data); err != nil {
			color.Red("⚠️  Auto-wiring failed for router.go: %v", err)
		} else {
			color.Green("✅ Routes registered in router.go")
		}

		color.Magenta("✨ Feature %s fully scaffolded and wired! 🚀", name)
		color.Yellow("Next steps:")
		color.Yellow("1. Add the domain entity to internal/core/domain")
		color.Yellow("2. Add interface definitions to internal/core/ports")
	},
}

var makeAuthCmd = &cobra.Command{
	Use:   "make:auth",
	Short: "Scaffold a complete authentication system (Backend & Frontend)",
	Run: func(cmd *cobra.Command, args []string) {
		color.Magenta("🔐 Scaffolding complete Authentication Ecosystem...")

		data := scaffolding.BuildData("Auth")

		// 1. Backend: Handler, Service, Repository
		scaffolding.Generate("auth_handler.tmpl", filepath.Join("backend", "internal", "adapters", "http", "handlers", "auth_handler.go"), data)
		scaffolding.Generate("auth_service.tmpl", filepath.Join("backend", "internal", "core", "services", "auth_service.go"), data)
		scaffolding.Generate("user_repo.tmpl", filepath.Join("backend", "internal", "adapters", "repository", "postgres", "user_repository.go"), data)
		scaffolding.Generate("refresh_token_repo.tmpl", filepath.Join("backend", "internal", "adapters", "repository", "postgres", "refresh_token_repository.go"), data)

		// 2. Backend: DTO, Middleware
		scaffolding.Generate("auth_dto.tmpl", filepath.Join("backend", "internal", "adapters", "http", "dto", "auth_dto.go"), data)
		scaffolding.Generate("auth_middleware.tmpl", filepath.Join("backend", "internal", "adapters", "http", "middleware", "auth.go"), data)

		// 3. Backend: Migrations
		upDest := filepath.Join("backend", "internal", "infrastructure", "database", "migrations", "sql", data.Timestamp+"_create_auth_tables.up.sql")
		downDest := filepath.Join("backend", "internal", "infrastructure", "database", "migrations", "sql", data.Timestamp+"_create_auth_tables.down.sql")
		scaffolding.Generate("auth_migration_up.tmpl", upDest, data)
		scaffolding.Generate("auth_migration_down.tmpl", downDest, data)

		// 4. Frontend: Pages & Store
		scaffolding.Generate("auth_frontend_login.tmpl", filepath.Join("frontend", "src", "routes", "(auth)", "login", "+page.svelte"), data)
		scaffolding.Generate("auth_frontend_register.tmpl", filepath.Join("frontend", "src", "routes", "(auth)", "register", "+page.svelte"), data)
		scaffolding.Generate("auth_frontend_store.tmpl", filepath.Join("frontend", "src", "lib", "stores", "auth.store.ts"), data)

		// 5. Auto-wiring magic
		color.Cyan("🪄  Performing auto-wiring magic...")

		mainPath := filepath.Join("backend", "cmd", "server", "main.go")
		if err := astutil.InjectAuth(mainPath); err != nil {
			color.Red("⚠️  Auto-wiring failed for main.go: %v", err)
		} else {
			color.Green("✅ Dependency injection registered in main.go")
		}

		routerPath := filepath.Join("backend", "internal", "adapters", "http", "router.go")
		if err := astutil.InjectAuthRoutes(routerPath); err != nil {
			color.Red("⚠️  Auto-wiring failed for router.go: %v", err)
		} else {
			color.Green("✅ Routes registered in router.go")
		}

		color.Magenta("✨ Auth system scaffolded and wired! 🚀")
		color.Yellow("Next steps:")
		color.Yellow("1. Ensure your DB is running and apply migrations: kodia migrate up")
		color.Yellow("2. Install frontend dependencies if needed: cd frontend && npm install")
	},
}

var makeMiddlewareCmd = &cobra.Command{
	Use:   "make:middleware [Name]",
	Short: "Create a new Gin middleware",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "adapters", "http", "middleware", data.LowerName+".go")
		
		color.Cyan("Generating middleware %s...", data.Name)
		if err := scaffolding.Generate("middleware.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		} else {
			color.Yellow("Next step: Apply it in router.go (e.g., group.Use(middleware.%s()))", data.Name)
		}
	},
}

var makeValidatorCmd = &cobra.Command{
	Use:   "make:validator [Name]",
	Short: "Create a new custom validator function",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "pkg", "validator", data.LowerName+".go")
		
		color.Cyan("Generating validator %s...", data.Name)
		if err := scaffolding.Generate("validator.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		} else {
			color.Yellow("Next step: Register it in main.go (e.g., v.RegisterValidation(\"%s\", validator.%s))", data.LowerName, data.Name)
		}
	},
}

var makeJobCmd = &cobra.Command{
	Use:   "make:job [Name]",
	Short: "Create a new background job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "core", "jobs", data.LowerName+"_job.go")
		
		color.Cyan("Generating background job %s...", data.Name)
		if err := scaffolding.Generate("job.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
			return
		}

		// Auto-wiring magic for worker
		color.Cyan("🪄  Performing auto-wiring magic for worker...")
		workerMainPath := filepath.Join("backend", "cmd", "worker", "main.go")
		if err := astutil.InjectJobRegistration(workerMainPath, data, false); err != nil {
			color.Red("⚠️  Auto-wiring failed for worker: %v", err)
		} else {
			color.Green("✅ Job registered in worker processor")
		}
	},
}

var makeCronCmd = &cobra.Command{
	Use:   "make:cron [Name]",
	Short: "Create a new scheduled cron job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		data := scaffolding.BuildData(name)
		dest := filepath.Join("backend", "internal", "core", "jobs", data.LowerName+"_cron.go")
		
		color.Cyan("Generating cron job %s...", data.Name)
		if err := scaffolding.Generate("cron.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
			return
		}

		// Auto-wiring magic for worker
		color.Cyan("🪄  Performing auto-wiring magic for worker...")
		workerMainPath := filepath.Join("backend", "cmd", "worker", "main.go")
		if err := astutil.InjectJobRegistration(workerMainPath, data, true); err != nil {
			color.Red("⚠️  Auto-wiring failed for worker: %v", err)
		} else {
			color.Green("✅ Cron job registered in worker processor")
		}
	},
}

var makeComponentCmd = &cobra.Command{
	Use:   "make:component [path/Name]",
	Short: "Create a new reusable Svelte component",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathName := args[0]
		data := scaffolding.BuildData(filepath.Base(pathName))
		
		// Determine directory and filename
		dir := filepath.Dir(pathName)
		var dest string
		if dir == "." {
			dest = filepath.Join("frontend", "src", "lib", "components", data.Name+".svelte")
		} else {
			dest = filepath.Join("frontend", "src", "lib", "components", dir, data.Name+".svelte")
		}
		
		color.Cyan("Generating component %s...", data.Name)
		if err := scaffolding.Generate("component.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makeLayoutCmd = &cobra.Command{
	Use:   "make:layout [route]",
	Short: "Create a new SvelteKit layout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		route := args[0]
		data := scaffolding.BuildData("Layout")
		dest := filepath.Join("frontend", "src", "routes", route, "+layout.svelte")
		
		color.Cyan("Generating layout for route %s...", route)
		if err := scaffolding.Generate("layout.tmpl", dest, data); err != nil {
			color.Red("Error: %v", err)
		}
	},
}

var makeTestCmd = &cobra.Command{
	Use:   "make:test [type] [name]",
	Short: "Create a new unit test for service or handler",
	Long:  "Type can be 'service' or 'handler'. Name is the feature name (e.g., Product).",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		testType := args[0]
		name := args[1]
		data := scaffolding.BuildData(name)

		var dest string
		var template string

		switch testType {
		case "service":
			dest = filepath.Join("backend", "internal", "core", "services", data.LowerName+"_service_test.go")
			template = "service_test.tmpl"
		case "handler":
			dest = filepath.Join("backend", "internal", "adapters", "http", "handlers", data.LowerName+"_handler_test.go")
			template = "handler_test.tmpl"
		default:
			color.Red("Error: Invalid test type. Use 'service' or 'handler'.")
			return
		}

		color.Cyan("Generating %s test for %s...", testType, data.Name)
		if err := scaffolding.Generate(template, dest, data); err != nil {
			color.Red("Error: %v", err)
		}
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
	rootCmd.AddCommand(makeAuthCmd)
	rootCmd.AddCommand(makeMiddlewareCmd)
	rootCmd.AddCommand(makeValidatorCmd)
	rootCmd.AddCommand(makeJobCmd)
	rootCmd.AddCommand(makeCronCmd)
	rootCmd.AddCommand(makeComponentCmd)
	rootCmd.AddCommand(makeLayoutCmd)
	rootCmd.AddCommand(makeTestCmd)
}
