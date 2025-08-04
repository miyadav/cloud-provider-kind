// pkg/testing/example_usage.go
package testing

import (
	"context"
	"fmt"
	"log"

	testing "github.com/miyadav/cloud-provider-testing-interface"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cloud-provider-kind/pkg/provider"
)

// ExampleUsage demonstrates how to use the complete testing interface implementation
// This shows how to handle cases where the cloud provider doesn't support certain functionality
func ExampleUsage() {
	// 1. Create a cloud provider instance
	cloudProvider := provider.New("test-cluster", nil)

	// 2. Create the test implementation
	testImpl := NewCloudProviderTestImplementation(cloudProvider)

	// 3. Example: Check what features are supported
	fmt.Println("=== Feature Support Check ===")
	features := []string{"loadbalancer", "clusters", "instances", "provider", "unsupported_feature"}
	for _, feature := range features {
		supported := testImpl.IsFeatureSupported(feature)
		fmt.Printf("Feature '%s': %t\n", feature, supported)
	}

	// 4. Example: Get environment information
	fmt.Println("\n=== Environment Information ===")
	envInfo := testImpl.GetTestEnvironmentInfo()
	for key, value := range envInfo {
		fmt.Printf("%s: %v\n", key, value)
	}

	// 5. Example: Test cloud provider interface support
	fmt.Println("\n=== Cloud Provider Interface Support ===")
	cloud := testImpl.GetCloudProvider()

	// Test LoadBalancer support - returns (interface, bool)
	_, lbSupported := cloud.LoadBalancer()
	if lbSupported {
		fmt.Println("✓ LoadBalancer interface is supported")
		// You can use lb interface here
	} else {
		fmt.Println("✗ LoadBalancer interface is NOT supported")
		// lb will be nil, but that's okay - the interface contract allows this
	}

	// Test Clusters support
	_, clustersSupported := cloud.Clusters()
	if clustersSupported {
		fmt.Println("✓ Clusters interface is supported")
		// You can use clusters interface here
	} else {
		fmt.Println("✗ Clusters interface is NOT supported")
		// clusters will be nil, but that's okay
	}

	// Test InstancesV2 support
	_, instancesSupported := cloud.InstancesV2()
	if instancesSupported {
		fmt.Println("✓ InstancesV2 interface is supported")
		// You can use instances interface here
	} else {
		fmt.Println("✗ InstancesV2 interface is NOT supported")
		// instances will be nil, but that's okay
	}

	// Test Instances (v1) support - typically false for KIND
	_, instancesV1Supported := cloud.Instances()
	if instancesV1Supported {
		fmt.Println("✓ Instances (v1) interface is supported")
	} else {
		fmt.Println("✗ Instances (v1) interface is NOT supported")
		// instancesV1 will be nil, but that's okay
	}

	// Test Zones support - typically false for KIND
	_, zonesSupported := cloud.Zones()
	if zonesSupported {
		fmt.Println("✓ Zones interface is supported")
	} else {
		fmt.Println("✗ Zones interface is NOT supported")
		// zones will be nil, but that's okay
	}

	// Test Routes support - typically false for KIND
	_, routesSupported := cloud.Routes()
	if routesSupported {
		fmt.Println("✓ Routes interface is supported")
	} else {
		fmt.Println("✗ Routes interface is NOT supported")
		// routes will be nil, but that's okay
	}

	// Test HasClusterID
	hasClusterID := cloud.HasClusterID()
	fmt.Printf("HasClusterID: %t\n", hasClusterID)

	// 6. Example: Safe interface usage pattern
	fmt.Println("\n=== Safe Interface Usage Pattern ===")
	ExampleSafeInterfaceUsage(testImpl)
}

