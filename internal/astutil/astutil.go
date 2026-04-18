package astutil

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"strings"

	"github.com/kodia-studio/cli/internal/scaffolding"
	"github.com/kodia-studio/cli/internal/validation"
)

// validateTemplateData ensures template data is safe for code generation.
// This prevents injection attacks through AST manipulation.
func validateTemplateData(data scaffolding.TemplateData) error {
	// Validate Name - used in code generation (PascalCase)
	if err := validation.ValidateName(data.Name); err != nil {
		return fmt.Errorf("invalid template data - Name: %w", err)
	}

	// Validate LowerName - used as variable names (lowercase, alphanumeric + underscore)
	if err := validation.ValidateIdentifier(data.LowerName); err != nil {
		return fmt.Errorf("invalid template data - LowerName: %w", err)
	}

	// Validate Plural - used in code generation
	if err := validation.ValidateName(data.Plural); err != nil {
		return fmt.Errorf("invalid template data - Plural: %w", err)
	}

	// Validate LowerPlural - used as variable/constant names
	if err := validation.ValidateIdentifier(data.LowerPlural); err != nil {
		return fmt.Errorf("invalid template data - LowerPlural: %w", err)
	}

	// ProjectName should be a valid identifier
	if err := validation.ValidateIdentifier(data.ProjectName); err != nil {
		return fmt.Errorf("invalid template data - ProjectName: %w", err)
	}

	return nil
}

// InjectDependencyInjection updates main.go with new repo, service, and handler
func InjectDependencyInjection(mainPath string, data scaffolding.TemplateData) error {
	// Validate template data to prevent code injection
	if err := validateTemplateData(data); err != nil {
		return err
	}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, mainPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse main.go: %w", err)
	}

	mainFunc := findFunc(node, "main")
	if mainFunc == nil {
		return fmt.Errorf("could not find main function")
	}

	// 1. Inject Repository
	repoStmt := parseStmt(fmt.Sprintf("%sRepo := postgres.New%sRepository(db)", data.LowerName, data.Name))
	insertInBlock(node, mainFunc.Body, "// 7. Initialize repositories", repoStmt)

	// 2. Inject Service
	serviceStmt := parseStmt(fmt.Sprintf("%sService := services.New%sService(%sRepo, log)", data.LowerName, data.Name, data.LowerName))
	insertInBlock(node, mainFunc.Body, "// 8. Initialize services", serviceStmt)

	// 3. Inject Handler
	handlerStmt := parseStmt(fmt.Sprintf("%sHandler := handlers.New%sHandler(%sService, validate, log)", data.LowerName, data.Name, data.LowerName))
	insertInBlock(node, mainFunc.Body, "// 9. Initialize handlers", handlerStmt)

	// 4. Update NewRouter call
	updateNewRouterCall(mainFunc.Body, data.LowerName+"Handler")

	return saveFile(mainPath, fset, node)
}

// InjectAuth handles the specialized injection for the authentication system
func InjectAuth(mainPath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, mainPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse main.go: %w", err)
	}

	mainFunc := findFunc(node, "main")
	if mainFunc == nil {
		return fmt.Errorf("could not find main function")
	}

	// 1. Inject Repositories
	userRepoStmt := parseStmt("userRepo := postgres.NewUserRepository(db)")
	insertInBlock(node, mainFunc.Body, "// 7. Initialize repositories", userRepoStmt)
	refreshRepoStmt := parseStmt("refreshRepo := postgres.NewRefreshTokenRepository(db)")
	insertInBlock(node, mainFunc.Body, "// 7. Initialize repositories", refreshRepoStmt)

	// 2. Inject Service
	authServiceStmt := parseStmt("authService := services.NewAuthService(userRepo, refreshRepo, jwtManager, log)")
	insertInBlock(node, mainFunc.Body, "// 8. Initialize services", authServiceStmt)

	// 3. Inject Handler
	authHandlerStmt := parseStmt("authHandler := handlers.NewAuthHandler(authService, validate, log)")
	insertInBlock(node, mainFunc.Body, "// 9. Initialize handlers", authHandlerStmt)

	// 4. Update NewRouter call
	updateNewRouterCall(mainFunc.Body, "authHandler")

	return saveFile(mainPath, fset, node)
}

