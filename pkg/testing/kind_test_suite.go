package testing

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"sigs.k8s.io/kind/pkg/cluster"

	// Import the external testing interface
	exttesting "github.com/miyadav/cloud-provider-testing-interface"
)

// KindTestSuite provides a comprehensive test suite for KIND cloud provider
type KindTestSuite struct {
	testInterface exttesting.TestInterface
	kubeClient    kubernetes.Interface
}

// NewKindTestSuite creates a new KIND test suite
func NewKindTestSuite(clusterName string, kindProvider *cluster.Provider, kubeClient kubernetes.Interface) *KindTestSuite {
	testInterface := NewKindTestInterface(clusterName, kindProvider, kubeClient)
	return &KindTestSuite{
		testInterface: testInterface,
		kubeClient:    kubeClient,
	}
}

// GetLoadBalancerTestSuite returns a test suite for load balancer functionality
func (kts *KindTestSuite) GetLoadBalancerTestSuite() exttesting.TestSuite {
	return exttesting.TestSuite{
		Name:        "Load Balancer Tests",
		Description: "Tests for load balancer functionality",
		Setup:       kts.setupLoadBalancerTests,
		Teardown:    kts.teardownLoadBalancerTests,
		Tests: []exttesting.Test{
			{
				Name:        "Basic Load Balancer Creation",
				Description: "Tests basic load balancer creation and provisioning",
				Run:         kts.testBasicLoadBalancerCreation,
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "Load Balancer Update",
				Description: "Tests load balancer update functionality",
				Run:         kts.testLoadBalancerUpdate,
				Timeout:     3 * time.Minute,
			},
			{
				Name:        "Load Balancer Deletion",
				Description: "Tests load balancer deletion and cleanup",
				Run:         kts.testLoadBalancerDeletion,
				Timeout:     2 * time.Minute,
			},
		},
	}
}

// GetInstancesTestSuite returns a test suite for instance functionality
func (kts *KindTestSuite) GetInstancesTestSuite() exttesting.TestSuite {
	return exttesting.TestSuite{
		Name:        "Instance Tests",
		Description: "Tests for instance management functionality",
		Setup:       kts.setupInstanceTests,
		Teardown:    kts.teardownInstanceTests,
		Tests: []exttesting.Test{
			{
				Name:        "Instance Existence",
				Description: "Tests instance existence verification",
				Run:         kts.testInstanceExistence,
				Timeout:     1 * time.Minute,
			},
			{
				Name:        "Instance Metadata",
				Description: "Tests instance metadata retrieval",
				Run:         kts.testInstanceMetadata,
				Timeout:     1 * time.Minute,
			},
			{
				Name:        "Instance Shutdown Status",
				Description: "Tests instance shutdown status verification",
				Run:         kts.testInstanceShutdownStatus,
				Timeout:     1 * time.Minute,
			},
		},
	}
}

// GetClustersTestSuite returns a test suite for cluster functionality
func (kts *KindTestSuite) GetClustersTestSuite() exttesting.TestSuite {
	return exttesting.TestSuite{
		Name:        "Cluster Tests",
		Description: "Tests for cluster management functionality",
		Setup:       kts.setupClusterTests,
		Teardown:    kts.teardownClusterTests,
		Tests: []exttesting.Test{
			{
				Name:        "Cluster Listing",
				Description: "Tests cluster listing functionality",
				Run:         kts.testClusterListing,
				Timeout:     1 * time.Minute,
			},
			{
				Name:        "Master Endpoint",
				Description: "Tests master endpoint retrieval",
				Run:         kts.testMasterEndpoint,
				Timeout:     1 * time.Minute,
			},
		},
	}
}

// GetProviderNameTestSuite returns a test suite for provider name functionality
func (kts *KindTestSuite) GetProviderNameTestSuite() exttesting.TestSuite {
	return exttesting.TestSuite{
		Name:        "Provider Name Tests",
		Description: "Tests for provider name functionality",
		Setup:       kts.setupProviderNameTests,
		Teardown:    kts.teardownProviderNameTests,
		Tests: []exttesting.Test{
			{
				Name:        "Provider Name Verification",
				Description: "Tests provider name verification",
				Run:         kts.testProviderName,
				Timeout:     30 * time.Second,
			},
		},
	}
}

// Setup functions for test suites
func (kts *KindTestSuite) setupLoadBalancerTests(testInterface exttesting.TestInterface) error {
	klog.Info("Setting up load balancer tests")
	return nil
}

func (kts *KindTestSuite) teardownLoadBalancerTests(testInterface exttesting.TestInterface) error {
	klog.Info("Tearing down load balancer tests")
	return nil
}

func (kts *KindTestSuite) setupInstanceTests(testInterface exttesting.TestInterface) error {
	klog.Info("Setting up instance tests")
	return nil
}

func (kts *KindTestSuite) teardownInstanceTests(testInterface exttesting.TestInterface) error {
	klog.Info("Tearing down instance tests")
	return nil
}

func (kts *KindTestSuite) setupClusterTests(testInterface exttesting.TestInterface) error {
	klog.Info("Setting up cluster tests")
	return nil
}

