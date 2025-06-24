package model

import "time"

// PodMetric 结构体用于存储 pod_metrics 表的数据
type PodMetrics struct {
	PodName     string    `json:"pod_name"`
	Namespace   string    `json:"namespace"`
	CPUUsage    string    `json:"cpu_usage"`
	MemoryUsage string    `json:"memory_usage"`
	Timestamp   time.Time `json:"timestamp"`
}
