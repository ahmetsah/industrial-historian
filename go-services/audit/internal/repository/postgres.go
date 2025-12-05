package repository

import (
	"context"
	"fmt"

	"github.com/ahmetsah/industrial-historian/go-services/audit/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AppendLog(ctx context.Context, log *core.LogEntry, hasher core.Hasher) error
	IterateLogs(ctx context.Context, callback func(*core.LogEntry) error) error
	Close()
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(ctx context.Context, dbUrl string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	migrationSQL := `
	CREATE TABLE IF NOT EXISTS audit_logs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		actor VARCHAR(255) NOT NULL,
		action VARCHAR(255) NOT NULL,
		details JSONB,
		prev_hash VARCHAR(64) NOT NULL,
		curr_hash VARCHAR(64) NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_logs(actor);
	`

	_, err = pool.Exec(ctx, migrationSQL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}

func (r *PostgresRepository) IterateLogs(ctx context.Context, callback func(*core.LogEntry) error) error {
	// TODO: Add pagination support. For now, limit to 1000 to prevent DoS.
	query := `SELECT id, timestamp, actor, action, details, prev_hash, curr_hash FROM audit_logs ORDER BY timestamp ASC LIMIT 1000`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var log core.LogEntry
		err := rows.Scan(&log.ID, &log.Timestamp, &log.Actor, &log.Action, &log.Details, &log.PrevHash, &log.CurrHash)
		if err != nil {
			return err
		}
		if err := callback(&log); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (r *PostgresRepository) AppendLog(ctx context.Context, log *core.LogEntry, hasher core.Hasher) error {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := r.appendLogTx(ctx, log, hasher)
		if err == nil {
			return nil
		}

		// Check for serialization failure (SQLState 40001)
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "40001" {
			// Exponential backoff could be added here, but simple retry is often enough for optimistic locking
			continue
		}

		return err
	}
	return fmt.Errorf("failed to append log after %d retries", maxRetries)
}

func (r *PostgresRepository) appendLogTx(ctx context.Context, log *core.LogEntry, hasher core.Hasher) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get last log
	var prevHash string
	query := `SELECT curr_hash FROM audit_logs ORDER BY timestamp DESC LIMIT 1`
	err = tx.QueryRow(ctx, query).Scan(&prevHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			prevHash = "0000000000000000000000000000000000000000000000000000000000000000" // Genesis hash
		} else {
			return fmt.Errorf("failed to get last log: %w", err)
		}
	}

	// Calculate new hash (using the timestamp already set in log.Timestamp)
	log.PrevHash = prevHash
	log.CurrHash = hasher.Hash(prevHash, log)

	// Debug: log what we're inserting
	// fmt.Printf("DEBUG INSERT: timestamp=%v, actor=%s, action=%s, details=%s, hash=%s\n",
	// 	log.Timestamp, log.Actor, log.Action, string(log.Details), log.CurrHash)

	// Insert with explicit timestamp to match hash calculation
	insertQuery := `INSERT INTO audit_logs (timestamp, actor, action, details, prev_hash, curr_hash) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = tx.QueryRow(ctx, insertQuery, log.Timestamp, log.Actor, log.Action, log.Details, log.PrevHash, log.CurrHash).Scan(&log.ID)
	if err != nil {
		return err // Let the caller handle wrapping/checking
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