func (kts *KindTestSuite) teardownClusterTests(testInterface exttesting.TestInterface) error {
	klog.Info("Tearing down cluster tests")
	return nil
}

func (kts *KindTestSuite) setupProviderNameTests(testInterface exttesting.TestInterface) error {
	klog.Info("Setting up provider name tests")
	return nil
}

func (kts *KindTestSuite) teardownProviderNameTests(testInterface exttesting.TestInterface) error {
	klog.Info("Tearing down provider name tests")
	return nil
}

// Test implementations
func (kts *KindTestSuite) testBasicLoadBalancerCreation(testInterface exttesting.TestInterface) error {
	klog.Info("Running basic load balancer creation test")

	// Create a test service
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-basic",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
		Labels: map[string]string{
			"app": "test-app",
		},
	}

	service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Wait for load balancer to be provisioned
	condition := exttesting.TestCondition{
		Type:    "LoadBalancerReady",
		Timeout: 2 * time.Minute,
		CheckFunction: func() (bool, error) {
			cloud := testInterface.GetCloudProvider()
			lb, ok := cloud.LoadBalancer()
			if !ok {
				return false, fmt.Errorf("load balancer interface not supported")
			}

			status, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
			if err != nil {
				return false, nil // Continue waiting
			}
			return exists && len(status.Ingress) > 0, nil
		},
	}

	if err := testInterface.WaitForCondition(context.Background(), condition); err != nil {
		return fmt.Errorf("load balancer not ready: %w", err)
	}

	// Clean up
	if err := testInterface.DeleteTestService(context.Background(), service.Name); err != nil {
		klog.Warningf("Failed to delete test service: %v", err)
	}

	klog.Info("Basic load balancer creation test passed")
	return nil
}

func (kts *KindTestSuite) testLoadBalancerUpdate(testInterface exttesting.TestInterface) error {
	klog.Info("Running load balancer update test")

	// Create initial service
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-update",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Wait for initial load balancer
	condition := exttesting.TestCondition{
		Type:    "LoadBalancerReady",
		Timeout: 2 * time.Minute,
		CheckFunction: func() (bool, error) {
			cloud := testInterface.GetCloudProvider()
			lb, ok := cloud.LoadBalancer()
			if !ok {
				return false, fmt.Errorf("load balancer interface not supported")
			}

			status, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
			if err != nil {
				return false, nil
			}
			return exists && len(status.Ingress) > 0, nil
		},
	}

	if err := testInterface.WaitForCondition(context.Background(), condition); err != nil {
		return fmt.Errorf("initial load balancer not ready: %w", err)
	}

	// Update service (add a new port)
	service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{
		Port:     443,
		Protocol: v1.ProtocolTCP,
	})

	_, err = kts.kubeClient.CoreV1().Services(service.Namespace).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// Wait for update to be processed
	updateCondition := exttesting.TestCondition{
		Type:    "LoadBalancerUpdated",
		Timeout: 1 * time.Minute,
		CheckFunction: func() (bool, error) {
			cloud := testInterface.GetCloudProvider()
			lb, ok := cloud.LoadBalancer()
			if !ok {
				return false, fmt.Errorf("load balancer interface not supported")
			}

			_, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
			return exists && err == nil, nil
		},
	}

	if err := testInterface.WaitForCondition(context.Background(), updateCondition); err != nil {
		return fmt.Errorf("load balancer update not processed: %w", err)
	}

	// Clean up
	if err := testInterface.DeleteTestService(context.Background(), service.Name); err != nil {
		klog.Warningf("Failed to delete test service: %v", err)
	}

	klog.Info("Load balancer update test passed")
	return nil
}

func (kts *KindTestSuite) testLoadBalancerDeletion(testInterface exttesting.TestInterface) error {
	klog.Info("Running load balancer deletion test")

	// Create service
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-delete",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Wait for load balancer to be created
	condition := exttesting.TestCondition{
		Type:    "LoadBalancerReady",
		Timeout: 2 * time.Minute,
		CheckFunction: func() (bool, error) {
			cloud := testInterface.GetCloudProvider()
			lb, ok := cloud.LoadBalancer()
			if !ok {
				return false, fmt.Errorf("load balancer interface not supported")
			}

			status, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
			if err != nil {
				return false, nil
			}
			return exists && len(status.Ingress) > 0, nil
		},
	}

	if err := testInterface.WaitForCondition(context.Background(), condition); err != nil {
		return fmt.Errorf("load balancer not ready: %w", err)
	}

	// Delete service
	if err := testInterface.DeleteTestService(context.Background(), service.Name); err != nil {
		return fmt.Errorf("failed to delete test service: %w", err)
	}

	// Wait for load balancer to be deleted
	deleteCondition := exttesting.TestCondition{
		Type:    "LoadBalancerDeleted",
		Timeout: 1 * time.Minute,
		CheckFunction: func() (bool, error) {
			cloud := testInterface.GetCloudProvider()
			lb, ok := cloud.LoadBalancer()
			if !ok {
				return false, fmt.Errorf("load balancer interface not supported")
			}

			_, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
			return !exists && err == nil, nil
		},
	}

	if err := testInterface.WaitForCondition(context.Background(), deleteCondition); err != nil {
		return fmt.Errorf("load balancer not deleted: %w", err)
	}

	klog.Info("Load balancer deletion test passed")
	return nil
}

