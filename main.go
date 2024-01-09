package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	namespace := "openshift-ingress-operator"
	kubeconfigPath := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	chErr(err)
	startWatchOperator(ctx, config, namespace, "default")

	namespace = "openshift-ingress"
	startWatchSVC(ctx, config, namespace, "router-default")

	time.Sleep(time.Second * 600)
	cancel()
	time.Sleep(time.Second * 5)
	fmt.Println("end of main")
}
