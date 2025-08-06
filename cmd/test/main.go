package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cloud-provider-kind/pkg/testing"
	"sigs.k8s.io/kind/pkg/cluster"

	// Import the external testing interface
	exttesting "github.com/miyadav/cloud-provider-testing-interface"
)

func main() {
	var (
		kubeconfig  = flag.String("kubeconfig", "", "Path to kubeconfig file (defaults to ~/.kube/config)")
		clusterName = flag.String("cluster", "kind", "Name of the KIND cluster")
		suiteName   = flag.String("suite", "", "Specific test suite to run (optional)")
		help        = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		fmt.Printf("Cloud Provider KIND Testing Tool\n\n")
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
		fmt.Printf("\nAvailable test suites:\n")
		fmt.Printf("  loadbalancer  - Tests load balancer creation, updates, and deletion\n")
		fmt.Printf("  instances     - Tests instance existence, metadata, and shutdown status\n")
		fmt.Printf("  clusters      - Tests cluster listing and master endpoint\n")
		fmt.Printf("  provider-name - Tests provider name functionality\n\n")
		fmt.Printf("Examples:\n")
		fmt.Printf("  %s                                    # Run all tests using default kubeconfig\n", os.Args[0])
		fmt.Printf("  %s -suite loadbalancer               # Run only load balancer tests\n", os.Args[0])
		fmt.Printf("  %s -kubeconfig /path/to/config       # Use specific kubeconfig\n", os.Args[0])
		fmt.Printf("  %s -cluster my-cluster               # Test specific KIND cluster\n", os.Args[0])
		os.Exit(0)
	}

	// Create Kubernetes client
	fmt.Printf("Connecting to Kubernetes cluster...\n")
	if *kubeconfig != "" {
		fmt.Printf("Using kubeconfig: %s\n", *kubeconfig)
	} else {
		fmt.Printf("Using default kubeconfig: ~/.kube/config\n")
	}

	kubeClient, err := createKubeClient(*kubeconfig)
	if err != nil {
		klog.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	fmt.Printf("Successfully connected to Kubernetes cluster\n")

	// Initialize KIND provider
	fmt.Printf("Testing KIND cluster: %s\n", *clusterName)
	kindProvider := cluster.NewProvider()

	// Create test suite
	testSuite := testing.NewKindTestSuite(*clusterName, kindProvider, kubeClient)

	// Create test runner using external interface
	testInterface := testing.NewKindTestInterface(*clusterName, kindProvider, kubeClient)
	runner := exttesting.NewTestRunner(testInterface)

	// Add test suites based on command line arguments
	if *suiteName != "" {
		// Run specific test suite
		fmt.Printf("Running test suite: %s\n", *suiteName)
		switch *suiteName {
		case "loadbalancer":
			runner.AddTestSuite(testSuite.GetLoadBalancerTestSuite())
		case "instances":
			runner.AddTestSuite(testSuite.GetInstancesTestSuite())
		case "clusters":
			runner.AddTestSuite(testSuite.GetClustersTestSuite())
		case "provider-name":
			runner.AddTestSuite(testSuite.GetProviderNameTestSuite())
		default:
			klog.Fatalf("Unknown test suite: %s. Available suites: loadbalancer, instances, clusters, provider-name", *suiteName)
		}
	} else {
		// Run all test suites
		fmt.Printf("Running all test suites: loadbalancer, instances, clusters, provider-name\n")
		runner.AddTestSuite(testSuite.GetLoadBalancerTestSuite())
		runner.AddTestSuite(testSuite.GetInstancesTestSuite())
		runner.AddTestSuite(testSuite.GetClustersTestSuite())
		runner.AddTestSuite(testSuite.GetProviderNameTestSuite())
	}

	// Setup test environment
	config := &exttesting.TestConfig{
		ProviderName:         "kind",
		ClusterName:          *clusterName,
		TestTimeout:          10 * time.Minute,
		CleanupResources:     true,
		MockExternalServices: false,
	}

	if err := testInterface.SetupTestEnvironment(config); err != nil {
		klog.Fatalf("Failed to setup test environment: %v", err)
	}

	// Run tests
	ctx := context.Background()
	if err := runner.RunTests(ctx); err != nil {
		klog.Errorf("Tests failed: %v", err)

		// Print test results
		results := runner.GetResults()
		for _, result := range results {
			if !result.Success {
				klog.Errorf("Test %s failed: %v", result.Test.Name, result.Error)
			}
		}

		os.Exit(1)
	}

	// Print test summary
	summary := runner.GetSummary()
	fmt.Printf("Test Summary:\n")
	fmt.Printf("  Total Tests: %d\n", summary.TotalTests)
	fmt.Printf("  Passed: %d\n", summary.PassedTests)
	fmt.Printf("  Failed: %d\n", summary.FailedTests)
	fmt.Printf("  Skipped: %d\n", summary.SkippedTests)
	fmt.Printf("  Total Duration: %v\n", summary.TotalDuration)

	// Cleanup
	if err := testInterface.TeardownTestEnvironment(); err != nil {
		klog.Warningf("Failed to teardown test environment: %v", err)
	}

	fmt.Println("All tests passed successfully!")
}

func createKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	var config *rest.Config
	var err error

	if kubeconfig != "" {
		// Use provided kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		// Use default kubeconfig location (~/.kube/config)
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
		config, err = clientConfig.ClientConfig()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return kubeClient, nil
}
