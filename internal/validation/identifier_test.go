package validation

import (
	"strings"
	"testing"
)

func TestValidateIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		identifier string
		wantError bool
	}{
		// Valid identifiers
		{
			name:      "simple letter",
			identifier: "a",
			wantError: false,
		},
		{
			name:      "multiple letters",
			identifier: "abc",
			wantError: false,
		},
		{
			name:      "with numbers",
			identifier: "var123",
			wantError: false,
		},
		{
			name:      "with underscores",
			identifier: "my_var",
			wantError: false,
		},
		{
			name:      "start with underscore",
			identifier: "_private",
			wantError: false,
		},
		{
			name:      "multiple underscores",
			identifier: "_var_name_",
			wantError: false,
		},
		{
			name:      "mixed case",
			identifier: "myVarName",
			wantError: false,
		},
		{
			name:      "single uppercase",
			identifier: "A",
			wantError: false,
		},
		{
			name:      "Pascal case",
			identifier: "MyStructName",
			wantError: false,
		},

		// Invalid identifiers
		{
			name:      "empty string",
			identifier: "",
			wantError: true,
		},
		{
			name:      "starts with number",
			identifier: "123var",
			wantError: true,
		},
		{
			name:      "starts with hyphen",
			identifier: "-var",
			wantError: true,
		},
		{
			name:      "contains space",
			identifier: "var name",
			wantError: true,
		},
		{
			name:      "contains special char dash",
			identifier: "var-name",
			wantError: true,
		},
		{
			name:      "contains special char dot",
			identifier: "var.name",
			wantError: true,
		},
		{
			name:      "contains special char colon",
			identifier: "var:name",
			wantError: true,
		},
		{
			name:      "contains special char semicolon",
			identifier: "var;name",
			wantError: true,
		},
		{
			name:      "contains parentheses",
			identifier: "var(name)",
			wantError: true,
		},
		{
			name:      "contains curly braces",
			identifier: "var{name}",
			wantError: true,
		},

		// Go reserved keywords
		{
			name:      "keyword break",
			identifier: "break",
			wantError: true,
		},
		{
			name:      "keyword func",
			identifier: "func",
			wantError: true,
		},
		{
			name:      "keyword interface",
			identifier: "interface",
			wantError: true,
		},
		{
			name:      "keyword return",
			identifier: "return",
			wantError: true,
		},

		// Length limits
		{
			name:      "max length valid (255 chars)",
			identifier: "a" + strings.Repeat("b", 254),
			wantError: false,
		},
		{
			name:      "exceeds max length (256 chars)",
			identifier: "a" + strings.Repeat("b", 255),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIdentifier(tt.identifier)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateIdentifier() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateEventName(t *testing.T) {
	tests := []struct {
		name      string
		eventName string
		wantError bool
	}{
		// Valid event names (PascalCase)
		{
			name:      "simple event",
			eventName: "User",
			wantError: false,
		},
		{
			name:      "multi-word event",
			eventName: "UserCreated",
			wantError: false,
		},
		{
			name:      "complex event",
			eventName: "OrderPaymentProcessed",
			wantError: false,
		},

		// Invalid event names
		{
			name:      "lowercase start",
			eventName: "userCreated",
			wantError: true,
		},
		{
			name:      "all lowercase",
			eventName: "user",
			wantError: true,
		},
		{
			name:      "all uppercase",
			eventName: "USER",
			wantError: false, // Technically valid Pascal case
		},
		{
			name:      "underscore start",
			eventName: "_User",
			wantError: true,
		},
		{
			name:      "empty",
			eventName: "",
			wantError: true,
		},
		{
			name:      "contains hyphen",
			eventName: "User-Created",
			wantError: true,
		},
		{
			name:      "looks like keyword but capitalized (valid)",
			eventName: "Return",
			wantError: false, // "Return" is valid; "return" is reserved
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEventName(tt.eventName)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateEventName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		// Valid names (PascalCase)
		{
			name:      "simple name",
			input:     "Product",
			wantError: false,
		},
		{
			name:      "multi-word name",
			input:     "ProductService",
			wantError: false,
		},

		// Invalid names
		{
			name:      "lowercase",
			input:     "product",
			wantError: true,
		},
		{
			name:      "empty",
			input:     "",
			wantError: true,
		},
		{
			name:      "with dash",
			input:     "Product-Service",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateName() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestSanitizeIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantError bool
	}{
		// Valid inputs that don't need sanitization
		{
			name:      "valid identifier",
			input:     "myVar",
			want:      "myVar",
			wantError: false,
		},

		// Inputs that need sanitization
		{
			name:      "spaces to underscores",
			input:     "my var",
			want:      "my_var",
			wantError: false,
		},
		{
			name:      "dashes to underscores",
			input:     "my-var",
			want:      "my_var",
			wantError: false,
		},
		{
			name:      "leading numbers removed",
			input:     "123myVar",
			want:      "myVar",
			wantError: false,
		},
		{
			name:      "dots to underscores",
			input:     "my.var.name",
			want:      "my_var_name",
			wantError: false,
		},
		{
			name:      "mixed special chars",
			input:     "my-var.name_123",
			want:      "my_var_name_123",
			wantError: false,
		},

		// Invalid inputs that cannot be sanitized
		{
			name:      "only numbers",
			input:     "123",
			wantError: true,
		},
		{
			name:      "only special chars",
			input:     "---",
			wantError: true,
		},
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeIdentifier(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("SanitizeIdentifier() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("SanitizeIdentifier() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkValidateIdentifier(b *testing.B) {
	identifiers := []string{
		"myVar",
		"MyStructName",
		"_privateVar",
		"valid_identifier_123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, id := range identifiers {
			ValidateIdentifier(id)
		}
	}
}

func BenchmarkValidateEventName(b *testing.B) {
	eventNames := []string{
		"UserCreated",
		"OrderPaymentProcessed",
		"DataSynced",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, name := range eventNames {
			ValidateEventName(name)
		}
	}
}

func BenchmarkSanitizeIdentifier(b *testing.B) {
	inputs := []string{
		"my-var",
		"my var",
		"my.var.name",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			SanitizeIdentifier(input)
		}
	}
}
