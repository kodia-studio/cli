package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidateIdentifier checks if the given string is a valid Go identifier.
// Valid identifiers must:
// - Start with a letter or underscore
// - Contain only letters, digits, and underscores
// - Not be a Go reserved keyword
// - Not exceed 255 characters
func ValidateIdentifier(identifier string) error {
	if identifier == "" {
		return fmt.Errorf("identifier cannot be empty")
	}

	if len(identifier) > 255 {
		return fmt.Errorf("identifier is too long (max 255 characters, got %d)", len(identifier))
	}

	// Check first character is letter or underscore
	r := rune(identifier[0])
	if !unicode.IsLetter(r) && r != '_' {
		return fmt.Errorf("identifier must start with a letter or underscore, got '%c'", r)
	}

	// Check remaining characters are letters, digits, or underscores
	for i, r := range identifier {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("identifier contains invalid character '%c' at position %d", r, i)
		}
	}

	// Check if it's a Go reserved keyword
	if isReservedKeyword(identifier) {
		return fmt.Errorf("'%s' is a reserved Go keyword and cannot be used", identifier)
	}

	return nil
}

// ValidateEventName checks if the given string is a valid event name.
// Event names follow the same rules as Go identifiers, but are more restrictive:
// - Must be PascalCase (start with uppercase letter)
// - Used for type/struct names in the codebase
func ValidateEventName(eventName string) error {
	// First validate as identifier
	if err := ValidateIdentifier(eventName); err != nil {
		return err
	}

	// Event names should be PascalCase (start with uppercase)
	if !unicode.IsUpper(rune(eventName[0])) {
		return fmt.Errorf("event name must start with an uppercase letter (PascalCase), got '%s'", eventName)
	}

	return nil
}

// ValidateName checks if the given string is a valid name for code generation.
// Names must be valid Go identifiers and follow PascalCase convention.
func ValidateName(name string) error {
	// First validate as identifier
	if err := ValidateIdentifier(name); err != nil {
		return err
	}

	// Names should be PascalCase (start with uppercase)
	if !unicode.IsUpper(rune(name[0])) {
		return fmt.Errorf("name must start with an uppercase letter (PascalCase), got '%s'", name)
	}

	return nil
}

// ValidateRoute checks if the given string is a valid SvelteKit route name.
// Routes can be lowercase and follow kebab-case or underscore patterns.
func ValidateRoute(route string) error {
	if route == "" {
		return fmt.Errorf("route cannot be empty")
	}

	// Basic validation for route characters
	for i, r := range route {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' && r != '(' && r != ')' {
			return fmt.Errorf("route contains invalid character '%c' at position %d", r, i)
		}
	}

	return nil
}

// SanitizeIdentifier removes invalid characters from an identifier and returns a valid one.
// If the result would be invalid, returns an error.
func SanitizeIdentifier(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input cannot be empty")
	}

	// Remove leading/trailing whitespace
	input = strings.TrimSpace(input)

	// Replace invalid characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	sanitized := re.ReplaceAllString(input, "_")

	// Ensure it starts with letter or underscore (remove leading digits)
	i := 0
	for i < len(sanitized) && (unicode.IsDigit(rune(sanitized[i])) || sanitized[i] == '_') {
		i++
	}
	if i > 0 {
		sanitized = sanitized[i:]
	}

	// If we removed everything, return error
	if sanitized == "" {
		return "", fmt.Errorf("input '%s' cannot be converted to a valid identifier", input)
	}

	// Validate the sanitized result
	if err := ValidateIdentifier(sanitized); err != nil {
		return "", fmt.Errorf("sanitized identifier '%s' is invalid: %w", sanitized, err)
	}

	return sanitized, nil
}

// Go reserved keywords that cannot be used as identifiers
var reservedKeywords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

func isReservedKeyword(identifier string) bool {
	return reservedKeywords[identifier]
}
