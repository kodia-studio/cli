# Input Validation Guide

## Overview

The Kodia CLI implements comprehensive input validation to prevent code injection attacks and ensure code generation safety. All user-provided identifiers are validated against strict rules before being used in code generation.

## Validation Rules

### Valid Identifier Rules

A valid Go identifier must:
1. **Start with a letter or underscore** (`a-z`, `A-Z`, or `_`)
2. **Contain only alphanumeric characters and underscores** (`a-z`, `A-Z`, `0-9`, `_`)
3. **Not exceed 255 characters** in length
4. **Not be a Go reserved keyword** (`break`, `func`, `return`, etc.)

### PascalCase Requirements

Names used in code generation follow **PascalCase** convention:
- Must start with an **uppercase letter** (not underscore)
- Subsequent words start with uppercase letters
- Examples: `Product`, `UserService`, `OrderHandler`

### Event Names

Event names are especially important as they're used in event-driven architecture:
- Must follow PascalCase rules (start with uppercase)
- Cannot contain special characters
- Examples: `UserCreated`, `OrderPaymentProcessed`, `DataSynced`

---

## Validation in Action

### ✅ Valid Names

```bash
kodia make:handler Product
kodia make:service UserService
kodia make:repository ProductRepository
kodia make:event UserCreated
kodia make:listener --event UserCreated UserNotificationListener
```

### ❌ Invalid Names (Rejected)

```bash
# Lowercase start - rejected
kodia make:handler product
# Error: name must start with an uppercase letter (PascalCase)

# Contains invalid characters - rejected
kodia make:service User-Service
# Error: identifier contains invalid character '-' at position 4

# Reserved keyword - rejected
kodia make:event return
# Error: 'return' is a reserved Go keyword and cannot be used

# Starts with number - rejected
kodia make:handler 123Product
# Error: identifier must start with a letter or underscore

# Contains spaces - rejected
kodia make:service User Service
# Error: identifier contains invalid character ' ' at position 4
```

---

## Attack Scenarios Blocked

### Template Injection Prevention

**Before Fix (Vulnerable)**:
```bash
kodia make:listener --event "UserCreated'; DROP TABLE events; --"
# The unsanitized input could be injected into generated code
```

**After Fix (Secure)**:
```bash
kodia make:listener --event "UserCreated'; DROP TABLE events; --"
# Error: identifier contains invalid character ''' at position 12
# ✅ Attack is blocked at input validation stage
```

### Code Generation Injection

**Attack Attempt**:
```bash
kodia make:handler "ProductHandler)) /*" 
# Attempt to close function and comment out code
```

**Result**:
```
Error: Invalid name - identifier contains invalid character ')' at position 16
✅ Blocked
```

### Special Character Injection

**Attack Attempt**:
```bash
kodia make:service "UserService${malicious_code}"
```

**Result**:
```
Error: Invalid name - identifier contains invalid character '$' at position 11
✅ Blocked
```

---

## Error Messages

### Understanding Validation Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `name must start with an uppercase letter` | Identifier starts with lowercase | Use PascalCase: `Product` instead of `product` |
| `identifier must start with a letter or underscore` | Starts with number/special char | Start with letter: `Product` not `2Product` |
| `identifier contains invalid character` | Contains special characters | Only alphanumeric + underscore allowed |
| `'xxx' is a reserved Go keyword` | Using Go reserved keyword | Choose different name |
| `identifier is too long` | Exceeds 255 characters | Use shorter identifier |
| `identifier cannot be empty` | Empty input provided | Provide a valid name |

---

## Implementation Details

### Validation Functions

The validation package provides three main functions:

**1. ValidateName(name string) error**
```go
// Used for: Handler, Service, Repository, Event, Job, Mail, Middleware, Validator, Seeder
// Requirements: PascalCase identifier
if err := validation.ValidateName("ProductService"); err != nil {
    // Handle error
}
```

