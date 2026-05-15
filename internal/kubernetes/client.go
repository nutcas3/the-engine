package kubernetes

import (
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Client holds Kubernetes clients
type Client struct {
	K8sClient     *kubernetes.Clientset
	DynamicClient dynamic.Interface
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		K8sClient:     clientset,
		DynamicClient: dynamicClient,
	}, nil
}

// NewClientOrMock creates a Kubernetes client or returns nil if it fails
func NewClientOrMock() *Client {
	client, err := NewClient()
	if err != nil {
		log.Printf("Warning: Could not connect to Kubernetes: %v", err)
		return nil
	}
	return client
}
