package db

import (
	"context"
	"fmt"
)

func init() {
	pool := dbPool()
	pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS namespace (id SERIAL PRIMARY KEY,namespace TEXT NOT NULL,UNIQUE(namespace))")
}

func InsertNamespace(namespace string) error {
	pool := dbPool()
	query := `INSERT INTO namespace (namespace) VALUES ($1) ON CONFLICT (namespace) DO NOTHING`
	_, err := pool.Exec(context.Background(), query, namespace)
	if err != nil {
		return fmt.Errorf("InsertNamespace failed: %w", err)
	}
	return nil
}
