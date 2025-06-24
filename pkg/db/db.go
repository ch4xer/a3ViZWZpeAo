package db

import (
	"context"
	"kubefix-cli/conf"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func dbPool() *pgxpool.Pool {
	once.Do(func() {
		var err error
		ctx := context.Background()
		config, err := pgxpool.ParseConfig(conf.Database)
		if err != nil {
			log.Fatalf("Unable to parse db config: %v\n", err)
		}
		pool, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
	})
	return pool
}

// ClosePool 关闭数据库连接池
func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

// WithTransaction 提供了一个事务上下文的辅助函数
func WithTransaction(ctx context.Context, txFunc func(pgx.Tx) error) error {
	pool := dbPool()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			// 发生 panic，回滚事务
			tx.Rollback(ctx)
			panic(r) // 重新抛出 panic
		} else if err != nil {
			// 发生错误，回滚事务
			tx.Rollback(ctx)
		} else {
			// 一切正常，提交事务
			err = tx.Commit(ctx)
		}
	}()

	err = txFunc(tx)
	return err
}
