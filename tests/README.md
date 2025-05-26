# Microservices Testing Framework

This directory contains testing infrastructure for the microservices-based project. The tests are organized independently from the main service code to promote isolation and reusability.

## Directory Structure

```
/tests
├── unit/               # Unit tests for individual components
│   ├── api-gateway/    # API Gateway unit tests
│   ├── inventory-service/
│   ├── order-service/
│   └── user-service/
├── integration/        # Integration tests between components
│   ├── api-gateway/
│   ├── inventory-service/
│   ├── order-service/
│   ├── user-service/
│   └── workflows/      # Cross-service workflow tests
├── e2e/                # End-to-end application tests
│   └── scenarios/      # Business-focused test scenarios
├── shared/             # Shared testing utilities
│   ├── mocks/          # Mock implementations
│   ├── fixtures/       # Test fixtures and data
│   └── utils/          # Common test utilities
├── testdata/           # Static test data files
└── scripts/            # Test runner scripts
    ├── setup.sh        # Environment setup script
    └── run-tests.sh    # Test execution script
```

## Testing Frameworks

This testing framework uses the following libraries and tools:

### Unit Tests
- Go's built-in `testing` package
- `github.com/stretchr/testify` for assertions and mocking
- `github.com/golang/mock/gomock` for interface mocking

### Integration Tests
- `github.com/testcontainers/testcontainers-go` for containerized dependencies
- `github.com/gavv/httpexpect/v2` for API testing

### E2E Tests
- `github.com/cucumber/godog` for BDD-style tests

## Configuration

For each service, create a `.env.test` file in the service root directory with test-specific configurations.

## Running Tests

Execute tests using the provided script:

```bash
# Run all tests
./scripts/run-tests.sh

# Run tests for a specific service
./scripts/run-tests.sh --service=inventory-service

# Run specific test types
./scripts/run-tests.sh --type=unit
./scripts/run-tests.sh --type=integration
./scripts/run-tests.sh --type=e2e

# Run tests with specific tags
./scripts/run-tests.sh --tags=redis,cache

# Run tests in parallel
./scripts/run-tests.sh --parallel

# Get verbose output
./scripts/run-tests.sh --verbose
```

## Adding New Tests

### Unit Tests

Unit tests should be placed in the corresponding service directory under `/tests/unit/`. Follow Go's standard test naming conventions:

```go
func TestMyFunction(t *testing.T) {
    // Test code here
}
```

### Integration Tests

Integration tests should be placed in the `/tests/integration/` directory. For service-specific integration tests, use the appropriate service subdirectory.

For cross-service workflows, place tests in the `/tests/integration/workflows/` directory.

### End-to-End Tests

End-to-end tests should be written in BDD style using GoDoc and placed in the `/tests/e2e/scenarios/` directory.

## Best Practices

1. Keep tests isolated from production code
2. Use mocks for external dependencies in unit tests
3. Use test containers for integration tests
4. Keep test data in the `testdata/` directory
5. Reuse test utilities from `shared/` directory
6. Tag tests appropriately for selective execution 