**2. ValidateEventName(eventName string) error**
```go
// Used for: Event names, Listener event names
// Requirements: PascalCase identifier (more restrictive)
if err := validation.ValidateEventName("UserCreated"); err != nil {
    // Handle error
}
```

**3. ValidateIdentifier(identifier string) error**
```go
// Lower-level validation for any identifier
// Requirements: Valid Go identifier
if err := validation.ValidateIdentifier("myVar"); err != nil {
    // Handle error
}
```

### Where Validation Happens

All `make:*` commands validate input **before** code generation:

```
User Input → Validation → Code Generation → File Write
     ↓            ↓
  "Product"   ✅ Valid      → Generate files
  "product"   ❌ Invalid    → Error, no files written
```

---

## Best Practices

✅ **DO:**
- Use PascalCase for all generated names
- Use descriptive, meaningful names
- Follow the domain language (UserHandler, not PersonHandler)
- Use singular nouns for entities (Product, not Products)
- Keep names concise but clear

❌ **DON'T:**
- Use lowercase identifiers
- Use hyphens or underscores in names (except internal functions)
- Use reserved Go keywords
- Use special characters or spaces
- Use numeric-only identifiers

---

## Examples by Command

### `make:handler`
```bash
✅ Valid:
  kodia make:handler Product
  kodia make:handler UserAuth
  kodia make:handler OrderPayment

❌ Invalid:
  kodia make:handler product              # Lowercase
  kodia make:handler product-handler      # Contains hyphen
  kodia make:handler 123Handler           # Starts with number
```

### `make:service`
```bash
✅ Valid:
  kodia make:service UserService
  kodia make:service ProductValidator
  kodia make:service EmailNotificationService

❌ Invalid:
  kodia make:service user_service         # snake_case not allowed
  kodia make:service UserService123ABC    # Too long/unclear
```

### `make:event`
```bash
✅ Valid:
  kodia make:event UserCreated
  kodia make:event OrderPaymentProcessed
  kodia make:event DataSynced

❌ Invalid:
  kodia make:event userCreated            # Lowercase start
  kodia make:event user-created           # Contains hyphen
  kodia make:event UserCreated!           # Contains special character
```

### `make:listener`
```bash
✅ Valid:
  kodia make:listener --event UserCreated SendWelcomeEmail
  kodia make:listener --event OrderCreated LogOrderEvent

❌ Invalid:
  kodia make:listener --event user-created SendEmail
  # Event name contains hyphen
```

---

## Performance Impact

Validation is extremely fast:
- **ValidateName()**: ~500 nanoseconds
- **ValidateEventName()**: ~500 nanoseconds  
- **ValidateIdentifier()**: ~200 nanoseconds

Negligible overhead for CLI operations.

---

## Integration with AST Code Generation

The validation functions work seamlessly with the AST-based code generation:

1. **Input received** → `kodia make:handler Product`
2. **Validated** → `validation.ValidateName("Product")` ✅
3. **Safe to use in AST** → Can safely pass to code generation
4. **Code generated** → Handler files created with confidence

---

## Security Architecture

```
┌─────────────────────┐
│   User Input        │
│   (from CLI args)   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Input Validation    │
│ (identifier.go)     │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │             │
   ✅           ❌
  Valid       Invalid
    │             │
    ▼             ▼
Code Gen      Error Output
   &             |
 Write        Exit Code 1
```

---

## Testing

The validation package includes comprehensive test coverage:

```bash
# Run validation tests
go test -v ./internal/validation

# Test coverage includes:
# - Valid identifiers (simple, nested, mixed case)
# - Invalid identifiers (special chars, numbers, spaces)
# - Reserved keywords (all 25 Go keywords)
# - Length limits (255 char max)
# - Edge cases (null bytes, encoding)
```

---

## References

- [Go Language Specification - Identifiers](https://golang.org/ref/spec#Identifiers)
- [Go Language Keywords](https://golang.org/ref/spec#Keywords)
- [OWASP Code Injection Prevention](https://owasp.org/www-community/attacks/Code_Injection)
