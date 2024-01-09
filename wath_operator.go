package main

import (
	"context"
	"fmt"
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

func startWatchOperator(ctx context.Context, config *rest.Config, namespace, name string) {
	dynamicClient, err := dynamic.NewForConfig(config)
	chErr(err)
	go watchOperatorByInformer(ctx, dynamicClient, namespace, name)
}

func watchOperatorByInformer(ctx context.Context, client *dynamic.DynamicClient, namespace, name string) {
	fac := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, namespace, nil)
	informer := fac.ForResource(schema.GroupVersionResource{
		Group:    "operator.openshift.io",
		Version:  "v1",
		Resource: "ingresscontrollers",
	}).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			typedObj := obj.(*unstructured.Unstructured)
			bytes, _ := typedObj.MarshalJSON()

			crdObj := operatorv1.IngressController{}
			json.Unmarshal(bytes, &crdObj)
			if crdObj.Name == name {
				fmt.Println(time.Now().String(), " ingresscontrollers added: ", crdObj.Name, " CreationTimestamp: ", crdObj.CreationTimestamp)
			}
		},
		DeleteFunc: func(obj interface{}) {
			typedObj := obj.(*unstructured.Unstructured)
			bytes, _ := typedObj.MarshalJSON()

			crdObj := operatorv1.IngressController{}
			json.Unmarshal(bytes, &crdObj)
			if crdObj.Name == name {
				fmt.Println(time.Now().String(), " ingresscontrollers deleted: ", crdObj.Name, " DeletionTimestamp: ", crdObj.DeletionTimestamp)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			typedObj := new.(*unstructured.Unstructured)
			bytes, _ := typedObj.MarshalJSON()

			crdObj := operatorv1.IngressController{}
			json.Unmarshal(bytes, &crdObj)
			if crdObj.Name == name {
				fmt.Println(time.Now().String(), " ingresscontrollers updated new: ", crdObj.Name, " DeletionTimestamp: ", crdObj.DeletionTimestamp)
			}
		},
	})
	fmt.Println("running informer of ingresscontrollers")
	informer.Run(ctx.Done())
	fmt.Println("ctx Done, exit watching", ctx.Err())
}

func watchOperator(ctx context.Context, client *dynamic.DynamicClient, namespace, name string) {
	w, err := client.Resource(schema.GroupVersionResource{
		Group:    "operator.openshift.io",
		Version:  "v1",
		Resource: "ingresscontrollers",
	}).Namespace(namespace).Watch(context.Background(), metav1.ListOptions{})
	chErr(err)

	w = watch.Filter(w, func(in watch.Event) (watch.Event, bool) {
		return in, filterToSystemNamespaces(in.Object, name)
	})
	for {
		select {
		case event := <-w.ResultChan():
			switch event.Type {
			case watch.Deleted, watch.Added, watch.Modified:
				typedObj := event.Object.(*unstructured.Unstructured)
				bytes, _ := typedObj.MarshalJSON()
				obj := operatorv1.IngressController{}
				json.Unmarshal(bytes, &obj)
				logPrefix := time.Now().String() + " | " + "event: " + string(event.Type) + " | " + "name: " + name + " | "
				fmt.Println(logPrefix, "create: ", obj.CreationTimestamp, " | delete: ", obj.DeletionTimestamp)
			}
		case <-ctx.Done():
			fmt.Println("ctx exited watching ", ctx.Err())
			return
		}
	}
}
