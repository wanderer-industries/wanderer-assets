# Contributing to wanderer-sde

Thank you for your interest in contributing to wanderer-sde! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.25.4 or later
- Git
- Make (optional but recommended)

### Setting Up the Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/wanderer-sde.git
   cd wanderer-sde
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/guarzo/wanderer-sde.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Build the project:
   ```bash
   make build
   ```

## Development Workflow

### Creating a Branch

Create a feature branch from `main`:

```bash
git checkout main
git pull upstream main
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-xyz` for new features
- `fix/issue-123` for bug fixes
- `docs/update-readme` for documentation changes

### Making Changes

1. Write your code following the project's coding standards
2. Add or update tests as needed
3. Ensure all tests pass:
   ```bash
   make test
   ```
4. Run the linter (if available):
   ```bash
   make lint
   ```

### Commit Messages

Follow conventional commit message format:

```
type(scope): short description

Longer description if needed.

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

Examples:
```
feat(parser): add support for new SDE field

fix(transformer): correct security calculation for edge cases

docs(readme): update installation instructions
```

### Submitting a Pull Request

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
2. Open a Pull Request against `main` on the upstream repository
3. Fill out the PR template with:
   - Description of changes
   - Related issues
   - Testing performed
4. Wait for review and address any feedback

## Code Standards

### Go Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format code
- Run `go vet` to catch common issues
- Keep functions focused and small
- Add comments for exported functions and types

### Project Structure

```
internal/       # Private application code
  config/       # Configuration handling
  downloader/   # SDE download logic
  parser/       # YAML parsing
  transformer/  # Data transformation
  writer/       # JSON output
  models/       # Data structures
pkg/            # Public libraries (if any)
cmd/            # Application entry points
```

### Testing

- Write unit tests for new functionality
- Place tests in `_test.go` files alongside the code
- Use table-driven tests where appropriate
- Aim for meaningful test coverage, not just high numbers

Example test structure:

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    validInput,
            expected: expectedOutput,
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    invalidInput,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if !tt.wantErr && result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Error Handling

- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Define custom error types in `internal/config/errors.go` for domain-specific errors

## Testing with Real Data

### Downloading the SDE

For integration testing, you can download the actual SDE:

```bash
./bin/sdeconvert --download --output ./test-output --verbose
```

### Sample Data

For unit tests, use minimal sample YAML data in test fixtures rather than the full SDE.

## Reporting Issues

### Bug Reports

Include:
- Go version (`go version`)
- Operating system and architecture
- Steps to reproduce
- Expected vs actual behavior
- Any error messages or logs

### Feature Requests

Include:
- Use case description
- Proposed solution (if any)
- Alternatives considered

## Questions?

If you have questions about contributing:

1. Check existing issues and discussions
2. Open a new issue with the `question` label

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