// InjectJobRegistration updates worker/main.go with new job handler
func InjectJobRegistration(workerMainPath string, data scaffolding.TemplateData, isCron bool) error {
	// Validate template data to prevent code injection
	if err := validateTemplateData(data); err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, workerMainPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse worker/main.go: %w", err)
	}

	mainFunc := findFunc(node, "main")
	if mainFunc == nil {
		return fmt.Errorf("could not find main function in worker/main.go")
	}

	importPath := "github.com/kodia-studio/kodia/internal/core/jobs"
	addImport(node, importPath)

	var registrationCode string
	if isCron {
		registrationCode = fmt.Sprintf("processor.Register(jobs.Type%sCron, jobs.Handle%sCronTask)", data.Name, data.Name)
	} else {
		registrationCode = fmt.Sprintf("processor.Register(jobs.Type%s, jobs.Handle%sTask)", data.Name, data.Name)
	}

	registrationStmt := parseStmt(registrationCode)
	insertInBlock(node, mainFunc.Body, "// --- Job Registration Start ---", registrationStmt)

	return saveFile(workerMainPath, fset, node)
}

// addImport adds a new import to the file if it doesn't exist
func addImport(node *ast.File, path string) {
	for _, imp := range node.Imports {
		if imp.Path.Value == fmt.Sprintf("\"%s\"", path) {
			return
		}
	}

	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%s\"", path),
		},
	}

	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{importSpec},
	}

	// Insert after existing imports
	idx := -1
	for i, decl := range node.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
			idx = i
		}
	}

	if idx == -1 {
		node.Decls = append([]ast.Decl{importDecl}, node.Decls...)
	} else {
		node.Decls = append(node.Decls[:idx+1], append([]ast.Decl{importDecl}, node.Decls[idx+1:]...)...)
	}
}

// Helper to find a function by name
func findFunc(node *ast.File, name string) *ast.FuncDecl {
	for _, decl := range node.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == name {
			return fd
		}
	}
	return nil
}

// Helper to parse a single statement string into AST
func parseStmt(code string) ast.Stmt {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package p; func _() { "+code+" }", 0)
	if err != nil {
		return nil
	}
	if len(f.Decls) == 0 {
		return nil
	}
	fd, ok := f.Decls[0].(*ast.FuncDecl)
	if !ok || len(fd.Body.List) == 0 {
		return nil
	}
	return fd.Body.List[0]
}

// Helper to insert a statement after a comment marker
func insertInBlock(file *ast.File, body *ast.BlockStmt, marker string, stmt ast.Stmt) {
	if stmt == nil {
		return
	}

	// Duplicate check for variable assignments
	if as, ok := stmt.(*ast.AssignStmt); ok {
		for _, nameExpr := range as.Lhs {
			if ident, ok := nameExpr.(*ast.Ident); ok {
				// Check if this variable is already defined in the block
				for _, existingStmt := range body.List {
					if eas, ok := existingStmt.(*ast.AssignStmt); ok {
						for _, eNameExpr := range eas.Lhs {
							if eIdent, ok := eNameExpr.(*ast.Ident); ok && eIdent.Name == ident.Name {
								return // Already exists
							}
						}
					}
				}
			}
		}
	}

	var markerPos token.Pos
	for _, cg := range file.Comments {
		for _, c := range cg.List {
			if strings.Contains(c.Text, marker) {
				markerPos = c.Pos()
				break
			}
		}
	}

	if markerPos == token.NoPos {
		body.List = append(body.List, stmt)
		return
	}

	insertIdx := -1
	for i, s := range body.List {
		if s.Pos() > markerPos {
			insertIdx = i
			break
		}
	}

	if insertIdx == -1 {
		body.List = append(body.List, stmt)
	} else {
		body.List = append(body.List[:insertIdx], append([]ast.Stmt{stmt}, body.List[insertIdx:]...)...)
	}
}

// Helper: Update NewRouter call in main.go
func updateNewRouterCall(body *ast.BlockStmt, handlerVar string) {
	for _, stmt := range body.List {
		ast.Inspect(stmt, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if sel.Sel.Name == "NewRouter" {
				// Avoid duplicate
				for _, arg := range call.Args {
					if ident, ok := arg.(*ast.Ident); ok && ident.Name == handlerVar {
						return false
					}
				}
				call.Args = append(call.Args, ast.NewIdent(handlerVar))
				return false
			}
			return true
		})
	}
}

// InjectRouteRegistration updates router.go with new handler field and routes
func InjectRouteRegistration(routerPath string, data scaffolding.TemplateData) error {
	// Validate template data to prevent code injection
	if err := validateTemplateData(data); err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse router.go: %w", err)
	}

	// 1. Add field to Router struct
	addStructField(node, "Router", data.LowerName+"Handler", "*handlers."+data.Name+"Handler")

	// 2. Add param to NewRouter function
	addFuncParam(node, "NewRouter", data.LowerName+"Handler", "*handlers."+data.Name+"Handler")

	// 3. Add assignment in NewRouter
	addFuncAssignment(node, "NewRouter", data.LowerName+"Handler")

	// 4. Add route group in Setup
	addRouteGroup(node, data)

	return saveFile(routerPath, fset, node)
}

