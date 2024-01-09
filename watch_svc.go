package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func startWatchSVC(ctx context.Context, config *rest.Config, namespace string) {
	clientSet, err := kubernetes.NewForConfig(config)
	chErr(err)
	go watchSVC(ctx, clientSet, namespace)
}

func watchSVC(ctx context.Context, clientSet *kubernetes.Clientset, namespace string) {
	//selector, err := labels.NewRequirement("name", selection.Equals, []string{"ingress-operator"})
	//chErr(err)
	fac := informers.NewSharedInformerFactoryWithOptions(clientSet, 0, informers.WithNamespace(namespace)) //informers.WithTweakListOptions(func(options *metav1.ListOptions) {
	//options.LabelSelector = selector.String()})

	svcInformer := fac.Core().V1().Services().Informer()

	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svcObj, ok := obj.(*corev1.Service)
			if ok != true {
				return
			}

			fmt.Println("svc added: ", svcObj.Name, " CreationTimestamp: ", svcObj.CreationTimestamp)
		},
		DeleteFunc: func(obj interface{}) {
			svcObj, ok := obj.(*corev1.Service)
			if ok != true {
				return
			}

			fmt.Println("svc deleted: ", svcObj.Name, " DeletionTimestamp: ", svcObj.DeletionTimestamp)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			svcObj, ok := new.(*corev1.Service)
			if ok != true {
				return
			}

			fmt.Println("svc updated new: ", svcObj.Name, " DeletionTimestamp: ", svcObj.DeletionTimestamp)
		},
	})

	fmt.Println("running informer of svc")
	svcInformer.Run(ctx.Done())
	fmt.Println("ctx Done, exit watching", ctx.Err())
}
