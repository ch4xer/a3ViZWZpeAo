package db

import (
	"context"
	"fmt"
	"slices"
)

func init() {
	pool := dbPool()
	pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS pod_capability (id SERIAL PRIMARY KEY,pod TEXT NOT NULL,namespace TEXT NOT NULL,caps TEXT[])")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_pod_capability_pod ON pod_capability(pod)")
}

func UpdateCaps(pod, namespace, cap string) error {
	pool := dbPool()
	caps, err := GetCaps(pod, namespace)
	if err != nil {
		return err
	}
	if slices.Contains(caps, cap) {
		return nil // cap already exists, no need to update
	}
	caps = append(caps, cap)
	query := `INSERT INTO pod_capability (pod, namespace, caps) VALUES ($1, $2, $3) ON CONFLICT (pod, namespace) DO UPDATE SET caps = $3`
	_, err = pool.Exec(context.Background(), query, pod, namespace, caps)
	if err != nil {
		return fmt.Errorf("UpdateCaps failed: %w", err)
	}
	return nil
}

func GetCaps(pod, namespace string) ([]string, error) {
	pool := dbPool()
	query := `SELECT caps FROM pod_capability WHERE pod = $1 AND namespace = $2`
	row := pool.QueryRow(context.Background(), query, pod, namespace)

	var caps []string
	err := row.Scan(&caps)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return []string{}, nil // No syscall found for this pod and namespace
		}
		return []string{}, fmt.Errorf("GetCaps query failed: %w", err)
	}
	return caps, nil
}
