package db

import (
	"context"
	"fmt"
	"kubefix-cli/pkg/model"
	"log"
	"time"
)

func init() {
	pool := dbPool()
	pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pod_metrics (id SERIAL PRIMARY KEY,pod_name TEXT NOT NULL,namespace TEXT NOT NULL,cpu_usage TEXT NOT NULL,memory_usage TEXT NOT NULL,timestamp TIMESTAMP NOT NULL)")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_pod_metrics_pod_name ON pod_metrics(pod_name)")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_pod_metrics_timestamp ON pod_metrics(timestamp)")
}

func InsertMetrics(metrics model.PodMetrics) error {
	pool := dbPool()
	query := `INSERT INTO pod_metrics (pod_name, namespace, cpu_usage, memory_usage, timestamp) VALUES ($1, $2, $3, $4, $5)`
	_, err := pool.Exec(context.Background(), query, metrics.PodName, metrics.Namespace, metrics.CPUUsage, metrics.MemoryUsage, time.Now())
	if err != nil {
		return fmt.Errorf("InsertMetric failed: %w", err)
	}
	return nil
}

// GetMetrics 查询指定 pod 和 namespace 的最近 limit 条指标记录
func GetMetrics(pod, namespace string, limit int) string {
	pool := dbPool()
	query := `SELECT pod_name, namespace, cpu_usage, memory_usage, timestamp FROM pod_metrics WHERE pod_name = $1 AND namespace = $2 ORDER BY timestamp DESC LIMIT $3`

	rows, err := pool.Query(context.Background(), query, pod, namespace, limit)
	if err != nil {
		log.Fatalf("GetMetrics query failed: %v\n", err)
	}
	defer rows.Close()

	var result string
	for rows.Next() {
		var metric model.PodMetrics
		err := rows.Scan(&metric.PodName, &metric.Namespace, &metric.CPUUsage, &metric.MemoryUsage, &metric.Timestamp)
		if err != nil {
			log.Fatalf("Error scanning row: %v\n", err)
		}
		result += fmt.Sprintf("Timestamp: %s CPUUsage: %s MemoryUsage: %s\n", metric.Timestamp, metric.CPUUsage, metric.MemoryUsage)
	}

	return result
}