// ExampleSafeInterfaceUsage demonstrates the safe pattern for using cloud provider interfaces
func ExampleSafeInterfaceUsage(testImpl *CloudProviderTestImplementation) {
	cloud := testImpl.GetCloudProvider()
	ctx := context.Background()

	// Pattern 1: Check support before using interface
	if lb, supported := cloud.LoadBalancer(); supported {
		fmt.Println("LoadBalancer is supported, can use it safely")
		// Use lb interface here
		_ = lb // Placeholder for actual usage
	} else {
		fmt.Println("LoadBalancer is NOT supported, skipping load balancer operations")
		// lb is nil, but that's expected and safe
	}

	// Pattern 2: Handle unsupported functionality gracefully
	if clusters, supported := cloud.Clusters(); supported {
		fmt.Println("Clusters is supported, attempting to list clusters")
		clusterList, err := clusters.ListClusters(ctx)
		if err != nil {
			fmt.Printf("Error listing clusters: %v\n", err)
		} else {
			fmt.Printf("Successfully listed %d clusters\n", len(clusterList))
		}
	} else {
		fmt.Println("Clusters is NOT supported, skipping cluster operations")
		// This is perfectly fine - the cloud provider doesn't support clusters
	}

	// Pattern 3: Test node creation (always supported for testing)
	fmt.Println("Creating test node...")
	nodeConfig := &testing.TestNodeConfig{
		Name:         "test-node-1",
		ProviderID:   "kind://test-node-1",
		InstanceType: "kind-node",
		Zone:         "local",
		Region:       "local",
	}

	node, err := testImpl.CreateTestNode(ctx, nodeConfig)
	if err != nil {
		fmt.Printf("Error creating test node: %v\n", err)
	} else {
		fmt.Printf("Successfully created test node: %s\n", node.Name)
	}

	// Pattern 4: Test service creation (always supported for testing)
	fmt.Println("Creating test service...")
	serviceConfig := &testing.TestServiceConfig{
		Name:      "test-service",
		Namespace: "default",
	}

	service, err := testImpl.CreateTestService(ctx, serviceConfig)
	if err != nil {
		fmt.Printf("Error creating test service: %v\n", err)
	} else {
		fmt.Printf("Successfully created test service: %s/%s\n", service.Namespace, service.Name)
	}
}

// ExampleTestFunction demonstrates how to write a test function that handles unsupported features
func ExampleTestFunction(ti testing.TestInterface) error {
	ctx := context.Background()
	cloud := ti.GetCloudProvider()

	// Get test results for logging
	results := ti.GetTestResults()

	// Example: Test load balancer functionality with proper error handling
	results.AddLog("Testing load balancer functionality...")

	lb, supported := cloud.LoadBalancer()
	if !supported {
		results.AddLog("LoadBalancer interface is not supported by this cloud provider")
		results.AddLog("This is expected for some cloud providers - test will be skipped")
		return fmt.Errorf("load balancer not supported by cloud provider")
	}

	// Only reach this code if LoadBalancer is supported
	results.AddLog("LoadBalancer interface is supported, proceeding with test")

	// Create test resources
	node, err := ti.CreateTestNode(ctx, &testing.TestNodeConfig{
		Name:       "test-node",
		ProviderID: "test://node-1",
	})
	if err != nil {
		return fmt.Errorf("failed to create test node: %w", err)
	}

	service, err := ti.CreateTestService(ctx, &testing.TestServiceConfig{
		Name:      "test-service",
		Namespace: "default",
	})
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Test the load balancer functionality
	// Note: This might fail due to environment constraints, not interface support
	_, err = lb.EnsureLoadBalancer(ctx, "test-cluster", service, []*v1.Node{node})
	if err != nil {
		// Check if this is an environment error vs. actual test failure
		if IsEnvironmentError(err) {
			results.AddLog("Load balancer test failed due to environment constraints")
			results.AddLog("This is different from interface support - the interface works but environment doesn't support it")
			return fmt.Errorf("environment does not support load balancer creation: %w", err)
		}
		return fmt.Errorf("load balancer test failed: %w", err)
	}

	results.AddLog("Load balancer test completed successfully")
	return nil
}

// ExampleMain shows how to use the testing interface in a main function
func ExampleMain() {
	// Create cloud provider
	cloudProvider := provider.New("example-cluster", nil)

	// Create test implementation
	testImpl := NewCloudProviderTestImplementation(cloudProvider)

	// Create test runner
	runner := testing.NewTestRunner(testImpl)

	// Add test suites
	runner.AddTestSuite(CreateProviderTestSuite()) // Always supported

	// Only add load balancer tests if supported
	if testImpl.IsFeatureSupported("loadbalancer") {
		runner.AddTestSuite(CreateLoadBalancerTestSuite())
	} else {
		log.Println("Load balancer tests skipped - not supported in this environment")
	}

	// Only add cluster tests if supported
	if testImpl.IsFeatureSupported("clusters") {
		runner.AddTestSuite(CreateClustersTestSuite())
	} else {
		log.Println("Cluster tests skipped - not supported in this environment")
	}

	// Run tests
	ctx := context.Background()
	err := runner.RunTests(ctx)
	if err != nil {
		log.Printf("Test execution failed: %v", err)
		return
	}

	// Get and display results
	summary := runner.GetSummary()
	log.Printf("Test Summary: %d total, %d passed, %d failed, %d skipped",
		summary.TotalTests, summary.PassedTests, summary.FailedTests, summary.SkippedTests)
}
