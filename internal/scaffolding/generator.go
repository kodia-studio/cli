package scaffolding

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
)

// TemplateData holds the variables passed into the `.tmpl` files
type TemplateData struct {
	Name        string // e.g., "Product"
	LowerName   string // e.g., "product"
	Plural      string // e.g., "Products"
	LowerPlural string // e.g., "products"
	Timestamp   string // e.g., "20231024150405"
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

	// Read template content dynamically. In a real-world CLI, 
	// we would bundle templates using `go:embed`. For now, we read from disk.
	// We'll read from `internal/scaffolding/templates/`
	
	// Ensure we read from the current CLI root
	pwd, _ := os.Getwd()
	// Detect if we are in the 'kodia-cli' folder or the root folder
	tmplPath := filepath.Join("kodia-cli", "internal", "scaffolding", "templates", templatePath)
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		// Fallback to local
		tmplPath = filepath.Join(pwd, "internal", "scaffolding", "templates", templatePath)
	}

	tmplContent, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	t, err := template.New("scaffold").Parse(string(tmplContent))
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
func BuildData(name string) TemplateData {
	lowerName := strings.ToLower(name)
	// Simple pluralization (Not perfect, but works for basic cases)
	plural := name + "s"
	lowerPlural := lowerName + "s"

	if strings.HasSuffix(name, "y") {
		plural = strings.TrimSuffix(name, "y") + "ies"
		lowerPlural = strings.TrimSuffix(lowerName, "y") + "ies"
	}

	return TemplateData{
		Name:        name,
		LowerName:   lowerName,
		Plural:      plural,
		LowerPlural: lowerPlural,
		Timestamp:   time.Now().Format("20060102150405"),
	}
}
