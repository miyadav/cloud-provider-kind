# Cloud Provider Testing Interface

This package provides a comprehensive testing interface for Kubernetes cloud providers, demonstrating how to implement and run tests using the standardized external testing interface from [github.com/miyadav/cloud-provider-testing-interface](https://github.com/miyadav/cloud-provider-testing-interface).

## Overview

The implementation demonstrates how cloud providers can use the external testing interface to achieve standardization across all cloud provider implementations. It provides:

- **Standardized Testing Interface**: Uses the external `TestInterface` for consistent testing across providers
- **Load Balancer Testing**: Tests for basic load balancer operations, updates, and deletions
- **Instance Testing**: Tests for instance existence, metadata, and shutdown status
- **Cluster Testing**: Tests for cluster listing and master endpoint retrieval
- **Provider Name Testing**: Basic provider identification testing

## Architecture

### Core Components

1. **`pkg/testing/kind_test_interface.go`** - KIND implementation of the external `TestInterface`
2. **`pkg/testing/kind_test_suite.go`** - Test suites using the external interface
3. **`cmd/test/main.go`** - Command-line tool using the external interface
4. **`examples/external-interface-example.go`** - Examples for other cloud providers

### Key Interfaces

#### External TestInterface
```go
// From github.com/miyadav/cloud-provider-testing-interface
type TestInterface interface {
    SetupTestEnvironment(config *TestConfig) error
    TeardownTestEnvironment() error
    GetCloudProvider() cloudprovider.Interface
    CreateTestNode(ctx context.Context, nodeConfig *TestNodeConfig) (*v1.Node, error)
    DeleteTestNode(ctx context.Context, nodeName string) error
    CreateTestService(ctx context.Context, serviceConfig *TestServiceConfig) (*v1.Service, error)
    DeleteTestService(ctx context.Context, serviceName string) error
    CreateTestRoute(ctx context.Context, routeConfig *TestRouteConfig) (*cloudprovider.Route, error)
    DeleteTestRoute(ctx context.Context, routeName string) error
    WaitForCondition(ctx context.Context, condition TestCondition) error
    GetTestResults() *TestResults
    ResetTestState() error
}
```

#### KIND Implementation
The `KindTestInterface` provides a concrete implementation of the external `TestInterface` for KIND cloud provider, demonstrating how to:
- Create and manage Kubernetes resources for testing
- Wait for cloud provider operations to complete
- Verify expected states and behaviors
- Clean up test resources

#### Test Suites
The `KindTestSuite` provides test suites that use the external interface, demonstrating how to:
- Create standardized test suites that work with any cloud provider
- Implement provider-specific test logic
- Use the external test runner for execution

## Usage

### Prerequisites

Before running the tests, ensure you have:
1. A Kubernetes cluster running (KIND cluster recommended)
2. A valid kubeconfig file (defaults to `~/.kube/config`)
3. The cloud-provider-kind running in your cluster

### Basic Usage with External Interface

```go
package main

import (
    "context"
    "sigs.k8s.io/cloud-provider-kind/pkg/testing"
    "sigs.k8s.io/kind/pkg/cluster"
    
    // Import the external testing interface
    exttesting "github.com/miyadav/cloud-provider-testing-interface"
)

func main() {
    // Create KIND provider
    kindProvider := cluster.NewProvider()
    kubeClient := createKubeClient()
    
    // Create test interface using external interface
    testInterface := testing.NewKindTestInterface("my-cluster", kindProvider, kubeClient)
    
    // Create test runner using external interface
    runner := exttesting.NewTestRunner(testInterface)
    
    // Add test suites
    testSuite := testing.NewKindTestSuite("my-cluster", kindProvider, kubeClient)
    runner.AddTestSuite(testSuite.GetLoadBalancerTestSuite())
    runner.AddTestSuite(testSuite.GetInstancesTestSuite())
    
    // Setup test environment
    config := &exttesting.TestConfig{
        ProviderName:        "kind",
        ClusterName:         "my-cluster",
        TestTimeout:         10 * time.Minute,
        CleanupResources:    true,
        MockExternalServices: false,
    }
    
    if err := testInterface.SetupTestEnvironment(config); err != nil {
        log.Fatalf("Failed to setup test environment: %v", err)
    }
    
    // Run tests
    if err := runner.RunTests(context.Background()); err != nil {
        log.Fatalf("Tests failed: %v", err)
    }
    
    fmt.Println("All tests passed!")
}
```

### Running Specific Test Suites

```bash
# Build the test tool
make test-cloud-provider

# Show help
./bin/cloud-provider-test -help

# Run all test suites (uses default kubeconfig ~/.kube/config)
./bin/cloud-provider-test

# Run specific test suite
./bin/cloud-provider-test -suite loadbalancer
./bin/cloud-provider-test -suite instances
./bin/cloud-provider-test -suite clusters
./bin/cloud-provider-test -suite provider-name

# Use specific configuration
./bin/cloud-provider-test -kubeconfig /path/to/kubeconfig -cluster my-cluster

# Examples
./bin/cloud-provider-test                                    # Run all tests using default kubeconfig
./bin/cloud-provider-test -suite loadbalancer               # Run only load balancer tests
./bin/cloud-provider-test -kubeconfig /path/to/config       # Use specific kubeconfig
./bin/cloud-provider-test -cluster my-cluster               # Test specific KIND cluster
```

### Available Test Suites

| Suite Name | Description |
|------------|-------------|
| `loadbalancer` | Tests load balancer creation, updates, and deletion |
| `instances` | Tests instance existence, metadata, and shutdown status |
| `clusters` | Tests cluster listing and master endpoint retrieval |
| `provider-name` | Tests provider name functionality |

## Standardization Benefits

### 1. **Consistent Interface**
All cloud providers implement the same `TestInterface`, ensuring:
- Consistent method signatures
- Standardized test configuration
- Uniform test results format

### 2. **Reusable Test Suites**
Test suites can be shared across providers:
```go
// This test suite works with any cloud provider
func createStandardLoadBalancerTestSuite() exttesting.TestSuite {
    return exttesting.TestSuite{
        Name: "Standard Load Balancer Tests",
        Tests: []exttesting.Test{
            {
                Name: "Load Balancer Creation",
                Run: func(testInterface exttesting.TestInterface) error {
                    // This works with AWS, GCP, Azure, KIND, etc.
                    cloud := testInterface.GetCloudProvider()
                    // ... test logic
                    return nil
                },
            },
        },
    }
}
```

### 3. **Easy Provider Comparison**
Run the same tests across different providers:
```go
// Test AWS
awsTestInterface := NewAWSTestInterface(awsCloud, awsKubeClient)
awsRunner := exttesting.NewTestRunner(awsTestInterface)
awsRunner.AddTestSuite(standardSuite)

// Test GCP
gcpTestInterface := NewGCPTestInterface(gcpCloud, gcpKubeClient)
gcpRunner := exttesting.NewTestRunner(gcpTestInterface)
gcpRunner.AddTestSuite(standardSuite)

// Same tests, different providers
```

## Implementing for Other Cloud Providers

### AWS Cloud Provider Example

```go
type AWSTestInterface struct {
    cloudProvider cloudprovider.Interface
    kubeClient    kubernetes.Interface
    config        *exttesting.TestConfig
    results       *exttesting.TestResults
}

func (aws *AWSTestInterface) SetupTestEnvironment(config *exttesting.TestConfig) error {
    aws.config = config
    aws.results.AddLog("AWS test environment setup completed")
    return nil
}

func (aws *AWSTestInterface) GetCloudProvider() cloudprovider.Interface {
    return aws.cloudProvider
}

func (aws *AWSTestInterface) CreateTestService(ctx context.Context, serviceConfig *exttesting.TestServiceConfig) (*v1.Service, error) {
    // AWS-specific service creation logic
    aws.results.IncrementResourceCount("aws-services")
    aws.results.AddLog(fmt.Sprintf("Created AWS test service: %s/%s", serviceConfig.Namespace, serviceConfig.Name))
    return &v1.Service{}, nil
}

// ... implement other methods
```

### GCP Cloud Provider Example

```go
type GCPTestInterface struct {
    cloudProvider cloudprovider.Interface
    kubeClient    kubernetes.Interface
    config        *exttesting.TestConfig
    results       *exttesting.TestResults
}

func (gcp *GCPTestInterface) SetupTestEnvironment(config *exttesting.TestConfig) error {
    gcp.config = config
    gcp.results.AddLog("GCP test environment setup completed")
    return nil
}

func (gcp *GCPTestInterface) GetCloudProvider() cloudprovider.Interface {
    return gcp.cloudProvider
}

func (gcp *GCPTestInterface) CreateTestService(ctx context.Context, serviceConfig *exttesting.TestServiceConfig) (*v1.Service, error) {
    // GCP-specific service creation logic
    gcp.results.IncrementResourceCount("gcp-services")
    gcp.results.AddLog(fmt.Sprintf("Created GCP test service: %s/%s", serviceConfig.Namespace, serviceConfig.Name))
    return &v1.Service{}, nil
}

// ... implement other methods
```

## Integration with CI/CD

The external testing interface can be easily integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
name: Cloud Provider Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.19
    
    - name: Run standardized cloud provider tests
      run: |
        go run cmd/test/main.go -suite loadbalancer
        go run cmd/test/main.go -suite instances
```

## Best Practices

1. **Implement External Interface**: Always implement the external `TestInterface` for consistency
2. **Provider-Specific Logic**: Keep provider-specific logic in the test interface implementation
3. **Standardized Test Suites**: Create test suites that work across all providers
4. **Resource Cleanup**: Always clean up test resources in teardown
5. **Error Handling**: Use consistent error handling patterns

## Extending the Testing Interface

To add new test cases using the external interface:

1. **Add to External Interface** (if needed):
```go
// In github.com/miyadav/cloud-provider-testing-interface
type TestInterface interface {
    // ... existing methods
    CreateTestVolume(ctx context.Context, volumeConfig *TestVolumeConfig) (*v1.PersistentVolume, error)
}
```

2. **Implement in Provider**:
```go
func (kti *KindTestInterface) CreateTestVolume(ctx context.Context, volumeConfig *exttesting.TestVolumeConfig) (*v1.PersistentVolume, error) {
    // Provider-specific implementation
    return &v1.PersistentVolume{}, nil
}
```

3. **Add to Test Suite**:
```go
func (kts *KindTestSuite) GetVolumeTestSuite() exttesting.TestSuite {
    return exttesting.TestSuite{
        Name: "Volume Tests",
        Tests: []exttesting.Test{
            {
                Name: "Volume Creation",
                Run: func(testInterface exttesting.TestInterface) error {
                    // Test logic
                    return nil
                },
            },
        },
    }
}
```

## Troubleshooting

### Common Issues

1. **External Interface Not Found**: Ensure `github.com/miyadav/cloud-provider-testing-interface` is in go.mod
2. **Test Interface Not Implemented**: Verify all methods of `TestInterface` are implemented
3. **Provider-Specific Errors**: Check provider-specific implementation details

### Debug Mode

Enable verbose logging for debugging:

```go
import "k8s.io/klog/v2"

// Set log level
klog.InitFlags(nil)
flag.Set("v", "4") // Verbose logging
flag.Parse()
```

## Contributing

When contributing to the testing interface:

1. Follow the external interface standards
2. Implement provider-specific logic in test interfaces
3. Create reusable test suites
4. Update documentation for new features
5. Ensure backward compatibility

## Related Links

- [External Testing Interface](https://github.com/miyadav/cloud-provider-testing-interface)
- [Kubernetes Cloud Provider Interface](https://kubernetes.io/docs/concepts/architecture/cloud-controller/)
- [KIND Cloud Provider](https://github.com/kubernetes-sigs/cloud-provider-kind) 