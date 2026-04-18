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

// TemplateData holds the variables passed into the `.tmpl` files
type TemplateData struct {
	Name        string // e.g., "Product"
	LowerName   string // e.g., "product"
	Plural      string // e.g., "Products"
	LowerPlural string // e.g., "products"
	Timestamp   string // e.g., "20231024150405"
	ProjectName string // e.g., "kodia-framework"
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
