package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmetsah/industrial-historian/go-services/alarm/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(dbUrl string) (*PostgresRepository, error) {
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return &PostgresRepository{pool: pool}, nil
}

func (r *PostgresRepository) Close() {
	r.pool.Close()
}

func (r *PostgresRepository) CreateDefinition(def *core.AlarmDefinition) error {
	query := `
		INSERT INTO alarm_definitions (tag, threshold, alarm_type, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(context.Background(), query, def.Tag, def.Threshold, def.Type, def.Priority).
		Scan(&def.ID, &def.CreatedAt, &def.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create definition: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetDefinition(id int) (*core.AlarmDefinition, error) {
	query := `
		SELECT id, tag, threshold, alarm_type, priority, created_at, updated_at
		FROM alarm_definitions
		WHERE id = $1
	`
	var def core.AlarmDefinition
	err := r.pool.QueryRow(context.Background(), query, id).
		Scan(&def.ID, &def.Tag, &def.Threshold, &def.Type, &def.Priority, &def.CreatedAt, &def.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Or specific error
		}
		return nil, fmt.Errorf("failed to get definition: %w", err)
	}
	return &def, nil
}

func (r *PostgresRepository) ListDefinitions() ([]*core.AlarmDefinition, error) {
	query := `
		SELECT id, tag, threshold, alarm_type, priority, created_at, updated_at
		FROM alarm_definitions
	`
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to list definitions: %w", err)
	}
	defer rows.Close()

	var defs []*core.AlarmDefinition
	for rows.Next() {
		var def core.AlarmDefinition
		if err := rows.Scan(&def.ID, &def.Tag, &def.Threshold, &def.Type, &def.Priority, &def.CreatedAt, &def.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan definition: %w", err)
		}
		defs = append(defs, &def)
	}
	return defs, nil
}

func (r *PostgresRepository) GetDefinitionsByTag(tag string) ([]*core.AlarmDefinition, error) {
	query := `
		SELECT id, tag, threshold, alarm_type, priority, created_at, updated_at
		FROM alarm_definitions
		WHERE tag = $1
	`
	rows, err := r.pool.Query(context.Background(), query, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to get definitions by tag: %w", err)
	}
	defer rows.Close()

	var defs []*core.AlarmDefinition
	for rows.Next() {
		var def core.AlarmDefinition
		if err := rows.Scan(&def.ID, &def.Tag, &def.Threshold, &def.Type, &def.Priority, &def.CreatedAt, &def.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan definition: %w", err)
		}
		defs = append(defs, &def)
	}
	return defs, nil
}

func (r *PostgresRepository) CreateActiveAlarm(alarm *core.ActiveAlarm) error {
	query := `
		INSERT INTO active_alarms (definition_id, state, activation_time, ack_time, shelved_until, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(context.Background(), query, alarm.DefinitionID, alarm.State, alarm.ActivationTime, alarm.AckTime, alarm.ShelvedUntil, alarm.Value).
		Scan(&alarm.ID, &alarm.CreatedAt, &alarm.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create active alarm: %w", err)
	}
	return nil
}

func (r *PostgresRepository) UpdateActiveAlarmState(id int, state string) error {
	query := `
		UPDATE active_alarms
		SET state = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.pool.Exec(context.Background(), query, state, id)
	if err != nil {
		return fmt.Errorf("failed to update active alarm state: %w", err)
	}
	return nil
}

func (r *PostgresRepository) AckActiveAlarm(id int, ackTime time.Time) error {
	query := `
		UPDATE active_alarms
		SET state = 'AckActive', ack_time = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.pool.Exec(context.Background(), query, ackTime, id)
	if err != nil {
		return fmt.Errorf("failed to ack active alarm: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ShelveActiveAlarm(id int, shelvedUntil time.Time) error {
	query := `
		UPDATE active_alarms
		SET state = 'Shelved', shelved_until = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.pool.Exec(context.Background(), query, shelvedUntil, id)
	if err != nil {
		return fmt.Errorf("failed to shelve active alarm: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetActiveAlarms() ([]*core.ActiveAlarm, error) {
	query := `
		SELECT id, definition_id, state, activation_time, ack_time, shelved_until, value, created_at, updated_at
		FROM active_alarms
		WHERE state != 'Normal'
	`
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active alarms: %w", err)
	}
	defer rows.Close()

	var alarms []*core.ActiveAlarm
	for rows.Next() {
		var alarm core.ActiveAlarm
		if err := rows.Scan(&alarm.ID, &alarm.DefinitionID, &alarm.State, &alarm.ActivationTime, &alarm.AckTime, &alarm.ShelvedUntil, &alarm.Value, &alarm.CreatedAt, &alarm.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan active alarm: %w", err)
		}
		alarms = append(alarms, &alarm)
	}
	return alarms, nil
}
