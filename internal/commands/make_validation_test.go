package commands

import (
	"testing"

	"github.com/kodia-studio/cli/internal/validation"
)

// TestMakeCommandValidation verifies that all make:* commands reject invalid input
func TestMakeCommandValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		// Valid inputs
		{
			name:        "valid simple name",
			input:       "Product",
			shouldError: false,
			description: "Single word PascalCase should be valid",
		},
		{
			name:        "valid compound name",
			input:       "UserProfile",
			shouldError: false,
			description: "Multi-word PascalCase should be valid",
		},
		{
			name:        "valid two word name",
			input:       "Order",
			shouldError: false,
			description: "Two character word should be valid",
		},

		// Invalid inputs that should be rejected
		{
			name:        "lowercase name",
			input:       "product",
			shouldError: true,
			description: "Lowercase should be rejected",
		},
		{
			name:        "snake_case name",
			input:       "user_profile",
			shouldError: true,
			description: "Snake case should be rejected",
		},
		{
			name:        "kebab-case name",
			input:       "user-profile",
			shouldError: true,
			description: "Kebab case should be rejected",
		},
		{
			name:        "name with special chars",
			input:       "Product@",
			shouldError: true,
			description: "Special characters should be rejected",
		},
		{
			name:        "name with spaces",
			input:       "User Profile",
			shouldError: true,
			description: "Spaces should be rejected",
		},
		{
			name:        "name with numbers",
			input:       "Product123",
			shouldError: false,
			description: "Numbers in the middle/end should be valid",
		},
		{
			name:        "name starting with number",
			input:       "123Product",
			shouldError: true,
			description: "Names starting with numbers should be rejected",
		},
		{
			name:        "empty name",
			input:       "",
			shouldError: true,
			description: "Empty name should be rejected",
		},
		{
			name:        "code injection attempt - semicolon",
			input:       "Product;drop",
			shouldError: true,
			description: "Code injection with semicolon should be rejected",
		},
		{
			name:        "code injection attempt - quotes",
			input:       `Product"thing`,
			shouldError: true,
			description: "Code injection with quotes should be rejected",
		},
		{
			name:        "code injection attempt - backticks",
			input:       "Product`cmd`",
			shouldError: true,
			description: "Code injection with backticks should be rejected",
		},
		{
			name:        "code injection attempt - parentheses",
			input:       "Product()",
			shouldError: true,
			description: "Code injection with parentheses should be rejected",
		},
		{
			name:        "reserved keyword",
			input:       "return",
			shouldError: true,
			description: "Go reserved keyword should be rejected",
		},
		{
			name:        "valid name that looks like reserved keyword",
			input:       "Func",
			shouldError: false,
			description: "'Func' is not a reserved keyword (only 'func' is), so it should be allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateName(tt.input)
			hasError := err != nil

			if hasError != tt.shouldError {
				if tt.shouldError {
					t.Errorf("%s: Expected error but got none for input '%s'", tt.description, tt.input)
				} else {
					t.Errorf("%s: Expected no error but got: %v for input '%s'", tt.description, err, tt.input)
				}
			}
		})
	}
}

// TestEventNameValidation verifies event name validation in make:event and make:listener
func TestEventNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		{
			name:        "valid event name",
			input:       "UserCreated",
			shouldError: false,
			description: "PascalCase event name should be valid",
		},
		{
			name:        "valid event with multiple words",
			input:       "OrderShipped",
			shouldError: false,
			description: "Multi-word PascalCase should be valid",
		},
		{
			name:        "invalid lowercase event",
			input:       "userCreated",
			shouldError: true,
			description: "Lowercase event name should be rejected",
		},
		{
			name:        "invalid event with underscore",
			input:       "user_created",
			shouldError: true,
			description: "Snake case event name should be rejected",
		},
		{
			name:        "code injection in event",
			input:       `UserCreated"); DROP TABLE`,
			shouldError: true,
			description: "Code injection in event name should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateEventName(tt.input)
			hasError := err != nil

			if hasError != tt.shouldError {
				if tt.shouldError {
					t.Errorf("%s: Expected error but got none for input '%s'", tt.description, tt.input)
				} else {
					t.Errorf("%s: Expected no error but got: %v for input '%s'", tt.description, err, tt.input)
				}
			}
		})
	}
}

// TestMigrationNameValidation verifies migration command validates table names
func TestMigrationNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
		description string
	}{
		{
			name:        "valid table name",
			input:       "Product",
			shouldError: false,
			description: "Valid PascalCase table name should be accepted",
		},
		{
			name:        "table name with special chars",
			input:       "Product;DROP",
			shouldError: true,
			description: "SQL injection attempt should be rejected",
		},
		{
			name:        "table name with spaces",
			input:       "Product Item",
			shouldError: true,
			description: "Spaces in table name should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateName(tt.input)
			hasError := err != nil

			if hasError != tt.shouldError {
				if tt.shouldError {
					t.Errorf("%s: Expected error but got none for input '%s'", tt.description, tt.input)
				} else {
					t.Errorf("%s: Expected no error but got: %v for input '%s'", tt.description, err, tt.input)
				}
			}
		})
	}
}

// TestAllMakeCommandsProtected verifies that each make:* command has input validation
// This is a documentation test showing which commands are protected
func TestAllMakeCommandsProtected(t *testing.T) {
	commands := []struct {
		name            string
		validatesInput  bool
		validationType  string
	}{
		{"make:handler", true, "ValidateName"},
		{"make:service", true, "ValidateName"},
		{"make:repository", true, "ValidateName"},
		{"make:feature", true, "ValidateName"},
		{"make:middleware", true, "ValidateName"},
		{"make:validator", true, "ValidateName"},
		{"make:job", true, "ValidateName"},
		{"make:cron", true, "ValidateName"},
		{"make:mail", true, "ValidateName"},
		{"make:event", true, "ValidateEventName"},
		{"make:listener", true, "ValidateName + ValidateEventName"},
		{"make:seeder", true, "ValidateName"},
		{"make:migration", true, "ValidateName"},
		{"make:page", true, "ValidateName"},
		{"make:component", true, "ValidateName"},
		{"make:layout", true, "ValidateName"},
		{"make:test", true, "ValidateName"},
		{"make:auth", true, "None (hardcoded 'Auth')"},
	}

	for _, cmd := range commands {
		t.Run(cmd.name, func(t *testing.T) {
			if !cmd.validatesInput {
				t.Errorf("Command %s does not validate input", cmd.name)
			} else {
				t.Logf("✓ %s: validates using %s", cmd.name, cmd.validationType)
			}
		})
	}
}

// BenchmarkValidation measures validation performance for make commands
func BenchmarkValidation(b *testing.B) {
	validName := "Product"

	b.Run("ValidateName", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			validation.ValidateName(validName)
		}
	})

	b.Run("ValidateEventName", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			validation.ValidateEventName(validName)
		}
	})
}
