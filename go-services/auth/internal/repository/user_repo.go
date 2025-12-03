package repository

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
}

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) CreateUser(user *User) error {
	// Validate Role
	validRoles := map[string]bool{
		"ADMIN":    true,
		"ENGINEER": true,
		"OPERATOR": true,
		"AUDITOR":  true,
		"SERVICE":  true,
	}
	if !validRoles[user.Role] {
		return errors.New("invalid role")
	}

	query := `
		INSERT INTO public.users (username, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.DB.QueryRow(query, user.Username, user.PasswordHash, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PostgresUserRepository) GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	user := &User{}
	err := r.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Or return a specific ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
