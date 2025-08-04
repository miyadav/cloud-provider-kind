// cmd/test-runner/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"sigs.k8s.io/cloud-provider-kind/pkg/provider"
	"sigs.k8s.io/cloud-provider-kind/pkg/testing"

	testingInterface "github.com/miyadav/cloud-provider-testing-interface"
)

func main() {
	var (
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		testType    = flag.String("test-type", "all", "Type of tests to run: all, loadbalancer, clusters, instances, provider")
		clusterName = flag.String("cluster", "test-cluster", "Name of the test cluster")
		showEnv     = flag.Bool("show-env", false, "Show environment information and exit")
	)
	flag.Parse()

	// Show environment information if requested
	if *showEnv {
		showEnvironmentInfo()
		return
	}

	// Create cloud provider instance
	cloudProvider := provider.New(*clusterName, nil)

	// Create test implementation
	testImpl := testing.NewCloudProviderTestImplementation(cloudProvider)

	// Create test runner
	runner := testingInterface.NewTestRunner(testImpl)

	// Add test suites based on test type
	switch strings.ToLower(*testType) {
	case "all":
		addAllTestSuites(runner)
	case "loadbalancer":
		runner.AddTestSuite(testing.CreateLoadBalancerTestSuite())
	case "clusters":
		runner.AddTestSuite(testing.CreateClustersTestSuite())
	case "instances":
		runner.AddTestSuite(testing.CreateInstancesTestSuite())
	case "provider":
		runner.AddTestSuite(testing.CreateProviderTestSuite())
	default:
		log.Fatalf("Unknown test type: %s. Valid options: all, loadbalancer, clusters, instances, provider", *testType)
	}

	// Run tests
	ctx := context.Background()
	if *verbose {
		fmt.Println("Starting test execution...")
		fmt.Printf("Test type: %s\n", *testType)
		fmt.Printf("Cluster name: %s\n", *clusterName)
	}

	err := runner.RunTests(ctx)
	if err != nil {
		if testing.IsEnvironmentError(err) {
			fmt.Printf("Tests skipped due to environment constraints: %v\n", err)
			os.Exit(0)
		}
		log.Fatalf("Test execution failed: %v", err)
	}

	// Get and display results
	summary := runner.GetSummary()
	displayResults(&summary, *verbose)

	// Exit with appropriate code
	if summary.FailedTests > 0 {
		os.Exit(1)
	}
	os.Exit(0)
}

// addAllTestSuites adds all available test suites to the runner
func addAllTestSuites(runner *testingInterface.TestRunner) {
	// Add provider tests first (they don't require external dependencies)
	runner.AddTestSuite(testing.CreateProviderTestSuite())

	// Add clusters tests
	if testing.IsClustersSupported() {
		runner.AddTestSuite(testing.CreateClustersTestSuite())
	} else {
		fmt.Println("Skipping clusters tests - not supported in this environment")
	}

	// Add instances tests
	if testing.IsInstancesSupported() {
		runner.AddTestSuite(testing.CreateInstancesTestSuite())
	} else {
		fmt.Println("Skipping instances tests - not supported in this environment")
	}

	// Add load balancer tests (if supported)
	if testing.IsLoadBalancerSupported() {
		runner.AddTestSuite(testing.CreateLoadBalancerTestSuite())
	} else {
		fmt.Println("Skipping load balancer tests - not supported in this environment")
	}
}

// showEnvironmentInfo displays information about the test environment
func showEnvironmentInfo() {
	envInfo := testing.GetTestEnvironmentInfo()

	fmt.Println("=== Test Environment Information ===")
	fmt.Printf("Load Balancer Supported: %v\n", envInfo["loadbalancer_supported"])
	fmt.Printf("Clusters Supported: %v\n", envInfo["clusters_supported"])
	fmt.Printf("Instances Supported: %v\n", envInfo["instances_supported"])
	fmt.Printf("Provider Supported: %v\n", envInfo["provider_supported"])
	fmt.Printf("Docker Available: %v\n", envInfo["docker_available"])
	fmt.Printf("Podman Available: %v\n", envInfo["podman_available"])
	fmt.Printf("KIND Available: %v\n", envInfo["kind_available"])
	fmt.Printf("Container Runtime Available: %v\n", envInfo["container_runtime"])
	fmt.Println("=====================================")
}

// displayResults displays test results in a formatted way
func displayResults(summary *testingInterface.TestSummary, verbose bool) {
	fmt.Println("\n=== Test Results ===")
	fmt.Printf("Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("Passed: %d\n", summary.PassedTests)
	fmt.Printf("Failed: %d\n", summary.FailedTests)
	fmt.Printf("Skipped: %d\n", summary.SkippedTests)

	if summary.TotalTests == 0 {
		fmt.Println("No tests were executed - likely due to environment constraints")
		return
	}

	successRate := float64(summary.PassedTests) / float64(summary.TotalTests) * 100
	fmt.Printf("Success Rate: %.1f%%\n", successRate)

	if verbose {
		fmt.Println("\n=== Test Logs ===")
		fmt.Println("Detailed logs would be available in verbose mode")
	}

	if summary.FailedTests > 0 {
		fmt.Println("\n❌ Some tests failed!")
	} else {
		fmt.Println("\n✅ All tests passed!")
	}
	fmt.Println("===================")
}
