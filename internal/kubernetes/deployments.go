package kubernetes

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// GetDeployments fetches deployments from Kubernetes
// This is a simplified version - in production you'd use the Crossplane CRD client
func (c *Client) GetDeployments() ([]Deployment, error) {
	if c.K8sClient == nil {
		return nil, fmt.Errorf("kubernetes client not initialized")
	}

	// Try to get XCompute resources from Crossplane namespace
	// This is a simplified version - in production you'd use the Crossplane CRD client
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
