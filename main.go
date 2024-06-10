package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Unable to get kubeconfig: %v", err)
	}

	// Create a new Manager to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		log.Fatalf("Unable to set up overall controller manager: %v", err)
	}

	// Create a context that is cancelled when SIGINT or SIGTERM is received.
	ctx := signals.SetupSignalHandler()

	// Start the Manager in a separate goroutine
	go func() {
		if err := mgr.Start(ctx); err != nil {
			log.Fatalf("Failed to start manager: %v", err)
		}
	}()

	// Wait for the cache to sync
	if !mgr.GetCache().WaitForCacheSync(ctx) {
		log.Fatalf("Cache sync failed")
	}

	log.Println("Cache has been synced successfully")

	// At this point, the cache is synced, and you can read objects from the cache
	log.Println("Listing all pods in all namespaces")

	// Example: Get a client to interact with the cluster
	k8sClient := mgr.GetClient()
	

	// List all pods in all namespaces
	podList := &corev1.PodList{}
	err = k8sClient.List(context.Background(), podList, &client.ListOptions{
		Namespace: metav1.NamespaceAll,
	})
	if err != nil {
		log.Fatalf("Failed to list pods: %v", err)
	}

	// Print the names of the pods
	for _, pod := range podList.Items {
		fmt.Printf("Pod Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
	}
}
