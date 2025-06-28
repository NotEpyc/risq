package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, role, startup_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.PasswordHash, user.Role, user.StartupID, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, name, password_hash, role, startup_id, created_at, updated_at
		FROM users WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.StartupID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.Errorf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, name, password_hash, role, startup_id, created_at, updated_at
		FROM users WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.StartupID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.Errorf("Failed to get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users 
		SET name = $2, password_hash = $3, role = $4, startup_id = $5, updated_at = $6
		WHERE id = $1
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		user.ID, user.Name, user.PasswordHash, user.Role, user.StartupID, user.UpdatedAt)

	if err != nil {
		logger.Errorf("Failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
