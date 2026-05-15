package cleanup

import (
	"context"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (cm *CleanupManager) checkKubernetesDeployments(ctx context.Context, env string, threshold time.Duration) bool {
	client, err := getKubernetesClient()
	if err != nil {
		log.Printf("Failed to create Kubernetes client: %v", err)
		return false
	}

	selector := labels.Set{"environment": env}.AsSelector().String()
	deployments, err := client.AppsV1().Deployments(metav1.NamespaceAll).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		log.Printf("Failed to list deployments for environment %s: %v", env, err)
		return false
	}

	cutoff := time.Now().Add(-threshold)
	for _, deployment := range deployments.Items {
		if deploymentRecentlyUpdated(&deployment, cutoff) {
			log.Printf("Deployment %s/%s updated within threshold", deployment.Namespace, deployment.Name)
			return false
		}
	}

	return true
}

func deploymentRecentlyUpdated(dep *appsv1.Deployment, cutoff time.Time) bool {
	if dep.CreationTimestamp.After(cutoff) {
		return true
	}

	for _, condition := range dep.Status.Conditions {
		if condition.LastUpdateTime.After(cutoff) || condition.LastTransitionTime.After(cutoff) {
			return true
		}
	}

	return false
}
