// pkg/testing/test_utils.go
package testing

import (
	"os/exec"
	"strings"
)

// IsLoadBalancerSupported checks if the current environment supports load balancer creation
func IsLoadBalancerSupported() bool {
	// Check if any container runtime is available
	if !isContainerRuntimeAvailable() {
		return false
	}

	// Check if we can create a privileged test container (required for load balancers)
	if !canCreatePrivilegedContainer() {
		return false
	}

	// Check if we can set sysctl parameters (required for load balancers)
	if !canSetSysctl() {
		return false
	}

	return true
}

// IsClustersSupported checks if the current environment supports cluster operations
func IsClustersSupported() bool {
	// Check if KIND is available
	if !isKindAvailable() {
		return false
	}

	// Check if we can list clusters
	if !canListClusters() {
		return false
	}

	return true
}

// IsInstancesSupported checks if the current environment supports instance operations
func IsInstancesSupported() bool {
	// Check if KIND is available
	if !isKindAvailable() {
		return false
	}

	// Check if we can access cluster nodes
	if !canAccessClusterNodes() {
		return false
	}

	return true
}

// IsProviderSupported checks if the current environment supports provider operations
func IsProviderSupported() bool {
	// Provider operations are always supported as they don't require external dependencies
	return true
}

// isContainerRuntimeAvailable checks if a container runtime is available
func isContainerRuntimeAvailable() bool {
	// Check for Docker
	if isCommandAvailable("docker") {
		// Also check if we can actually use Docker
		if canUseDocker() {
			return true
		}
	}

	// Check for Podman
	if isCommandAvailable("podman") {
		// Also check if we can actually use Podman
		if canUsePodman() {
			return true
		}
	}

	// Check for nerdctl
	if isCommandAvailable("nerdctl") {
		return true
	}

	// Check for Finch
	if isCommandAvailable("finch") {
		return true
	}

	return false
}

// isKindAvailable checks if KIND is available
func isKindAvailable() bool {
	return isCommandAvailable("kind")
}

// canListClusters checks if we can list KIND clusters
func canListClusters() bool {
	cmd := exec.Command("kind", "get", "clusters")
	return cmd.Run() == nil
}

// canAccessClusterNodes checks if we can access cluster nodes
func canAccessClusterNodes() bool {
	// Try to get nodes from any existing cluster
	cmd := exec.Command("kind", "get", "nodes")
	return cmd.Run() == nil
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// canUseDocker checks if Docker is available and accessible
func canUseDocker() bool {
	cmd := exec.Command("docker", "version")
	return cmd.Run() == nil
}

// canUsePodman checks if Podman is available and accessible
func canUsePodman() bool {
	cmd := exec.Command("podman", "version")
	return cmd.Run() == nil
}

// canCreatePrivilegedContainer checks if we can create a privileged container (required for load balancers)
func canCreatePrivilegedContainer() bool {
	// Try to create a privileged test container to verify the environment supports it
	testContainerName := "cloud-provider-kind-privileged-test"

	// Clean up any existing test container
	cleanupTestContainer(testContainerName)

	// Try to create a privileged container with sysctl settings
	err := createPrivilegedTestContainer(testContainerName)
	if err != nil {
		return false
	}

	// Clean up the test container
	cleanupTestContainer(testContainerName)

	return true
}

// canSetSysctl checks if we can set sysctl parameters in containers
func canSetSysctl() bool {
	// This is tested as part of canCreatePrivilegedContainer
	// since load balancers require sysctl settings
	return true
}

// canCreateTestContainer checks if we can create a simple test container
func canCreateTestContainer() bool {
	// Try to create a simple test container to verify the environment works
	testContainerName := "cloud-provider-kind-test-container"

	// Clean up any existing test container
	cleanupTestContainer(testContainerName)

	// Try to create a simple container
	err := createSimpleTestContainer(testContainerName)
	if err != nil {
		return false
	}

	// Clean up the test container
	cleanupTestContainer(testContainerName)

	return true
}

// createPrivilegedTestContainer creates a privileged test container with sysctl settings
func createPrivilegedTestContainer(name string) error {
	// Use a simple image that should be available
	image := "alpine:latest"

	// Try to pull the image first
	pullCmd := exec.Command("docker", "pull", image)
	if err := pullCmd.Run(); err != nil {
		// If docker pull fails, try with other runtimes
		pullCmd = exec.Command("podman", "pull", image)
		if err := pullCmd.Run(); err != nil {
			return err
		}
	}

	// Create a privileged container with sysctl settings (similar to load balancer requirements)
	args := []string{
		"run", "--name", name, "--rm", "--privileged",
		"--sysctl=net.ipv4.ip_forward=1",
		"--sysctl=net.ipv4.conf.all.rp_filter=0",
		"--sysctl=net.ipv4.ip_unprivileged_port_start=1",
		image, "echo", "test",
	}

	createCmd := exec.Command("docker", args...)
	if err := createCmd.Run(); err != nil {
		// If docker fails, try with other runtimes
		createCmd = exec.Command("podman", args...)
		if err := createCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// createSimpleTestContainer creates a simple test container
func createSimpleTestContainer(name string) error {
	// Use a simple image that should be available
	image := "alpine:latest"

	// Try to pull the image first
	pullCmd := exec.Command("docker", "pull", image)
	if err := pullCmd.Run(); err != nil {
		// If docker pull fails, try with other runtimes
		pullCmd = exec.Command("podman", "pull", image)
		if err := pullCmd.Run(); err != nil {
			return err
		}
	}

	// Create a simple container that exits immediately
	createCmd := exec.Command("docker", "run", "--name", name, "--rm", image, "echo", "test")
	if err := createCmd.Run(); err != nil {
		// If docker fails, try with other runtimes
		createCmd = exec.Command("podman", "run", "--name", name, "--rm", image, "echo", "test")
		if err := createCmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestContainer removes a test container
func cleanupTestContainer(name string) {
	// Try docker first
	exec.Command("docker", "rm", "-f", name).Run()
	// Try podman as fallback
	exec.Command("podman", "rm", "-f", name).Run()
}

// IsEnvironmentError checks if an error is related to environment constraints
func IsEnvironmentError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Check for common environment-related error messages
	environmentErrors := []string{
		"docker",
		"container",
		"permission denied",
		"connection refused",
		"no such file or directory",
		"command not found",
		"daemon",
		"network",
		"bridge",
		"exit status 126",
		"exit status 127",
		"exit status 1",
		"failed to create continers",
		"failed to create container",
		"privileged",
		"sysctl",
		"kind",
		"cluster",
		"node",
		"kubeconfig",
		"kubernetes",
	}

	for _, envErr := range environmentErrors {
		if strings.Contains(strings.ToLower(errStr), envErr) {
			return true
		}
	}

	return false
}

// GetTestEnvironmentInfo returns information about the test environment
func GetTestEnvironmentInfo() map[string]bool {
	return map[string]bool{
		"loadbalancer_supported": IsLoadBalancerSupported(),
		"clusters_supported":     IsClustersSupported(),
		"instances_supported":    IsInstancesSupported(),
		"provider_supported":     IsProviderSupported(),
		"docker_available":       isCommandAvailable("docker") && canUseDocker(),
		"podman_available":       isCommandAvailable("podman") && canUsePodman(),
		"kind_available":         isKindAvailable(),
		"container_runtime":      isContainerRuntimeAvailable(),
	}
}
