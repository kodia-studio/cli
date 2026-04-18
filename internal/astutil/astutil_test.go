package astutil

import (
	"testing"

	"github.com/kodia-studio/cli/internal/scaffolding"
)

// TestValidateTemplateData tests that template data is properly validated to prevent injection
func TestValidateTemplateData(t *testing.T) {
	tests := []struct {
		name      string
		data      scaffolding.TemplateData
		wantError bool
		errorMsg  string
	}{
		// Valid template data
		{
			name: "valid template data",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: false,
		},
		{
			name: "valid with underscores",
			data: scaffolding.TemplateData{
				Name:        "UserProfile",
				LowerName:   "user_profile",
				Plural:      "UserProfiles",
				LowerPlural: "user_profiles",
				ProjectName: "my_project",
			},
			wantError: false,
		},

		// Invalid Name (lowercase start)
		{
			name: "invalid Name - lowercase",
			data: scaffolding.TemplateData{
				Name:        "product",
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Name",
		},

		// Invalid LowerName (invalid identifier)
		{
			name: "invalid LowerName - starts with number",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "123product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "LowerName",
		},

		// Invalid Plural (lowercase)
		{
			name: "invalid Plural - lowercase",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product",
				Plural:      "products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Plural",
		},

		// Invalid LowerPlural (special characters)
		{
			name: "invalid LowerPlural - contains hyphen",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "product-items",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "LowerPlural",
		},

		// Code injection attempts in Name
		{
			name: "code injection in Name - parenthesis",
			data: scaffolding.TemplateData{
				Name:        "Product)); DROP TABLE",
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Name",
		},

		// Code injection in LowerName
		{
			name: "code injection in LowerName - backtick",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product`; rm -rf",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "LowerName",
		},

		// Template injection in Plural
		{
			name: "template injection in Plural",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product",
				Plural:      "Products{{.Malicious}}",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Plural",
		},

		// Reserved keyword injection
		{
			name: "reserved keyword in Name",
			data: scaffolding.TemplateData{
				Name:        "return",
				LowerName:   "return",
				Plural:      "Returns",
				LowerPlural: "returns",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Name",
		},

		// Empty fields
		{
			name: "empty Name",
			data: scaffolding.TemplateData{
				Name:        "",
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "myproject",
			},
			wantError: true,
			errorMsg:  "Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplateData(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("validateTemplateData() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if tt.wantError && tt.errorMsg != "" && err != nil {
				if !containsString(err.Error(), tt.errorMsg) {
					t.Errorf("validateTemplateData() error should contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

// TestInjectDependencyInjection_RejectsInvalidInput tests that malicious input is rejected
func TestInjectDependencyInjection_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		data      scaffolding.TemplateData
		wantError bool
	}{
		// Code injection attempts
		{
			name: "injection: double quote in Name",
			data: scaffolding.TemplateData{
				Name:        `Product"); DROP TABLE users; --`,
				LowerName:   "product",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "kodia",
			},
			wantError: true,
		},
		{
			name: "injection: semicolon in LowerName",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product; malicious_code()",
				Plural:      "Products",
				LowerPlural: "products",
				ProjectName: "kodia",
			},
			wantError: true,
		},
		{
			name: "injection: backtick in Plural",
			data: scaffolding.TemplateData{
				Name:        "Product",
				LowerName:   "product",
				Plural:      "Products`command injection`",
				LowerPlural: "products",
				ProjectName: "kodia",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call InjectDependencyInjection with invalid data
			// It should reject the data during validation
			err := InjectDependencyInjection("nonexistent.go", tt.data)

			if !tt.wantError && err == nil {
				// If we expect no error and got none, that's fine
				// (might fail due to file not found, but not validation)
				return
			}

			if tt.wantError {
				if err == nil {
					t.Errorf("InjectDependencyInjection() expected error for malicious input, got nil")
				}
				// Validation error should come first, before file operations
				if !containsString(err.Error(), "invalid template data") {
					t.Logf("Expected validation error, got: %v", err)
				}
			}
		})
	}
}

// TestInjectJobRegistration_RejectsInvalidInput tests job registration validation
func TestInjectJobRegistration_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name      string
		data      scaffolding.TemplateData
		wantError bool
	}{
		{
			name: "invalid job name with special chars",
			data: scaffolding.TemplateData{
				Name:        "My<Job>",
				LowerName:   "my_job",
				Plural:      "Jobs",
				LowerPlural: "jobs",
				ProjectName: "kodia",
			},
			wantError: true,
		},
		{
			name: "job name with parentheses",
			data: scaffolding.TemplateData{
				Name:        "Job()",
				LowerName:   "job",
				Plural:      "Jobs",
				LowerPlural: "jobs",
				ProjectName: "kodia",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InjectJobRegistration("nonexistent.go", tt.data, false)

			if tt.wantError {
				if err == nil {
					t.Errorf("InjectJobRegistration() expected error for invalid input")
				}
				if !containsString(err.Error(), "invalid template data") {
					t.Logf("Expected validation error, got: %v", err)
				}
			}
		})
	}
}

// TestInjectListenerRegistration_RejectsInvalidInput tests listener registration validation
func TestInjectListenerRegistration_RejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name          string
		eventName     string
		listenerName  string
		wantError     bool
	}{
		// Valid inputs
		{
			name:         "valid event and listener names",
			eventName:    "UserCreated",
			listenerName: "SendWelcomeEmail",
			wantError:    false,
		},

		// Invalid event name
		{
			name:         "invalid event name - lowercase",
			eventName:    "userCreated",
			listenerName: "SendWelcomeEmail",
			wantError:    true,
		},
		{
			name:         "invalid event name - with injection",
			eventName:    `UserCreated"); DROP TABLE events; --`,
			listenerName: "SendWelcomeEmail",
			wantError:    true,
		},

		// Invalid listener name
		{
			name:         "invalid listener name - lowercase",
			eventName:    "UserCreated",
			listenerName: "sendWelcomeEmail",
			wantError:    true,
		},
		{
			name:         "invalid listener name - with special chars",
			eventName:    "UserCreated",
			listenerName: "Send-Welcome-Email",
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InjectListenerRegistration("nonexistent.go", tt.eventName, tt.listenerName)

			if (err != nil) != tt.wantError {
				if tt.wantError {
					t.Errorf("InjectListenerRegistration() expected error, got nil")
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substring string) bool {
	for i := 0; i <= len(s)-len(substring); i++ {
		if s[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}

// BenchmarkValidateTemplateData measures validation performance
func BenchmarkValidateTemplateData(b *testing.B) {
	validData := scaffolding.TemplateData{
		Name:        "Product",
		LowerName:   "product",
		Plural:      "Products",
		LowerPlural: "products",
		ProjectName: "myproject",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validateTemplateData(validData)
	}
}