func (kts *KindTestSuite) testInstanceExistence(testInterface exttesting.TestInterface) error {
	klog.Info("Running instance existence test")

	// Get nodes
	nodes, err := kts.kubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	if len(nodes.Items) == 0 {
		return fmt.Errorf("no nodes found")
	}

	// Test instance exists for each node
	cloud := testInterface.GetCloudProvider()
	instances, ok := cloud.InstancesV2()
	if !ok {
		return fmt.Errorf("instances v2 interface not supported")
	}

	for _, node := range nodes.Items {
		exists, err := instances.InstanceExists(context.Background(), &node)
		if err != nil {
			return fmt.Errorf("failed to check if instance exists for node %s: %w", node.Name, err)
		}
		if !exists {
			return fmt.Errorf("instance does not exist for node %s", node.Name)
		}
	}

	klog.Info("Instance existence test passed")
	return nil
}

func (kts *KindTestSuite) testInstanceMetadata(testInterface exttesting.TestInterface) error {
	klog.Info("Running instance metadata test")

	// Get nodes
	nodes, err := kts.kubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	if len(nodes.Items) == 0 {
		return fmt.Errorf("no nodes found")
	}

	// Test instance metadata for each node
	cloud := testInterface.GetCloudProvider()
	instances, ok := cloud.InstancesV2()
	if !ok {
		return fmt.Errorf("instances v2 interface not supported")
	}

	for _, node := range nodes.Items {
		metadata, err := instances.InstanceMetadata(context.Background(), &node)
		if err != nil {
			return fmt.Errorf("failed to get instance metadata for node %s: %w", node.Name, err)
		}
		if metadata == nil {
			return fmt.Errorf("instance metadata is nil for node %s", node.Name)
		}
		if metadata.ProviderID == "" {
			return fmt.Errorf("provider ID is empty for node %s", node.Name)
		}
	}

	klog.Info("Instance metadata test passed")
	return nil
}

func (kts *KindTestSuite) testInstanceShutdownStatus(testInterface exttesting.TestInterface) error {
	klog.Info("Running instance shutdown status test")

	// Get nodes
	nodes, err := kts.kubeClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	if len(nodes.Items) == 0 {
		return fmt.Errorf("no nodes found")
	}

	// Test instance shutdown status for each node
	cloud := testInterface.GetCloudProvider()
	instances, ok := cloud.InstancesV2()
	if !ok {
		return fmt.Errorf("instances v2 interface not supported")
	}

	for _, node := range nodes.Items {
		shutdown, err := instances.InstanceShutdown(context.Background(), &node)
		if err != nil {
			return fmt.Errorf("failed to check shutdown status for node %s: %w", node.Name, err)
		}
		// For KIND, nodes should not be shutdown
		if shutdown {
			return fmt.Errorf("node %s is unexpectedly shutdown", node.Name)
		}
	}

	klog.Info("Instance shutdown status test passed")
	return nil
}

func (kts *KindTestSuite) testClusterListing(testInterface exttesting.TestInterface) error {
	klog.Info("Running cluster listing test")

	cloud := testInterface.GetCloudProvider()
	clusters, ok := cloud.Clusters()
	if !ok {
		return fmt.Errorf("clusters interface not supported")
	}

	clusterList, err := clusters.ListClusters(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	if len(clusterList) == 0 {
		return fmt.Errorf("no clusters found")
	}

	klog.Infof("Found clusters: %v", clusterList)
	klog.Info("Cluster listing test passed")
	return nil
}

func (kts *KindTestSuite) testMasterEndpoint(testInterface exttesting.TestInterface) error {
	klog.Info("Running master endpoint test")

	cloud := testInterface.GetCloudProvider()
	clusters, ok := cloud.Clusters()
	if !ok {
		return fmt.Errorf("clusters interface not supported")
	}

	clusterList, err := clusters.ListClusters(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	for _, clusterName := range clusterList {
		master, err := clusters.Master(context.Background(), clusterName)
		if err != nil {
			return fmt.Errorf("failed to get master for cluster %s: %w", clusterName, err)
		}
		if master == "" {
			return fmt.Errorf("master endpoint is empty for cluster %s", clusterName)
		}
		klog.Infof("Cluster %s master: %s", clusterName, master)
	}

	klog.Info("Master endpoint test passed")
	return nil
}

func (kts *KindTestSuite) testProviderName(testInterface exttesting.TestInterface) error {
	klog.Info("Running provider name test")

	cloud := testInterface.GetCloudProvider()
	providerName := cloud.ProviderName()

	if providerName == "" {
		return fmt.Errorf("provider name is empty")
	}

	klog.Infof("Provider name: %s", providerName)
	klog.Info("Provider name test passed")
	return nil
}
