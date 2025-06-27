package startup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, startup *Startup) error
	GetByID(ctx context.Context, id uuid.UUID) (*Startup, error)
	Update(ctx context.Context, startup *Startup) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Startup, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, startup *Startup) error {
	query := `
		INSERT INTO startups (id, name, description, industry, funding_stage, location, founded_date, team_size, website, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	// Generate UUID if not provided
	if startup.ID == uuid.Nil {
		startup.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		startup.ID,
		startup.Name,
		startup.Description,
		startup.Industry,
		startup.FundingStage,
		startup.Location,
		startup.FoundedDate,
		startup.TeamSize,
		startup.Website,
		startup.CreatedAt,
		startup.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create startup: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Startup, error) {
	var startup Startup
	query := `
		SELECT id, name, description, industry, funding_stage, location, founded_date, team_size, website, created_at, updated_at
		FROM startups WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&startup.ID,
		&startup.Name,
		&startup.Description,
		&startup.Industry,
		&startup.FundingStage,
		&startup.Location,
		&startup.FoundedDate,
		&startup.TeamSize,
		&startup.Website,
		&startup.CreatedAt,
		&startup.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("startup not found")
		}
		return nil, fmt.Errorf("failed to get startup: %w", err)
	}
	return &startup, nil
}

func (r *repository) Update(ctx context.Context, startup *Startup) error {
	query := `
		UPDATE startups 
		SET name = $2, description = $3, industry = $4, funding_stage = $5, location = $6, 
		    founded_date = $7, team_size = $8, website = $9, updated_at = $10
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		startup.ID,
		startup.Name,
		startup.Description,
		startup.Industry,
		startup.FundingStage,
		startup.Location,
		startup.FoundedDate,
		startup.TeamSize,
		startup.Website,
		startup.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update startup: %w", err)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM startups WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete startup: %w", err)
	}
	return nil
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Startup, error) {
	var startup Startup
	query := `
		SELECT s.id, s.name, s.description, s.industry, s.funding_stage, s.location, 
		       s.founded_date, s.team_size, s.website, s.created_at, s.updated_at
		FROM startups s
		JOIN users u ON u.startup_id = s.id
		WHERE u.id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&startup.ID,
		&startup.Name,
		&startup.Description,
		&startup.Industry,
		&startup.FundingStage,
		&startup.Location,
		&startup.FoundedDate,
		&startup.TeamSize,
		&startup.Website,
		&startup.CreatedAt,
		&startup.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("startup not found")
		}
		return nil, fmt.Errorf("failed to get startup by user ID: %w", err)
	}
	return &startup, nil
}
