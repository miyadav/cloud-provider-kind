# Cloud Provider KIND Testing Package

This package provides comprehensive integration testing for the KIND cloud provider, covering all major cloud provider interfaces and functionality.

## Overview

The testing package implements a robust testing framework that:

- **Environment-Aware**: Automatically detects and adapts to the testing environment
- **Comprehensive Coverage**: Tests all major cloud provider interfaces
- **Graceful Degradation**: Skips tests that cannot run in the current environment
- **Detailed Reporting**: Provides clear test results and environment information

## Test Suites

### 1. Load Balancer Tests (`CreateLoadBalancerTestSuite`)

Tests load balancer creation, update, and deletion functionality.

**Requirements:**
- Container runtime (Docker, Podman, nerdctl, or Finch)
- Privileged container support
- Sysctl parameter setting capability

**Tests:**
- Load Balancer Creation
- Load Balancer Update
- Load Balancer Deletion

### 2. Clusters Tests (`CreateClustersTestSuite`)

Tests cluster management functionality including listing clusters and retrieving master node information.

**Requirements:**
- KIND CLI available
- Ability to list clusters

**Tests:**
- List Clusters
- Get Master Node
- Cluster Operations

### 3. Instances Tests (`CreateInstancesTestSuite`)

Tests instance management functionality including existence checks, shutdown status, and metadata retrieval.

**Requirements:**
- KIND CLI available
- Ability to access cluster nodes

**Tests:**
- Instance Exists
- Instance Shutdown Status
- Instance Metadata

### 4. Provider Tests (`CreateProviderTestSuite`)

Tests basic provider functionality including provider name, initialization, and interface support reporting.

**Requirements:**
- No external dependencies (always supported)

**Tests:**
- Provider Name
- Provider Initialization
- Provider Interface Support

## Usage

### Running Tests via Go Test

```bash
# Run all tests
go test ./pkg/testing -v

# Run specific test suites
go test ./pkg/testing -run TestLoadBalancerIntegration -v
go test ./pkg/testing -run TestClustersIntegration -v
go test ./pkg/testing -run TestInstancesIntegration -v
go test ./pkg/testing -run TestProviderIntegration -v
go test ./pkg/testing -run TestAllIntegration -v
```

### Running Tests via Test Runner

```bash
# Build the test runner
go build -o bin/test-runner cmd/test-runner/main.go

# Show environment information
./bin/test-runner -show-env

# Run all tests
./bin/test-runner

# Run specific test types
./bin/test-runner -test-type loadbalancer
./bin/test-runner -test-type clusters
./bin/test-runner -test-type instances
./bin/test-runner -test-type provider

# Run with verbose output
./bin/test-runner -verbose

# Specify cluster name
./bin/test-runner -cluster my-test-cluster
```

## Environment Detection

The testing package automatically detects the testing environment capabilities:

### Load Balancer Support
- Checks for container runtime availability
- Tests privileged container creation
- Verifies sysctl parameter setting

### Clusters Support
- Checks for KIND CLI availability
- Tests cluster listing capability

### Instances Support
- Checks for KIND CLI availability
- Tests cluster node access

### Provider Support
- Always available (no external dependencies)

## Test Configuration

Tests are configured via the `TestConfig` structure:

```go
config := &testing.TestConfig{
    ProviderName:         "kind",
    ClusterName:          "test-cluster",
    Region:               "us-west-1",
    Zone:                 "us-west-1a",
    TestTimeout:          5 * time.Minute,
    CleanupResources:     true,
    MockExternalServices: false,
}
```

## Error Handling

The testing framework provides intelligent error handling:

### Environment Errors
- Automatically detected and categorized
- Tests are skipped gracefully when environment constraints prevent execution
- Clear error messages indicate why tests were skipped

### Test Failures
- Detailed error reporting with context
- Distinguishes between environment issues and actual test failures
- Provides actionable feedback for debugging

## Integration with External Testing Interface

The testing package integrates with the `cloud-provider-testing-interface` to provide:

- Standardized test execution
- Consistent test reporting
- Reusable test infrastructure
- Cross-provider compatibility

## File Structure

```
pkg/testing/
├── integration_test.go      # Main integration test functions
├── test_suites.go          # Test suite definitions and test implementations
├── test_utils.go           # Environment detection and utility functions
├── cloud_provider_test_impl.go  # Cloud provider test implementation
├── config/
│   └── test-config.yaml    # Default test configuration
└── README.md               # This documentation
```

## Best Practices

### Writing New Tests

1. **Environment Detection**: Always check if the required functionality is supported
2. **Graceful Degradation**: Handle cases where external dependencies are unavailable
3. **Resource Cleanup**: Ensure proper cleanup of test resources
4. **Error Context**: Provide meaningful error messages for debugging

### Test Implementation

```go
func testExample(ti testing.TestInterface) error {
    ctx := context.Background()
    cloud := ti.GetCloudProvider()
    
    // Test implementation
    // ...
    
    ti.GetTestResults().AddLog("Test completed successfully")
    return nil
}
```

### Environment Checks

```go
// Check if functionality is supported
if !IsLoadBalancerSupported() {
    return fmt.Errorf("load balancer not supported in this environment")
}
```

## Troubleshooting

### Common Issues

1. **Tests Skipped**: Check environment requirements using `-show-env`
2. **Container Runtime Issues**: Verify Docker/Podman installation and permissions
3. **KIND Issues**: Ensure KIND CLI is installed and accessible
4. **Permission Errors**: Run with appropriate privileges for container operations

### Debug Mode

Use verbose mode for detailed output:

```bash
./bin/test-runner -verbose -test-type all
```

### Environment Information

Check environment capabilities:

```bash
./bin/test-runner -show-env
```

## Contributing

When adding new tests:

1. Follow the existing test structure and patterns
2. Add appropriate environment detection
3. Include comprehensive error handling
4. Update this documentation
5. Add test cases for both success and failure scenarios

## Dependencies

- `github.com/miyadav/cloud-provider-testing-interface`: External testing interface
- `k8s.io/cloud-provider`: Kubernetes cloud provider interfaces
- `k8s.io/api/core/v1`: Kubernetes core API types
- `sigs.k8s.io/kind/pkg/cluster`: KIND cluster management 