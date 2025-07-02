package db

import (
	"context"
	"fmt"
	"slices"
)

func init() {
	pool := dbPool()
	pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS file (pod TEXT NOT NULL,namespace TEXT NOT NULL,files TEXT[],UNIQUE(pod, namespace))")
	pool.Exec(context.Background(), "CREATE INDEX IF NOT EXISTS idx_file_pod ON file(pod)")
}

func UpdateFiles(pod, namespace, file string) error {
	pool := dbPool()
	files, err := GetFiles(pod, namespace)
	if err != nil {
		return err
	}
	if slices.Contains(files, file) {
		return nil
	}
	files = append(files, file)
	query := `INSERT INTO file (pod, namespace, files) VALUES ($1, $2, $3) ON CONFLICT (pod, namespace) DO UPDATE SET files = $3`
	_, err = pool.Exec(context.Background(), query, pod, namespace, files)
	if err != nil {
		return fmt.Errorf("UpdateFiles failed: %w", err)
	}
	return nil
}

func GetFiles(pod, namespace string) ([]string, error) {
	pool := dbPool()
	query := `SELECT files FROM file WHERE pod = $1 AND namespace = $2`
	row := pool.QueryRow(context.Background(), query, pod, namespace)

	var files []string
	err := row.Scan(&files)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return []string{}, nil // No files found for this pod and namespace
		}
		return []string{}, err // Return the error if it's not "no rows"
	}
	return files, nil // Return the list of files
}
