package common

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sprm-crawler/repository"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	Pool  *pgxpool.Pool
	Query *repository.Queries
)

func ConnectDatabase() error {
	// Build connection string from environment variables
	host := getEnv("DATABASE_HOST", "localhost")
	port := getEnv("DATABASE_PORT", "5432")
	user := getEnv("DATABASE_USER", "postgres")
	password := getEnv("DATABASE_PASSWORD", "password")
	dbname := getEnv("DATABASE_NAME", "sprm_crawler")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	newPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := newPool.Ping(ctx); err != nil {
		newPool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	Pool = newPool
	return nil
}

func SetQuery(newQuery *repository.Queries) error {
	if newQuery == nil {
		return errors.New("cannot assign nil query")
	}
	Query = newQuery
	return nil
}

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
