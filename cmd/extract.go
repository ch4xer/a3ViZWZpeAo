package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 声明命令行参数的变量
var (
	namespace   string
	kubeconfig  string
	outputDir   string
	cleanOutput bool
)

var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract native resources from Kubernetes",
	Run:   extract,
}

// extractCmd 处理提取资源的命令
func extract(cmd *cobra.Command, args []string) {
	// 创建输出目录
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// 如果用户选择清空输出目录
	if cleanOutput {
		// 清空输出目录中的所有文件
		err = cleanDirectory(outputDir)
		if err != nil {
			fmt.Printf("Error cleaning output directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Output directory cleared: %s\n", outputDir)
	}

	// 处理kubeconfig路径
	var kubeconfigPath string

	if kubeconfig == "" {
		// 如果用户没有提供kubeconfig路径，使用默认路径
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			os.Exit(1)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	} else if strings.HasPrefix(kubeconfig, "~/") {
		// 处理以~开头的路径
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting user home directory: %v\n", err)
			os.Exit(1)
		}
		kubeconfigPath = filepath.Join(homeDir, (kubeconfig)[2:])
	} else {
		// 使用用户提供的路径
		kubeconfigPath = kubeconfig
	}

	// 检查kubeconfig文件是否存在
	_, err = os.Stat(kubeconfigPath)
	if err != nil {
		fmt.Printf("Error accessing kubeconfig file %s: %v\n", kubeconfigPath, err)
		fmt.Println("Please check if the file exists and you have permission to read it.")
		os.Exit(1)
	}

	// 加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating dynamic client: %v\n", err)
		os.Exit(1)
	}

	// 创建发现客户端
	discoveryClient := clientset.Discovery()

	// 获取服务器支持的API资源类型
	fmt.Printf("Discovering API resource types in namespace %s...\n", namespace)
	resourceTypes, err := getNativeResourceTypes(discoveryClient)
	if err != nil {
		fmt.Printf("Error discovering API resource types: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Exporting resources:")
	exportedCount := 0

	for _, resourceType := range resourceTypes {
		// 只处理命名空间范围内的资源（跳过集群级别的资源）
		if resourceType.Namespaced {

			err := exportResource(dynamicClient, resourceType, namespace, outputDir)
			if err != nil {
				fmt.Printf("  - Error exporting %s: %v\n", resourceType.Kind, err)
			} else {
				exportedCount += 1
			}
		}
	}

	summaryMsg := fmt.Sprintf("\nSuccessfully exported %d resource types to %s", exportedCount, outputDir)
	fmt.Println(summaryMsg)
}

// isK8sNativeResource 判断资源是否为Kubernetes原生资源
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

// getNativeResourceTypes 获取集群支持的所有原生API资源类型（而非资源实例）
func getNativeResourceTypes(discoveryClient discovery.DiscoveryInterface) ([]metav1.APIResource, error) {
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

// exportResource 导出指定命名空间中的单个资源类型的所有资源实例
func exportResource(client dynamic.Interface, resourceType metav1.APIResource, namespace, outputDir string) error {
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

// cleanDirectory 清空目录中的所有文件，但会进行安全检查
func cleanDirectory(dirPath string) error {
	// 安全检查：防止清空关键目录
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("无法获取绝对路径: %v", err)
	}

	// 检查是否是危险目录
	dangerousPaths := []string{
		"/", "/bin", "/boot", "/etc", "/home", "/lib", "/lib64",
		"/opt", "/root", "/sbin", "/usr", "/var",
	}

	// 获取用户主目录并加入危险目录列表
	homeDir, err := os.UserHomeDir()
	if err == nil {
		dangerousPaths = append(dangerousPaths, homeDir)
	}

	if slices.Contains(dangerousPaths, absPath) {
		return fmt.Errorf("安全限制: 不允许清空系统关键目录 '%s'", absPath)
	}

	// 读取目录中的所有项目
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// 遍历并删除每个项目
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		// 如果是目录，递归删除
		if entry.IsDir() {
			if err := os.RemoveAll(fullPath); err != nil {
				return err
			}
		} else {
			// 如果是文件，直接删除
			if err := os.Remove(fullPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
