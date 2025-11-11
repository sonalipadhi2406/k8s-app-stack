package main

import (
	"log"
	"os"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Create in-cluster configuration (official method)
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Unable to create in-cluster config: %v", err)
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Unable to create Kubernetes client: %v", err)
	}

	// Read namespace from environment variable
	namespace := os.Getenv("WATCH_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}
	log.Printf("Starting Pod Watcher for namespace: %s", namespace)

	// Watch Pods
	watchPods(clientset, namespace)
}

func watchPods(clientset *kubernetes.Clientset, namespace string) {
	watchlist := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		namespace,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		30*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Printf("Pod Created: %s (%s)", pod.Name, pod.Status.Phase)
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				log.Printf("Pod Deleted: %s", pod.Name)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
				newPod := newObj.(*v1.Pod)
				if oldPod.Status.Phase != newPod.Status.Phase {
					log.Printf("Pod Updated: %s (%s → %s)",
						newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
				}
			},
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	controller.Run(stop)
}
