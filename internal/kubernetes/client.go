package kubernetes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client holds Kubernetes clients
type Client struct {
	K8sClient     *kubernetes.Clientset
	DynamicClient dynamic.Interface
	Transport     http.RoundTripper
}

// NewClient creates a new Kubernetes client
func NewClient() (*Client, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	config.Transport = createTransport()

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
		Transport:     config.Transport,
	}, nil
}

// createTransport creates a transport with connection pooling
func createTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
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
