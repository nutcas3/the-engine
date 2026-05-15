package cleanup

import (
	"context"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func (cm *CleanupManager) checkFluxDeployments(ctx context.Context, env string, threshold time.Duration) bool {
	client, err := getDynamicClient()
	if err != nil {
		log.Printf("Failed to create dynamic Kubernetes client: %v", err)
		return false
	}

	selector := labels.Set{"environment": env}.AsSelector().String()
	cutoff := time.Now().Add(-threshold)

	kustomizationGVR := schema.GroupVersionResource{Group: "kustomize.toolkit.fluxcd.io", Version: "v1", Resource: "kustomizations"}
	if !cm.fluxResourcesIdle(ctx, client, kustomizationGVR, selector, cutoff) {
		return false
	}

	helmGVR := schema.GroupVersionResource{Group: "helm.toolkit.fluxcd.io", Version: "v2", Resource: "helmreleases"}
	if !cm.fluxResourcesIdle(ctx, client, helmGVR, selector, cutoff) {
		return false
	}

	return true
}

func (cm *CleanupManager) fluxResourcesIdle(ctx context.Context, client dynamic.Interface, gvr schema.GroupVersionResource, selector string, cutoff time.Time) bool {
	resources, err := client.Resource(gvr).List(ctx, v1.ListOptions{LabelSelector: selector})
	if err != nil {
		log.Printf("Failed to list Flux resource %s: %v", gvr.Resource, err)
		return false
	}

	for _, item := range resources.Items {
		if fluxResourceUpdatedRecently(&item, cutoff) {
			log.Printf("Flux resource %s/%s updated within threshold", item.GetNamespace(), item.GetName())
			return false
		}
	}

	return true
}

func fluxResourceUpdatedRecently(obj *unstructured.Unstructured, cutoff time.Time) bool {
	status, ok := obj.Object["status"].(map[string]interface{})
	if !ok {
		return false
	}

	conditions, ok := status["conditions"].([]interface{})
	if !ok {
		return false
	}

	for _, cond := range conditions {
		condition, ok := cond.(map[string]interface{})
		if !ok {
			continue
		}
		timestamp, _ := condition["lastTransitionTime"].(string)
		if timestamp == "" {
			continue
		}
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil && t.After(cutoff) {
			return true
		}
	}

	return false
}
