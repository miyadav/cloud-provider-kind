// pkg/testing/integration_test.go
package testing

import (
	"context"
	"testing"

	"sigs.k8s.io/cloud-provider-kind/pkg/provider"

	testingInterface "github.com/miyadav/cloud-provider-testing-interface"
)

// TestLoadBalancerIntegration tests load balancer integration
func TestLoadBalancerIntegration(t *testing.T) {
	// Check if load balancer creation is supported in this environment
	if !IsLoadBalancerSupported() {
		t.Skip("Load balancer creation not supported in this environment - skipping tests")
		return
	}

	// Create your cloud provider instance
	cloudProvider := provider.New("test-cluster", nil)

	// Create test implementation
	testImpl := NewCloudProviderTestImplementation(cloudProvider)

	// Create test runner
	runner := testingInterface.NewTestRunner(testImpl)

	// Add test suites
	runner.AddTestSuite(CreateLoadBalancerTestSuite())

	// Run tests
	ctx := context.Background()
	err := runner.RunTests(ctx)
	if err != nil {
		// Check if this is an environment-related error
		if IsEnvironmentError(err) {
			t.Skipf("Load balancer tests skipped due to environment constraints: %v", err)
			return
		}
		t.Fatalf("Test execution failed: %v", err)
	}

	// Verify results
	summary := runner.GetSummary()

	if summary.TotalTests == 0 {
		t.Skip("No tests were executed - likely due to environment constraints")
		return
	}

	if summary.FailedTests > 0 {
		t.Errorf("Some tests failed: %d failed out of %d total", summary.FailedTests, summary.TotalTests)
	}

	// Print results for debugging
	t.Logf("Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)
}
