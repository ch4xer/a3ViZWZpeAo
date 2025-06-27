package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

func isK8sNativeResource(groupVersion string) bool {
	// 无组资源(核心资源) - 例如 Pod, Service, ConfigMap 等
	if groupVersion == "v1" {
		return true
	}

	// 只允许以下Kubernetes原生API组
	allowedGroups := map[string]bool{
		"apps":                      true, // Deployment, StatefulSet, DaemonSet
		"batch":                     true, // Job, CronJob
		"networking.k8s.io":         true, // Ingress, NetworkPolicy
		"storage.k8s.io":            true, // StorageClass
		"rbac.authorization.k8s.io": true, // Role, ClusterRole
		"policy":                    true, // PodDisruptionBudget
	}

	gv, err := schema.ParseGroupVersion(groupVersion)
	if err != nil {
		return false
	}

	return allowedGroups[gv.Group]
}


func NativeResourceTypes(discoveryClient discovery.DiscoveryInterface) ([]metav1.APIResource, error) {
	// 获取服务器上的所有API组和资源类型
	_, resourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		return nil, fmt.Errorf("failed to get server resource types: %v", err)
	}

	var resourceTypes []metav1.APIResource

	// 遍历所有API资源类型
	for _, list := range resourceList {
		// 解析组/版本
		groupVersion := list.GroupVersion

		for _, r := range list.APIResources {
			if strings.Contains(r.Name, "/") {
				continue
			}

			// 确保资源可列表和获取
			if !containsVerb(r.Verbs, "list") || !containsVerb(r.Verbs, "get") {
				continue
			}

			if isK8sNativeResource(groupVersion) {
				// 解析组和版本
				gv, err := schema.ParseGroupVersion(groupVersion)
				if err != nil {
					continue
				}

				// 存储组和版本信息
				r.Group = gv.Group
				r.Version = gv.Version
				resourceTypes = append(resourceTypes, r)
			}
		}
	}

	return resourceTypes, nil
}

func ExportResource(client dynamic.Interface, resourceType metav1.APIResource, namespace, outputDir string) error {
	// 使用已经解析好的组和版本创建GVR
	fmt.Printf("  - Processing resource type: %s (Group: '%s', Version: '%s')\n",
		resourceType.Kind, resourceType.Group, resourceType.Version)

	// 创建资源类型的GVR
	gvr := schema.GroupVersionResource{
		Group:    resourceType.Group,
		Version:  resourceType.Version,
		Resource: resourceType.Name,
	}

	// 使用动态客户端获取该类型的所有资源实例
	list, err := client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(list.Items) == 0 {
		return nil
	}

	fmt.Printf("  - Found %d %s resources\n", len(list.Items), resourceType.Kind)

	// 遍历每个资源并导出为YAML
	for _, item := range list.Items {
		cleanObject(&item)
		// 生成文件名
		name := item.GetName()
		kind := strings.ToLower(resourceType.Kind)
		filename := fmt.Sprintf("%s-%s.yaml", kind, name)
		filePath := filepath.Join(outputDir, filename)

		// 转换为YAML
		yamlBytes, err := yaml.Marshal(item.Object)
		if err != nil {
			fmt.Printf("    - Error marshaling %s/%s: %v\n", resourceType.Kind, name, err)
			continue
		}

		// 写入文件
		err = os.WriteFile(filePath, yamlBytes, 0644)
		if err != nil {
			fmt.Printf("    - Error writing %s/%s: %v\n", resourceType.Kind, name, err)
			continue
		}
		fmt.Printf("    - Exported: %s\n", filename)
	}

	return nil
}

// containsVerb 检查动词列表中是否包含特定动词
func containsVerb(verbs []string, verb string) bool {
	for _, v := range verbs {
		if v == verb || v == "*" {
			return true
		}
	}
	return false
}

// cleanObject 清理对象以适合导出
func cleanObject(obj *unstructured.Unstructured) {
	// 清理常见的服务器端字段
	metadata := obj.Object["metadata"].(map[string]any)

	// 移除字段
	delete(metadata, "creationTimestamp")
	delete(metadata, "resourceVersion")
	delete(metadata, "selfLink")
	delete(metadata, "uid")
	delete(metadata, "generation")
	delete(metadata, "managedFields")
	delete(metadata, "ownerReferences")

	// 如果存在状态字段，移除它
	delete(obj.Object, "status")
}
