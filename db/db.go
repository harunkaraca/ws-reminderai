package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates and returns a new connection pool
func NewPool() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://u5qu312i1icosp:pf11eb121e985099fbf38ea9681fdd66ff6d28484dfa5a57288651a94f91b4d3b@c3l5o0rb2a6o4l.cluster-czz5s0kz4scl.eu-west-1.rds.amazonaws.com:5432/d6nuqnl6o2lgp8"
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// You can configure pool settings here if needed
	// config.MaxConns = 10

	return pgxpool.NewWithConfig(context.Background(), config)
}

// PrintPoolStats prints connection pool statistics to console
func PrintPoolStats(pool *pgxpool.Pool) {
	stats := pool.Stat()
	fmt.Printf("Pool Stats - Total: %d, Idle: %d, In Use: %d, Max: %d\n",
		stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns(), stats.MaxConns())
}

// MonitorPoolStats periodically prints pool statistics to console
func MonitorPoolStats(pool *pgxpool.Pool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			PrintPoolStats(pool)
		}
	}()
}

// InitSchema creates the necessary tables if they don't exist
func InitSchema(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create books table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create logs table
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			message TEXT NOT NULL,
			level VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
