package client

import (
	"fmt"
	"kubefix-cli/conf"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/discovery"
)

func Client() (*kubernetes.Clientset, error) {
	var kubeconfigPath string
	if strings.HasPrefix(conf.Kubeconfig, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Error getting user home directory: %v\n", err)
		}
		kubeconfigPath = filepath.Join(homeDir, (conf.Kubeconfig)[2:])
	} else {
		kubeconfigPath = conf.Kubeconfig
	}

	// 检查kubeconfig文件是否存在
	_, err := os.Stat(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("Error accessing kubeconfig file %s: %v\n", kubeconfigPath, err)
	}

	// 加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("Error building kubeconfig: %v\n", err)
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernetes client: %v\n", err)
	}

	return clientset, nil
}

func DynamicClient() (dynamic.Interface, error) {
	var kubeconfigPath string
	if strings.HasPrefix(conf.Kubeconfig, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		kubeconfigPath = filepath.Join(homeDir, (conf.Kubeconfig)[2:])
	} else {
		kubeconfigPath = conf.Kubeconfig
	}

	// 检查kubeconfig文件是否存在
	_, err := os.Stat(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing kubeconfig file %s: %v", kubeconfigPath, err)
	}

	// 加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating dynamic client: %v", err)
	}

	return dynamicClient, nil
}

func DiscoveryClient() (*discovery.DiscoveryClient, error) {
	var kubeconfigPath string
	if strings.HasPrefix(conf.Kubeconfig, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting user home directory: %v", err)
		}
		kubeconfigPath = filepath.Join(homeDir, (conf.Kubeconfig)[2:])
	} else {
		kubeconfigPath = conf.Kubeconfig
	}

	// 检查kubeconfig文件是否存在
	_, err := os.Stat(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error accessing kubeconfig file %s: %v", kubeconfigPath, err)
	}

	// 加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig: %v", err)
	}
	// 创建DiscoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating DiscoveryClient: %v", err)
	}

	return discoveryClient, nil
}