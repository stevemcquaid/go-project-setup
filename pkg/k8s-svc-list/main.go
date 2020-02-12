package k8ssvclist

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func Run() {
	c, err := client.New(k8sconfig.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	services := &corev1.ServiceList{}
	if err := c.List(context.Background(), services); err != nil {
		log.Fatal(err)
	}
	for _, svc := range services.Items {
		log.Printf("Service Name: %s. Service IP: %s", svc.Name, svc.Spec.ClusterIP)
	}
}
