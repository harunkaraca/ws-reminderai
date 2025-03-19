package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"reminderai/model"
)

// LogRepository handles database operations for logs
type LogRepository struct {
	pool *pgxpool.Pool
}

func NewLogRepository(pool *pgxpool.Pool) *LogRepository {
	return &LogRepository{pool}
}

func (r *LogRepository) Create(message, level string) error {
	ctx := context.Background()
	_, err := r.pool.Exec(ctx,
		"INSERT INTO logs (message, level) VALUES ($1, $2)",
		message, level)
	return err
}

func (r *LogRepository) GetAll() ([]model.Log, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, "SELECT id, message, level, created_at FROM logs ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.Log
	for rows.Next() {
		var log model.Log
		if err := rows.Scan(&log.ID, &log.Message, &log.Level, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
