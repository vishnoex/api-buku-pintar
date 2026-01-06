# GitHub Copilot Guidelines for Buku Pintar API

## Project Overview
This is a Go-based API for the Buku Pintar platform, featuring OAuth2 authentication, RBAC (Role-Based Access Control), content management, and payment processing.

## Code Style & Conventions

### General Go Practices
- Follow standard Go formatting (gofmt/goimports)
- Use meaningful variable and function names
- Keep functions small and focused (single responsibility)
- Add comments for exported functions and complex logic
- Error messages should be lowercase and not end with punctuation

### Project Structure
```
internal/
  ├── constant/     # Application constants
  ├── delivery/     # HTTP handlers/controllers
  ├── domain/       # Domain models and interfaces
  ├── helper/       # Utility functions
  ├── repository/   # Database layer
  ├── service/      # Business logic
  └── usecase/      # Use case implementations
```

### Architecture Patterns
- **Clean Architecture**: Follow the dependency rule (domain → usecase → delivery)
- **Repository Pattern**: All database operations go through repositories
- **Dependency Injection**: Use constructors to inject dependencies
- **Interface-based Design**: Define interfaces in domain layer

## Naming Conventions

### Files
- Use snake_case: `user_repository.go`, `token_service.go`
- Test files: `user_repository_test.go`
- Interface files in domain: `user.go` (contains UserRepository interface)

### Functions & Methods
- Use camelCase for private: `getUserByID()`
- Use PascalCase for exported: `GetUserByID()`
- Handler methods: `HandleCreateUser()`, `HandleGetUserByID()`
- Repository methods: `Create()`, `FindByID()`, `Update()`, `Delete()`
- Service methods: Describe business action: `AuthenticateUser()`, `ValidateToken()`

### Variables
- Use camelCase: `userId`, `accessToken`
- Constants: Use PascalCase or UPPER_SNAKE_CASE
- Context variables: `ctx context.Context` (always first parameter)

### Types
- Use PascalCase: `UserRequest`, `TokenResponse`
- Suffix with purpose: `*Handler`, `*Service`, `*Repository`, `*UseCase`

## Authentication & Authorization

### OAuth2 Implementation
- Store tokens securely in database with encryption
- Use refresh tokens for long-lived sessions
- Implement token rotation on refresh
- Always validate tokens before processing requests

### RBAC System
- Use permission-based authorization (not just roles)
- Cache permissions for performance
- Check permissions in middleware
- Document required permissions for each endpoint

### Security Best Practices
- Never log sensitive data (tokens, passwords, credentials)
- Use prepared statements for SQL queries (prevent injection)
- Validate and sanitize all user inputs
- Use HTTPS in production
- Implement rate limiting on authentication endpoints

## Database Operations

### Migrations
- Use sequential numbering: `000001_`, `000002_`
- Always create both up and down migrations
- Test migrations before committing
- Never modify existing migrations in production

### Repository Layer
```go
// Example repository interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}
```

### Transaction Handling
- Use transactions for multi-step operations
- Always defer rollback with error check
- Commit only after all operations succeed

## Error Handling

### Error Patterns
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}

// Custom domain errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidToken = errors.New("invalid token")
)

// HTTP error responses
return c.JSON(http.StatusBadRequest, map[string]string{
    "error": "invalid request parameters",
})
```

### Error Response Format
```json
{
    "error": "descriptive error message",
    "code": "ERROR_CODE_IF_NEEDED"
}
```

## API Design

### RESTful Conventions
- Use plural nouns: `/users`, `/articles`, `/ebooks`
- Use HTTP methods appropriately: GET, POST, PUT, PATCH, DELETE
- Use path parameters for IDs: `/users/:id`
- Use query parameters for filtering: `/articles?status=published&category=tech`

### Request/Response Structure
```go
// Request DTOs
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}

// Response DTOs
type UserResponse struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Validation
- Use struct tags for validation
- Validate in handler layer before passing to service
- Return meaningful validation error messages

## Testing

### Test Structure
- Unit tests for services and use cases
- Integration tests for repositories
- Handler tests with mocked dependencies
- Use table-driven tests when appropriate

### Test Naming
```go
func TestUserService_AuthenticateUser(t *testing.T) { }
func TestUserService_AuthenticateUser_InvalidCredentials(t *testing.T) { }
```

### Mock Interfaces
- Create mocks for external dependencies
- Use interfaces to make code testable
- Keep mocks in `*_mock.go` files

## Configuration

### Environment Variables
- Use `config.json` for local development
- Use environment variables in production
- Never commit sensitive configuration
- Provide `example.config.json` as template

### Configuration Loading
- Load config at startup
- Validate required fields
- Use default values where appropriate

## Documentation

### Code Comments
```go
// GetUserByID retrieves a user by their unique identifier.
// Returns ErrUserNotFound if the user does not exist.
func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
    // implementation
}
```

### API Documentation
- Document new endpoints in relevant docs/ files
- Include request/response examples
- Document required permissions
- Update README.md when adding major features

### Changelog Files
- Create summary docs for major implementations
- Include quick reference guides for complex features
- Document architectural decisions

## Performance Considerations

### Database Optimization
- Use indexes appropriately
- Avoid N+1 queries
- Use pagination for large datasets
- Cache frequently accessed data

### Caching Strategy
- Cache permissions (with TTL)
- Cache static content metadata
- Implement cache invalidation properly

### Context Management
- Always use context.Context for cancellation
- Set appropriate timeouts
- Pass context through all layers

## Common Patterns

### Handler Pattern
```go
func (h *UserHandler) HandleGetUser(c echo.Context) error {
    id := c.Param("id")
    
    user, err := h.userService.GetUserByID(c.Request().Context(), id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            return c.JSON(http.StatusNotFound, map[string]string{
                "error": "user not found",
            })
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "internal server error",
        })
    }
    
    return c.JSON(http.StatusOK, user)
}
```

### Service Pattern
```go
func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Validate business rules
    existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
    if existing != nil {
        return nil, ErrUserAlreadyExists
    }
    
    // Hash password
    hashedPassword, err := s.cryptoService.HashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }
    
    // Create user
    user := &User{
        ID:       uuid.New().String(),
        Email:    req.Email,
        Name:     req.Name,
        Password: hashedPassword,
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}
```

## Git Commit Messages
- Use conventional commits format: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`
- Be descriptive: `feat: add user profile update endpoint`
- Reference issues when applicable: `fix: resolve token refresh issue #123`

## When Suggesting Code

### Always Consider
1. Is this following clean architecture principles?
2. Are errors properly handled and wrapped?
3. Is the code testable?
4. Are there security implications?
5. Is the code performant?
6. Does it follow project conventions?

### Prefer
- Explicit over implicit
- Interfaces over concrete types (in domain layer)
- Composition over inheritance
- Pure functions when possible
- Clear error messages

### Avoid
- Global variables (except constants)
- Magic numbers (use named constants)
- Commented-out code
- God objects or functions
- Tight coupling between layers

## Resources
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

**Remember**: Write code that is easy to read, maintain, and test. When in doubt, prioritize clarity over cleverness.
