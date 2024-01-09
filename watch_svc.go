package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

func startWatchSVC(ctx context.Context, config *rest.Config, namespace, name string) {
	clientSet, err := kubernetes.NewForConfig(config)
	chErr(err)
	go watchSVCByInformer(ctx, clientSet, namespace, name)
	//go watchSVC(ctx, clientSet, namespace, name)
}

func watchSVCByInformer(ctx context.Context, clientSet *kubernetes.Clientset, namespace, name string) {
	//selector, err := labels.NewRequirement("name", selection.Equals, []string{name})
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
			if svcObj.Name == name {
				fmt.Println(time.Now().String(), " svc added: ", svcObj.Name, " CreationTimestamp: ", svcObj.CreationTimestamp)
			}
		},
		DeleteFunc: func(obj interface{}) {
			svcObj, ok := obj.(*corev1.Service)
			if ok != true {
				return
			}
			if svcObj.Name == name {
				fmt.Println(time.Now().String(), " svc deleted: ", svcObj.Name, " DeletionTimestamp: ", svcObj.DeletionTimestamp)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			svcObj, ok := new.(*corev1.Service)
			if ok != true {
				return
			}
			if svcObj.Name == name {
				fmt.Println(time.Now().String(), " svc updated new: ", svcObj.Name, " DeletionTimestamp: ", svcObj.DeletionTimestamp)
			}
		},
	})

	fmt.Println("running informer of svc")
	svcInformer.Run(ctx.Done())
	fmt.Println("ctx Done, exit watching", ctx.Err())
}

func watchSVC(ctx context.Context, clientSet *kubernetes.Clientset, namespace, name string) {
	w, err := clientSet.CoreV1().Services(namespace).Watch(ctx, metav1.ListOptions{})
	chErr(err)
	defer w.Stop()
	w = watch.Filter(w, func(in watch.Event) (watch.Event, bool) {
		return in, filterToSystemNamespaces(in.Object, name)
	})

	fmt.Println("watching ResultChan of the watcher of svc")
	for {
		select {
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Deleted, watch.Added, watch.Modified:
				obj, ok := event.Object.(*corev1.Service)
				if !ok {
					continue
				}
				logPrefix := time.Now().String() + " | " + "event: " + string(event.Type) + " | " + "name: " + name + " | "
				fmt.Println(logPrefix, "create: ", obj.CreationTimestamp, " | delete: ", obj.DeletionTimestamp)
			}
		case <-ctx.Done():
			fmt.Println("ctx exited watching ", ctx.Err())
			return
		}
	}
}
