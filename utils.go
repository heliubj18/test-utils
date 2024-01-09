package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func chErr(err error) {
	if err != nil {
		panic(err)
	}
}

func filterToSystemNamespaces(obj runtime.Object, name string) bool {
	m, ok := obj.(metav1.Object)
	if !ok {
		return true
	}
	return m.GetName() == name
}
