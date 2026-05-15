package kubernetes

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Deployment represents a deployment in the system
type Deployment struct {
	ID        string    `json:"id"`
	Provider  string    `json:"provider"`
	Tier      string    `json:"tier"`
	Region    string    `json:"region"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetDeployments fetches deployments from Crossplane XCompute resources
func (c *Client) GetDeployments() ([]Deployment, error) {
	if c.DynamicClient == nil {
		return nil, fmt.Errorf("dynamic client not initialized")
	}

	// Try to fetch XCompute composite resources from Crossplane
	// The GVK (Group/Version/Kind) for XCompute would be defined in the XRD
	// For now, we'll try to fetch composite resources from the default Crossplane namespace
	// In production, you'd use the actual GVK from your XRD

	// Try to list composite resources - this is a common pattern for Crossplane
	// The actual GVK depends on your XRD definition
	list, err := c.DynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    "compute.example.org",
			Version:  "v1alpha1",
			Resource: "xcomputes",
		},
	).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		// If Crossplane resources don't exist, fall back to the simplified pod-based approach
		return c.getDeploymentsFromPods()
	}

	var deployments []Deployment
	for _, item := range list.Items {
		deployment := c.unstructuredToDeployment(item)
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

// getDeploymentsFromPods is the fallback simplified version
func (c *Client) getDeploymentsFromPods() ([]Deployment, error) {
	if c.K8sClient == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	namespaces, err := c.K8sClient.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var k8sDeployments []Deployment
	for _, ns := range namespaces.Items {
		pods, err := c.K8sClient.CoreV1().Pods(ns.Name).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, pod := range pods.Items {
			provider := pod.Labels["engine.io/provider"]
			tier := pod.Labels["engine.io/tier"]
			region := pod.Labels["topology.kubernetes.io/zone"]

			if provider == "" {
				provider = "unknown"
			}
			if tier == "" {
				tier = "micro"
			}
			if region == "" {
				region = "unknown"
			}

			status := "running"
			if pod.Status.Phase != "" {
				status = string(pod.Status.Phase)
			}

			k8sDeployments = append(k8sDeployments, Deployment{
				ID:        pod.Name,
				Provider:  provider,
				Tier:      tier,
				Region:    region,
				Status:    status,
				CreatedAt: pod.CreationTimestamp.Time,
			})
		}
	}

	return k8sDeployments, nil
}

// unstructuredToDeployment converts an unstructured Crossplane resource to a Deployment
func (c *Client) unstructuredToDeployment(item unstructured.Unstructured) Deployment {
	provider, _, _ := unstructured.NestedString(item.Object, "spec", "provider")
	tier, _, _ := unstructured.NestedString(item.Object, "spec", "tier")
	region, _, _ := unstructured.NestedString(item.Object, "spec", "region")

	if provider == "" {
		provider = "unknown"
	}
	if tier == "" {
		tier = "micro"
	}
	if region == "" {
		region = "unknown"
	}

	// Try to get status from the Crossplane resource
	status := "unknown"
	if statusVal, found, _ := unstructured.NestedString(item.Object, "status", "conditions", "0", "type"); found {
		status = statusVal
	}

	createdAt := time.Now()
	if ts := item.GetCreationTimestamp(); !ts.IsZero() {
		createdAt = ts.Time
	}

	return Deployment{
		ID:        item.GetName(),
		Provider:  provider,
		Tier:      tier,
		Region:    region,
		Status:    status,
		CreatedAt: createdAt,
	}
}
