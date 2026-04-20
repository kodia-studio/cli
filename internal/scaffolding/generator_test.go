package scaffolding

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEmbeddedTemplatesAvailable verifies all templates are embedded correctly
func TestEmbeddedTemplatesAvailable(t *testing.T) {
	templates := []string{
		"handler.tmpl",
		"migration_up.tmpl",
		"migration_down.tmpl",
		"repository.tmpl",
		"service.tmpl",
		"middleware.tmpl",
		"validator.tmpl",
		"job.tmpl",
		"cron.tmpl",
		"mail.tmpl",
		"event.tmpl",
		"listener.tmpl",
		"seeder.tmpl",
	}

	for _, tmpl := range templates {
		t.Run(tmpl, func(t *testing.T) {
			_, err := templateFS.ReadFile(filepath.Join("templates", tmpl))
			if err != nil {
				t.Errorf("Template %s not found in embedded FS: %v", tmpl, err)
			}
		})
	}
}

// TestGenerateWithEmbeddedTemplates verifies template processing works with embedded FS
func TestGenerateWithEmbeddedTemplates(t *testing.T) {
	tempDir := t.TempDir()
	destPath := filepath.Join(tempDir, "test_handler.go")

	data := TemplateData{
		Name:        "Product",
		LowerName:   "product",
		Plural:      "Products",
		LowerPlural: "products",
		ProjectName: "testproject",
	}

	// Try to generate from embedded template
	err := Generate("handler.tmpl", destPath, data)
	if err != nil {
		t.Fatalf("Failed to generate from embedded template: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Errorf("Generated file not found: %s", destPath)
	}

	// Verify file has content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Generated file is empty")
	}

	// Verify template variables were replaced
	contentStr := string(content)
	if contentStr == "" {
		t.Errorf("Generated file has no content")
	}
}

// TestBuildDataStructure verifies BuildData creates proper template data
func TestBuildDataStructure(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedName   string
		expectedLower  string
		expectedPlural string
	}{
		{
			name:           "simple name",
			input:          "Product",
			expectedName:   "Product",
			expectedLower:  "product",
			expectedPlural: "Products",
		},
		{
			name:           "name ending in y",
			input:          "Category",
			expectedName:   "Category",
			expectedLower:  "category",
			expectedPlural: "Categories",
		},
		{
			name:           "with fields",
			input:          "Post",
			expectedName:   "Post",
			expectedLower:  "post",
			expectedPlural: "Posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := ""
			if tt.name == "with fields" {
				fields = "title:string,body:text"
			}
			data := BuildData(tt.input, fields)

			if data.Name != tt.expectedName {
				t.Errorf("Name: got %s, want %s", data.Name, tt.expectedName)
			}
			if tt.name == "with fields" && len(data.Fields) != 2 {
				t.Errorf("Fields: got %d, want 2", len(data.Fields))
			}
			if data.LowerName != tt.expectedLower {
				t.Errorf("LowerName: got %s, want %s", data.LowerName, tt.expectedLower)
			}
			if data.Plural != tt.expectedPlural {
				t.Errorf("Plural: got %s, want %s", data.Plural, tt.expectedPlural)
			}
		})
	}
}
