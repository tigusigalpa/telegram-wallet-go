# Setup Guide for telegram-wallet-go

## Initial Setup

### 1. Initialize Go Module

```bash
cd public_html/packages/telegram-wallet-go
go mod tidy
```

This will download all dependencies including `testify` for tests.

### 2. Run Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 3. Verify Build

```bash
# Build examples
go build ./examples/create_order
go build ./examples/webhook_server

# Run linter (install golangci-lint first)
golangci-lint run
```

### 4. Publish to GitHub

```bash
git init
git add .
git commit -m "Initial commit: Telegram Wallet Pay Go SDK v1.0.0"
git branch -M main
git remote add origin https://github.com/tigusigalpa/telegram-wallet-go.git
git push -u origin main
```

### 5. Create GitHub Release and Tag

```bash
# Create and push tag
git tag v1.0.0
git push origin v1.0.0
```

Then create a release on GitHub:
1. Go to https://github.com/tigusigalpa/telegram-wallet-go/releases
2. Click "Create a new release"
3. Select tag: `v1.0.0`
4. Release title: `v1.0.0 - Initial Release`
5. Description: Copy from README features section
6. Publish release

### 6. Verify on pkg.go.dev

After pushing the tag, the package will automatically appear on pkg.go.dev within a few minutes:
- https://pkg.go.dev/github.com/tigusigalpa/telegram-wallet-go

You can manually request indexing:
- https://pkg.go.dev/github.com/tigusigalpa/telegram-wallet-go@v1.0.0

## Development Workflow

### Running Tests

```bash
# All tests
go test ./...

# Verbose output
go test -v ./...

# Specific package
go test -v ./middleware

# With race detector
go test -race ./...

# Benchmark tests (if you add any)
go test -bench=. ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (install golangci-lint)
golangci-lint run

# Check for security issues (install gosec)
gosec ./...
```

### Building Examples

```bash
# Create order example
go build -o create_order ./examples/create_order
./create_order

# Webhook server example
go build -o webhook_server ./examples/webhook_server
./webhook_server
```

### Testing with Different Frameworks

#### Gin Framework

```bash
# Install Gin
go get -u github.com/gin-gonic/gin

# Build with gin tag
go build -tags gin ./...

# Test
go test -tags gin ./middleware
```

#### Echo Framework

```bash
# Install Echo
go get -u github.com/labstack/echo/v4

# Build with echo tag
go build -tags echo ./...

# Test
go test -tags echo ./middleware
```

## Integration Testing

### Create a Test Application

```bash
mkdir test-app
cd test-app
go mod init example.com/test-app

# Use local package for testing
go mod edit -replace github.com/tigusigalpa/telegram-wallet-go=../telegram-wallet-go
go get github.com/tigusigalpa/telegram-wallet-go
```

Create `main.go`:
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/tigusigalpa/telegram-wallet-go"
)

func main() {
    apiKey := os.Getenv("WALLETPAY_API_KEY")
    client := walletpay.NewClient(apiKey)
    
    count, err := client.GetOrderAmount(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Total orders: %d\n", count)
}
```

Run:
```bash
export WALLETPAY_API_KEY=your_test_api_key
go run main.go
```

## Versioning

### Creating New Releases

```bash
# Update version in code/docs if needed
# Commit changes
git add .
git commit -m "Release v1.1.0: Add new features"

# Create and push tag
git tag v1.1.0
git push origin v1.1.0

# Create GitHub release
# pkg.go.dev will automatically index the new version
```

### Semantic Versioning

- `v1.0.0` - Initial release
- `v1.0.1` - Patch (bug fixes)
- `v1.1.0` - Minor (new features, backward compatible)
- `v2.0.0` - Major (breaking changes)

## Troubleshooting

### Module Not Found

If users report "module not found":
```bash
# Ensure tag is pushed
git push origin v1.0.0

# Verify on GitHub
# Check https://github.com/tigusigalpa/telegram-wallet-go/tags

# Request indexing on pkg.go.dev
# Visit https://pkg.go.dev/github.com/tigusigalpa/telegram-wallet-go@v1.0.0
```

### Test Dependencies Not Found

```bash
# Download dependencies
go mod download

# Or tidy
go mod tidy
```

### Import Cycle Errors

Ensure middleware packages don't import the main package circularly. Current structure is correct.

### Build Tags Not Working

Ensure you're using the correct build command:
```bash
go build -tags gin ./...
go build -tags echo ./...
```

## Documentation

### Update pkg.go.dev Documentation

Documentation is automatically generated from:
- Package comments in `.go` files
- README.md
- Go doc comments

To improve documentation:
1. Add comprehensive package-level comments
2. Document all exported types and functions
3. Include examples in doc comments

Example:
```go
// Package walletpay provides a client for the Telegram Wallet Pay API.
//
// Example usage:
//
//	client := walletpay.NewClient("YOUR_API_KEY")
//	order, err := client.CreateOrder(ctx, walletpay.CreateOrderRequest{
//	    Amount: walletpay.MoneyAmount{CurrencyCode: "USD", Amount: "9.99"},
//	    Description: "Premium subscription",
//	    ExternalID: "ORDER-123",
//	    TimeoutSeconds: 3600,
//	    CustomerTelegramUserID: 123456789,
//	})
package walletpay
```

## CI/CD Setup (Optional)

### GitHub Actions

Create `.github/workflows/test.yml`:
```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
    
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Test
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
```

## Support

For issues or questions:
- GitHub Issues: https://github.com/tigusigalpa/telegram-wallet-go/issues
- Email: sovletig@gmail.com
