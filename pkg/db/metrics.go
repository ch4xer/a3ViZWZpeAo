package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

func init() {
	pool := dbPool()
	pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pod_metrics (id SERIAL PRIMARY KEY,pod TEXT NOT NULL,namespace TEXT NOT NULL,cpu_usage TEXT NOT NULL,memory_usage TEXT NOT NULL,timestamp TIMESTAMP NOT NULL)")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_pod_metrics_pod ON pod_metrics(pod)")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_pod_metrics_timestamp ON pod_metrics(timestamp)")
}

func InsertMetrics(pod, namespace, cpu, memory string) error {
	pool := dbPool()
	query := `INSERT INTO pod_metrics (pod, namespace, cpu_usage, memory_usage, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err := pool.Exec(context.Background(), query, pod, namespace, cpu, memory, time.Now())
	if err != nil {
		return fmt.Errorf("InsertMetric failed: %w", err)
	}
	return nil
}

// GetMetrics 查询指定 pod 和 namespace 的最近 limit 条指标记录
func GetMetrics(pod, namespace string, limit int) string {
	pool := dbPool()
	query := `SELECT pod, namespace, cpu_usage, memory_usage, timestamp FROM pod_metrics WHERE pod = $1 AND namespace = $2 ORDER BY timestamp DESC LIMIT $3`

	rows, err := pool.Query(context.Background(), query, pod, namespace, limit)
	if err != nil {
		log.Fatalf("GetMetrics query failed: %v\n", err)
	}
	defer rows.Close()

	var result string
	type PodMetrics struct {
		Pod         string    `json:"pod"`
		Namespace   string    `json:"namespace"`
		CPUUsage    string    `json:"cpu_usage"`
		MemoryUsage string    `json:"memory_usage"`
		Timestamp   time.Time `json:"timestamp"`
	}

	for rows.Next() {
		var metric PodMetrics
		err := rows.Scan(&metric.Pod, &metric.Namespace, &metric.CPUUsage, &metric.MemoryUsage, &metric.Timestamp)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result += fmt.Sprintf("Timestamp: %s CPUUsage: %s MemoryUsage: %s\n", metric.Timestamp, metric.CPUUsage, metric.MemoryUsage)
	}

	return result
}
