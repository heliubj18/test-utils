package main

import (
	"context"
	"fmt"
	operatorv1 "github.com/openshift/api/operator/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func startWatchOperator(ctx context.Context, config *rest.Config, namespace string) {
	dynamicClient, err := dynamic.NewForConfig(config)
	chErr(err)
	go watchOperator(ctx, dynamicClient, namespace)
}

func watchOperator(ctx context.Context, client *dynamic.DynamicClient, namespace string) {
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
			fmt.Println("ingresscontrollers added: ", crdObj.Name, " CreationTimestamp: ", crdObj.CreationTimestamp)
		},
		DeleteFunc: func(obj interface{}) {
			typedObj := obj.(*unstructured.Unstructured)
			bytes, _ := typedObj.MarshalJSON()

			crdObj := operatorv1.IngressController{}
			json.Unmarshal(bytes, &crdObj)
			fmt.Println("ingresscontrollers deleted: ", crdObj.Name, " DeletionTimestamp: ", crdObj.DeletionTimestamp)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			typedObj := new.(*unstructured.Unstructured)
			bytes, _ := typedObj.MarshalJSON()

			crdObj := operatorv1.IngressController{}
			json.Unmarshal(bytes, &crdObj)
			fmt.Println("ingresscontrollers updated new: ", crdObj.Name, " DeletionTimestamp: ", crdObj.DeletionTimestamp)
		},
	})
	fmt.Println("running informer of ingresscontrollers")
	informer.Run(ctx.Done())
	fmt.Println("ctx Done, exit watching", ctx.Err())
}
