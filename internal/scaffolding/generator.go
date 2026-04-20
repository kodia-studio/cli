package scaffolding

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
)

// Embed all template files into the binary at compile time
//go:embed templates/*.tmpl
var templateFS embed.FS

// Field represents a property in a scaffolded entity
type Field struct {
	Name        string // e.g., "Name"
	LowerName   string // e.g., "name"
	Type        string // e.g., "string", "references"
	GoType      string // e.g., "string", "uint64"
	SQLType     string // e.g., "VARCHAR(255)", "BIGINT"
	Validation  string // e.g., "required", "email"
	IsRefs      bool   // true if references another model
	RefModel    string // e.g., "User"
	RefTable    string // e.g., "users"
	IsUnique    bool
	IsNullable  bool
}

// TemplateData holds the variables passed into the `.tmpl` files
type TemplateData struct {
	Name        string  // e.g., "Product"
	LowerName   string  // e.g., "product"
	Plural      string  // e.g., "Products"
	LowerPlural string  // e.g., "products"
	Timestamp   string  // e.g., "20231024150405"
	ProjectName string  // e.g., "kodia-framework"
	Fields      []Field // List of fields for the entity
}

// Generate processes a template file and writes it to the destination
func Generate(templatePath, destPath string, data TemplateData) error {
	// Create destination directory if it doesn't exist
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Check if file already exists
	if _, err := os.Stat(destPath); err == nil {
		color.Yellow("⚠️  Skipped: %s (File already exists)", destPath)
		return nil
	}

	// Read template content from embedded filesystem
	tmplFile := filepath.Join("templates", templatePath)
	tmplContent, err := templateFS.ReadFile(tmplFile)
	if err != nil {
		return fmt.Errorf("could not find template %s: %w", templatePath, err)
	}

	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"contains": strings.Contains,
		"lower":    strings.ToLower,
	}

	t, err := template.New("scaffold").Funcs(funcMap).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := os.WriteFile(destPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	color.Green("✅ Created: %s", destPath)
	return nil
}

// BuildData constructs the strings needed for templates
func BuildData(name string, fieldsRaw string) TemplateData {
	lowerName := strings.ToLower(name)
	// Simple pluralization
	plural := name + "s"
	lowerPlural := lowerName + "s"

	if strings.HasSuffix(name, "y") {
		plural = strings.TrimSuffix(name, "y") + "ies"
		lowerPlural = strings.TrimSuffix(lowerName, "y") + "ies"
	}

	fields := ParseFields(fieldsRaw)

	return TemplateData{
		Name:        name,
		LowerName:   lowerName,
		Plural:      plural,
		LowerPlural: lowerPlural,
		Timestamp:   time.Now().Format("20060102150405"),
		Fields:      fields,
	}
}

// ParseFields converts raw field strings (name:string,author:references:User) into Field structs
func ParseFields(raw string) []Field {
	if raw == "" {
		return []Field{}
	}

	parts := strings.Split(raw, ",")
	fields := make([]Field, 0, len(parts))

	for _, part := range parts {
		fieldParts := strings.Split(part, ":")
		if len(fieldParts) < 2 {
			continue
		}

		name := fieldParts[0]
		fieldType := fieldParts[1]
		
		field := Field{
			Name:      strings.Title(name),
			LowerName: strings.ToLower(name),
			Type:      fieldType,
		}

		// Handle data types and mapping
		switch fieldType {
		case "string":
			field.GoType = "string"
			field.SQLType = "VARCHAR(255)"
		case "text":
			field.GoType = "string"
			field.SQLType = "TEXT"
		case "integer", "int":
			field.GoType = "int"
			field.SQLType = "INTEGER"
		case "bigint":
			field.GoType = "int64"
			field.SQLType = "BIGINT"
		case "float", "decimal":
			field.GoType = "float64"
			field.SQLType = "DECIMAL(10,2)"
		case "boolean", "bool":
			field.GoType = "bool"
			field.SQLType = "BOOLEAN"
		case "timestamp", "datetime":
			field.GoType = "time.Time"
			field.SQLType = "TIMESTAMP"
		case "references":
			if len(fieldParts) >= 3 {
				model := fieldParts[2]
				field.IsRefs = true
				field.RefModel = model
				field.RefTable = strings.ToLower(model) + "s"
				field.GoType = "uint64"
				field.SQLType = "BIGINT"
				// Standard foreign key naming: author -> author_id
				if !strings.HasSuffix(field.LowerName, "_id") {
					field.Name = field.Name + "ID"
					field.LowerName = field.LowerName + "_id"
				}
			}
		case "enum":
			field.GoType = "string"
			field.SQLType = "VARCHAR(50)"
		default:
			field.GoType = "string"
			field.SQLType = "VARCHAR(255)"
		}

		// Check for additional options like unique, nullable
		for i := 2; i < len(fieldParts); i++ {
			opt := strings.ToLower(fieldParts[i])
			if opt == "unique" {
				field.IsUnique = true
			} else if opt == "nullable" {
				field.IsNullable = true
			} else if !field.IsRefs {
				// Assume it's a validation tag if not references
				if field.Validation == "" {
					field.Validation = opt
				} else {
					field.Validation += "," + opt
				}
			}
		}
		
		// Default validation for common fields
		if field.LowerName == "email" && !strings.Contains(field.Validation, "email") {
			if field.Validation == "" {
				field.Validation = "email"
			} else {
				field.Validation += ",email"
			}
		}

		fields = append(fields, field)
	}

	return fields
}
