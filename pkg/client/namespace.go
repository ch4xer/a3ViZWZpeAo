package client

import (
	"context"
	"fmt"
	"kubefix-cli/conf"
	"kubefix-cli/pkg/db"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Namespaces() ([]string, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	result := []string{}
	for _, ns := range namespaces.Items {
		if slices.Contains(conf.IgnoreNamespaces, ns.Name) {
			continue
		}
		result = append(result, ns.Name)
	}
	return result, nil
}

func CollectNamespace() {
	fmt.Println("Collecting namespaces...")
	namespaces, err := Namespaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, ns := range namespaces {
		err = db.InsertNamespace(ns)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Collecting namespaces finished.")
}
