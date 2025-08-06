package testing

import (
	"context"
	"fmt"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cloud-provider-kind/pkg/provider"
	"sigs.k8s.io/kind/pkg/cluster"

	// Import the external testing interface
	exttesting "github.com/miyadav/cloud-provider-testing-interface"
)

// KindTestInterface implements the external TestInterface for KIND cloud provider
type KindTestInterface struct {
	cloudProvider cloudprovider.Interface
	kindProvider  *cluster.Provider
	clusterName   string
	kubeClient    kubernetes.Interface
	config        *exttesting.TestConfig
	results       *exttesting.TestResults
	mu            sync.RWMutex
}

// NewKindTestInterface creates a new KIND test interface
func NewKindTestInterface(clusterName string, kindProvider *cluster.Provider, kubeClient kubernetes.Interface) *KindTestInterface {
	return &KindTestInterface{
		kindProvider: kindProvider,
		clusterName:  clusterName,
		kubeClient:   kubeClient,
		results:      &exttesting.TestResults{},
	}
}

// SetupTestEnvironment initializes the test environment
func (kti *KindTestInterface) SetupTestEnvironment(config *exttesting.TestConfig) error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	kti.config = config

	// Create cloud provider instance
	kti.cloudProvider = provider.New(kti.clusterName, kti.kindProvider)

	// Initialize the cloud provider
	kti.cloudProvider.Initialize(config.ClientBuilder, make(chan struct{}))

	klog.Infof("Test environment setup completed for cluster: %s", kti.clusterName)
	kti.results.AddLog(fmt.Sprintf("Test environment setup completed for cluster: %s", kti.clusterName))

	return nil
}

// TeardownTestEnvironment cleans up the test environment
func (kti *KindTestInterface) TeardownTestEnvironment() error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	// Clean up any test resources
	if kti.config != nil && kti.config.CleanupResources {
		// Delete test namespaces, services, etc.
		klog.Info("Cleaning up test resources")
		kti.results.AddLog("Test environment teardown completed")
	}

	return nil
}

// GetCloudProvider returns the cloud provider instance
func (kti *KindTestInterface) GetCloudProvider() cloudprovider.Interface {
	kti.mu.RLock()
	defer kti.mu.RUnlock()
	return kti.cloudProvider
}

// CreateTestNode creates a test node
func (kti *KindTestInterface) CreateTestNode(ctx context.Context, nodeConfig *exttesting.TestNodeConfig) (*v1.Node, error) {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	node := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:        nodeConfig.Name,
			Labels:      nodeConfig.Labels,
			Annotations: nodeConfig.Annotations,
		},
		Spec: v1.NodeSpec{
			ProviderID: nodeConfig.ProviderID,
		},
		Status: v1.NodeStatus{
			Addresses:  nodeConfig.Addresses,
			Conditions: nodeConfig.Conditions,
			NodeInfo: v1.NodeSystemInfo{
				KubeletVersion: "v1.24.0",
			},
		},
	}

	createdNode, err := kti.kubeClient.CoreV1().Nodes().Create(ctx, node, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test node: %w", err)
	}

	kti.results.IncrementResourceCount("nodes")
	kti.results.AddLog(fmt.Sprintf("Created test node: %s", nodeConfig.Name))

	return createdNode, nil
}

// DeleteTestNode deletes a test node
func (kti *KindTestInterface) DeleteTestNode(ctx context.Context, nodeName string) error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	err := kti.kubeClient.CoreV1().Nodes().Delete(ctx, nodeName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete test node: %w", err)
	}

	kti.results.AddLog(fmt.Sprintf("Deleted test node: %s", nodeName))
	return nil
}

