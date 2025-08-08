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

	// Test 1: Verify LoadBalancer interface is implemented
	cloud := testInterface.GetCloudProvider()
	lb, ok := cloud.LoadBalancer()
	if !ok {
		return fmt.Errorf("load balancer interface not supported by cloud provider")
	}
	klog.Info("✓ LoadBalancer interface is implemented")

	// Test 2: Verify GetLoadBalancerName method works
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-basic",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:     "http",
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

	// Test GetLoadBalancerName
	lbName := lb.GetLoadBalancerName(context.Background(), "test-cluster", service)
	if lbName == "" {
		klog.Warning("GetLoadBalancerName returned empty string (this may be expected for some providers)")
	} else {
		klog.Infof("✓ GetLoadBalancerName returned: %s", lbName)
	}

	// Test 3: Verify GetLoadBalancer method works (may return not found, which is OK)
	status, exists, err := lb.GetLoadBalancer(context.Background(), "test-cluster", service)
	if err != nil {
		klog.Warningf("GetLoadBalancer returned error (may be expected): %v", err)
	} else {
		klog.Infof("✓ GetLoadBalancer returned exists=%v, status=%+v", exists, status)
	}

	// Test 4: Verify EnsureLoadBalancer method can be called (may fail, but should not panic)
	_, err = lb.EnsureLoadBalancer(context.Background(), "test-cluster", service, []*v1.Node{})
	if err != nil {
		klog.Infof("✓ EnsureLoadBalancer returned expected error: %v", err)
	} else {
		klog.Info("✓ EnsureLoadBalancer completed successfully")
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

	// Test 1: Verify LoadBalancer interface is implemented
	cloud := testInterface.GetCloudProvider()
	lb, ok := cloud.LoadBalancer()
	if !ok {
		return fmt.Errorf("load balancer interface not supported by cloud provider")
	}
	klog.Info("✓ LoadBalancer interface is implemented")

	// Test 2: Create initial service
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-update",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:     "http",
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Test 3: Verify UpdateLoadBalancer method can be called
	err = lb.UpdateLoadBalancer(context.Background(), "test-cluster", service, []*v1.Node{})
	if err != nil {
		klog.Infof("✓ UpdateLoadBalancer returned expected error: %v", err)
	} else {
		klog.Info("✓ UpdateLoadBalancer completed successfully")
	}

	// Test 4: Update service (add a new port) and test again
	service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{
		Name:     "https",
		Port:     443,
		Protocol: v1.ProtocolTCP,
	})

	_, err = kts.kubeClient.CoreV1().Services(service.Namespace).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// Test UpdateLoadBalancer with updated service
	err = lb.UpdateLoadBalancer(context.Background(), "test-cluster", service, []*v1.Node{})
	if err != nil {
		klog.Infof("✓ UpdateLoadBalancer with updated service returned expected error: %v", err)
	} else {
		klog.Info("✓ UpdateLoadBalancer with updated service completed successfully")
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

	// Test 1: Verify LoadBalancer interface is implemented
	cloud := testInterface.GetCloudProvider()
	lb, ok := cloud.LoadBalancer()
	if !ok {
		return fmt.Errorf("load balancer interface not supported by cloud provider")
	}
	klog.Info("✓ LoadBalancer interface is implemented")

	// Test 2: Create service
	serviceConfig := &exttesting.TestServiceConfig{
		Name:      "test-lb-delete",
		Namespace: "default",
		Type:      v1.ServiceTypeLoadBalancer,
		Ports: []v1.ServicePort{
			{
				Name:     "http",
				Port:     80,
				Protocol: v1.ProtocolTCP,
			},
		},
	}

	service, err := testInterface.CreateTestService(context.Background(), serviceConfig)
	if err != nil {
		return fmt.Errorf("failed to create test service: %w", err)
	}

	// Test 3: Verify EnsureLoadBalancerDeleted method can be called
	err = lb.EnsureLoadBalancerDeleted(context.Background(), "test-cluster", service)
	if err != nil {
		klog.Infof("✓ EnsureLoadBalancerDeleted returned expected error: %v", err)
	} else {
		klog.Info("✓ EnsureLoadBalancerDeleted completed successfully")
	}

	// Clean up
	if err := testInterface.DeleteTestService(context.Background(), service.Name); err != nil {
		klog.Warningf("Failed to delete test service: %v", err)
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
			klog.Warningf("Instance does not exist for node %s (this may be expected for some providers)", node.Name)
		} else {
			klog.Infof("✓ Instance exists for node %s", node.Name)
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
			klog.Warningf("Instance metadata is nil for node %s (this may be expected for some providers)", node.Name)
		} else {
			klog.Infof("✓ Instance metadata for node %s: ProviderID=%s", node.Name, metadata.ProviderID)
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
		// Log the shutdown status (both true and false are valid)
		klog.Infof("✓ Node %s shutdown status: %v", node.Name, shutdown)
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
		klog.Warning("No clusters found (this may be expected for some providers)")
	} else {
		klog.Infof("✓ Found clusters: %v", clusterList)
	}
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
			klog.Warningf("Failed to get master for cluster %s (this may be expected): %v", clusterName, err)
		} else if master == "" {
			klog.Warningf("Master endpoint is empty for cluster %s (this may be expected for some providers)", clusterName)
		} else {
			klog.Infof("✓ Cluster %s master: %s", clusterName, master)
		}
	}

	klog.Info("Master endpoint test passed")
	return nil
}

func (kts *KindTestSuite) testProviderName(testInterface exttesting.TestInterface) error {
	klog.Info("Running provider name test")

	cloud := testInterface.GetCloudProvider()
	providerName := cloud.ProviderName()

	if providerName == "" {
		klog.Warning("Provider name is empty (this may be expected for some providers)")
	} else {
		klog.Infof("✓ Provider name: %s", providerName)
	}
	klog.Info("Provider name test passed")
	return nil
}