// InjectAuthRoutes handles the specialized route registration for authentication
func InjectAuthRoutes(routerPath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, routerPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse router.go: %w", err)
	}

	// 1. Add field to Router struct
	addStructField(node, "Router", "authHandler", "*handlers.AuthHandler")

	// 2. Add param to NewRouter function
	addFuncParam(node, "NewRouter", "authHandler", "*handlers.AuthHandler")

	// 3. Add assignment in NewRouter
	addFuncAssignment(node, "NewRouter", "authHandler")

	// 4. Add auth route group in Setup
	addAuthRouteGroup(node)

	return saveFile(routerPath, fset, node)
}

func addAuthRouteGroup(file *ast.File) {
	fd := findFunc(file, "Setup")
	if fd == nil {
		return
	}

	// Find the 'api' block
	var targetBlock *ast.BlockStmt
	for _, stmt := range fd.Body.List {
		if bs, ok := stmt.(*ast.BlockStmt); ok {
			targetBlock = bs
		}
	}
	if targetBlock == nil {
		targetBlock = fd.Body
	}

	authRouteCode := `{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/logout", r.jwtManagerAuthMiddleware(), r.authHandler.Logout)

			protectedAuth := auth.Group("")
			protectedAuth.Use(r.jwtManagerAuthMiddleware())
			{
				protectedAuth.POST("/logout-all", r.authHandler.LogoutAll)
				protectedAuth.GET("/me", r.authHandler.Me)
			}
		}
	}`

	stmt := parseStmt(authRouteCode)
	if stmt != nil {
		// Insert before the last statement (return engine)
		if len(targetBlock.List) > 0 {
			targetBlock.List = append(targetBlock.List[:len(targetBlock.List)-1], stmt, targetBlock.List[len(targetBlock.List)-1])
		} else {
			targetBlock.List = append(targetBlock.List, stmt)
		}
	}
}

// Helper: Add field to struct
func addStructField(node *ast.File, structName string, fieldName string, fieldType string) {
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name.Name != structName {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			
			// Duplicate check
			for _, f := range st.Fields.List {
				for _, name := range f.Names {
					if name.Name == fieldName {
						return
					}
				}
			}

			// Add field using parser.ParseExpr for safety
			expr, _ := parser.ParseExpr(fieldType)
			st.Fields.List = append(st.Fields.List, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent(fieldName)},
				Type:  expr,
			})
		}
	}
}

// Helper: Add param to function
func addFuncParam(node *ast.File, funcName string, paramName string, paramType string) {
	fd := findFunc(node, funcName)
	if fd == nil {
		return
	}

	// Duplicate check
	for _, f := range fd.Type.Params.List {
		for _, name := range f.Names {
			if name.Name == paramName {
				return
			}
		}
	}

	expr, _ := parser.ParseExpr(paramType)
	fd.Type.Params.List = append(fd.Type.Params.List, &ast.Field{
		Names: []*ast.Ident{ast.NewIdent(paramName)},
		Type:  expr,
	})
}

// Helper: Add assignment in function body
func addFuncAssignment(node *ast.File, funcName string, fieldName string) {
	fd := findFunc(node, funcName)
	if fd == nil {
		return
	}

	for _, stmt := range fd.Body.List {
		ret, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			continue
		}

		for _, expr := range ret.Results {
			unary, ok := expr.(*ast.UnaryExpr)
			if !ok || unary.Op != token.AND {
				continue
			}

			cl, ok := unary.X.(*ast.CompositeLit)
			if !ok {
				continue
			}
			
			// Duplicate check
			for _, elt := range cl.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if ident, ok := kv.Key.(*ast.Ident); ok && ident.Name == fieldName {
						return
					}
				}
			}

			cl.Elts = append(cl.Elts, &ast.KeyValueExpr{
				Key:   ast.NewIdent(fieldName),
				Value: ast.NewIdent(fieldName),
			})
		}
	}
}