// CreateTestService creates a test service
func (kti *KindTestInterface) CreateTestService(ctx context.Context, serviceConfig *exttesting.TestServiceConfig) (*v1.Service, error) {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        serviceConfig.Name,
			Namespace:   serviceConfig.Namespace,
			Labels:      serviceConfig.Labels,
			Annotations: serviceConfig.Annotations,
		},
		Spec: v1.ServiceSpec{
			Type:                  serviceConfig.Type,
			Ports:                 serviceConfig.Ports,
			LoadBalancerIP:        serviceConfig.LoadBalancerIP,
			ExternalTrafficPolicy: serviceConfig.ExternalTrafficPolicy,
			InternalTrafficPolicy: serviceConfig.InternalTrafficPolicy,
			Selector: map[string]string{
				"app": "test-app",
			},
		},
	}

	createdService, err := kti.kubeClient.CoreV1().Services(serviceConfig.Namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create test service: %w", err)
	}

	kti.results.IncrementResourceCount("services")
	kti.results.AddLog(fmt.Sprintf("Created test service: %s/%s", serviceConfig.Namespace, serviceConfig.Name))

	return createdService, nil
}

// DeleteTestService deletes a test service
func (kti *KindTestInterface) DeleteTestService(ctx context.Context, serviceName string) error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	// Try to delete from default namespace first
	err := kti.kubeClient.CoreV1().Services("default").Delete(ctx, serviceName, metav1.DeleteOptions{})
	if err != nil {
		// Try other common namespaces
		namespaces := []string{"kube-system", "cloud-provider-test"}
		for _, ns := range namespaces {
			if deleteErr := kti.kubeClient.CoreV1().Services(ns).Delete(ctx, serviceName, metav1.DeleteOptions{}); deleteErr == nil {
				kti.results.AddLog(fmt.Sprintf("Deleted test service: %s/%s", ns, serviceName))
				return nil
			}
		}
		return fmt.Errorf("failed to delete test service: %w", err)
	}

	kti.results.AddLog(fmt.Sprintf("Deleted test service: %s", serviceName))
	return nil
}

// CreateTestRoute creates a test route
func (kti *KindTestInterface) CreateTestRoute(ctx context.Context, routeConfig *exttesting.TestRouteConfig) (*cloudprovider.Route, error) {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	route := &cloudprovider.Route{
		Name:            routeConfig.Name,
		TargetNode:      routeConfig.TargetNode,
		DestinationCIDR: routeConfig.DestinationCIDR,
		Blackhole:       routeConfig.Blackhole,
	}

	kti.results.IncrementResourceCount("routes")
	kti.results.AddLog(fmt.Sprintf("Created test route: %s", routeConfig.Name))

	return route, nil
}

// DeleteTestRoute deletes a test route
func (kti *KindTestInterface) DeleteTestRoute(ctx context.Context, routeName string) error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	kti.results.AddLog(fmt.Sprintf("Deleted test route: %s", routeName))
	return nil
}

// WaitForCondition waits for a condition to be met
func (kti *KindTestInterface) WaitForCondition(ctx context.Context, condition exttesting.TestCondition) error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	timeout := condition.Timeout
	if timeout == 0 {
		timeout = 2 * time.Minute
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if condition.CheckFunction != nil {
				met, err := condition.CheckFunction()
				if err != nil {
					klog.V(2).Infof("Condition check failed: %v", err)
					continue
				}
				if met {
					kti.results.AddLog(fmt.Sprintf("Condition met: %s", condition.Type))
					return nil
				}
			}
		}
	}

	return fmt.Errorf("timeout waiting for condition: %s", condition.Type)
}

// GetTestResults returns the test results
func (kti *KindTestInterface) GetTestResults() *exttesting.TestResults {
	kti.mu.RLock()
	defer kti.mu.RUnlock()
	return kti.results
}

// ResetTestState resets the test state
func (kti *KindTestInterface) ResetTestState() error {
	kti.mu.Lock()
	defer kti.mu.Unlock()

	kti.results = &exttesting.TestResults{}
	kti.results.AddLog("Test state reset")

	return nil
}
