package client

import (
	"context"
	"fmt"
	"kubefix-cli/conf"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


func GetNamespaces() ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", conf.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("加载 kubeconfig 失败: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	} 
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取命名空间列表失败: %v", err)
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