// Helper: Add route group in Setup method
func addRouteGroup(file *ast.File, data scaffolding.TemplateData) {
	fd := findFunc(file, "Setup")
	if fd == nil {
		return
	}

	// Heuristic to find the 'api' grouping block
	var targetBlock *ast.BlockStmt
	for _, stmt := range fd.Body.List {
		// Look for blocks that aren't the top level return or simple expressions
		if bs, ok := stmt.(*ast.BlockStmt); ok {
			targetBlock = bs
		}
	}

	// If we still didn't find a nested block, use the body
	if targetBlock == nil {
		targetBlock = fd.Body
	}

	routeCode := fmt.Sprintf(`{
		%s := api.Group("/%s")
		{
			%s.GET("", r.%sHandler.GetAll)
			%s.GET("/:id", r.%sHandler.GetByID)
			%s.POST("", r.%sHandler.Create)
			%s.PATCH("/:id", r.%sHandler.Update)
			%s.DELETE("/:id", r.%sHandler.Delete)
		}
	}`, data.LowerPlural, data.LowerPlural, data.LowerPlural, data.LowerName, data.LowerPlural, data.LowerName, data.LowerPlural, data.LowerName, data.LowerPlural, data.LowerName, data.LowerPlural, data.LowerName)

	stmt := parseStmt(routeCode)
	if stmt != nil {
		// Insert before the last statement of the target block (usually 'return engine')
		if len(targetBlock.List) > 0 {
			targetBlock.List = append(targetBlock.List[:len(targetBlock.List)-1], stmt, targetBlock.List[len(targetBlock.List)-1])
		} else {
			targetBlock.List = append(targetBlock.List, stmt)
		}
	}
}

// InjectListenerRegistration registers a listener for an event in registry.go
func InjectListenerRegistration(registryPath string, eventName string, listenerName string) error {
	// Validate event name and listener name to prevent code injection
	if err := validation.ValidateEventName(eventName); err != nil {
		return fmt.Errorf("invalid event name: %w", err)
	}
	if err := validation.ValidateName(listenerName); err != nil {
		return fmt.Errorf("invalid listener name: %w", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, registryPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse registry.go: %w", err)
	}

	// 1. Add imports
	addImport(node, "github.com/kodia-studio/kodia/internal/core/listeners")

	// 2. Find Registry variable
	var registryExpr *ast.CompositeLit
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, name := range vs.Names {
				if name.Name == "Registry" && i < len(vs.Values) {
					if cl, ok := vs.Values[i].(*ast.CompositeLit); ok {
						registryExpr = cl
						break
					}
				}
			}
		}
	}

	if registryExpr == nil {
		return fmt.Errorf("could not find Registry variable in %s", registryPath)
	}

	// 3. Find or create the event entry
	newListenerExpr, _ := parser.ParseExpr(fmt.Sprintf("&listeners.%sListener{}", listenerName))
	
	found := false
	for _, elt := range registryExpr.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		
		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Value != fmt.Sprintf("\"%s\"", eventName) {
			continue
		}
		
		// Found the event, add listener to slice
		valCl, ok := kv.Value.(*ast.CompositeLit)
		if !ok {
			continue
		}
		
		// Check for duplicate
		isDuplicate := false
		for _, v := range valCl.Elts {
			if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", newListenerExpr) {
				isDuplicate = true
				break
			}
		}
		
		if !isDuplicate {
			valCl.Elts = append(valCl.Elts, newListenerExpr)
		}
		found = true
		break
	}

	if !found {
		// Create new entry
		registryExpr.Elts = append(registryExpr.Elts, &ast.KeyValueExpr{
			Key: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", eventName)},
			Value: &ast.CompositeLit{
				Type: &ast.ArrayType{Elt: &ast.SelectorExpr{
					X:   ast.NewIdent("ports"),
					Sel: ast.NewIdent("Listener"),
				}},
				Elts: []ast.Expr{newListenerExpr},
			},
		})
	}

	return saveFile(registryPath, fset, node)
}

// InjectSeederRegistration registers a new seeder in registry.go
func InjectSeederRegistration(registryPath string, name string) error {
	// Validate seeder name to prevent code injection
	if err := validation.ValidateName(name); err != nil {
		return fmt.Errorf("invalid seeder name: %w", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, registryPath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse registry.go: %w", err)
	}

	// 1. Find Registry variable (which is a slice in this case)
	var registryExpr *ast.CompositeLit
	for _, decl := range node.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, varName := range vs.Names {
				if varName.Name == "Registry" && i < len(vs.Values) {
					if cl, ok := vs.Values[i].(*ast.CompositeLit); ok {
						registryExpr = cl
						break
					}
				}
			}
		}
	}

	if registryExpr == nil {
		return fmt.Errorf("could not find Registry variable in %s", registryPath)
	}

	// 2. Add seeder instance to the slice
	newSeederExpr, _ := parser.ParseExpr(fmt.Sprintf("&%sSeeder{}", name))

	// Duplicate check
	found := false
	for _, elt := range registryExpr.Elts {
		if fmt.Sprintf("%v", elt) == fmt.Sprintf("%v", newSeederExpr) {
			found = true
			break
		}
	}

	if !found {
		registryExpr.Elts = append(registryExpr.Elts, newSeederExpr)
	}

	return saveFile(registryPath, fset, node)
}


func saveFile(path string, fset *token.FileSet, node *ast.File) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := format.Node(f, fset, node); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Run go fmt for final polish
	exec.Command("go", "fmt", path).Run()

	return nil
}
