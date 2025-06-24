package metrics

import (
	"context"
	"fmt"
	"kubefix-cli/pkg/model"
	"time"
	"kubefix-cli/conf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

// GetPodCPUAndMemoryUsage 获取指定命名空间下所有 Pod 的 CPU 和内存使用量（单位：mCPU 和 MiB）
func GetPodCPUAndMemoryUsage(namespace string) ([]model.PodMetrics, error) {
	config, err := clientcmd.BuildConfigFromFlags("", conf.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("加载 kubeconfig 失败: %v", err)
	}
	metricsClient, err := metricsclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 metrics client 失败: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取 Pod 指标失败: %v", err)
	}
	result := []model.PodMetrics{}
	for _, podMetric := range podMetricsList.Items {
		cpuTotal := int64(0)
		memTotal := int64(0)
		for _, c := range podMetric.Containers {
			cpu := c.Usage.Cpu().MilliValue()               // mCPU
			mem := c.Usage.Memory().Value() / (1024 * 1024) // MiB
			cpuTotal += cpu
			memTotal += mem
		}

		var podMatrics model.PodMetrics

		podMatrics.PodName = podMetric.Name
		podMatrics.Namespace = podMetric.Namespace
		podMatrics.CPUUsage = fmt.Sprintf("%dm", cpuTotal)
		podMatrics.MemoryUsage = fmt.Sprintf("%dMiB", memTotal)

		result = append(result, podMatrics)
	}
	return result, nil
}
