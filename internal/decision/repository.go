package decision

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, decision *Decision) error
	GetByID(ctx context.Context, id uuid.UUID) (*Decision, error)
	Update(ctx context.Context, decision *Decision) error
	GetByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error)
	GetSpeculativeByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, decision *Decision) error {
	query := `
		INSERT INTO decisions (id, startup_id, title, description, category, urgency, status, context, reasoning, created_at, updated_at, confirmed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// Generate UUID if not provided
	if decision.ID == uuid.Nil {
		decision.ID = uuid.New()
	}

	// Map fields to database schema
	title := decision.Description // Use description as title for now
	urgency := "medium"           // Default urgency

	_, err := r.db.ExecContext(ctx, query,
		decision.ID,
		decision.StartupID,
		title,
		decision.Description,
		decision.Category,
		urgency,
		decision.Status,
		decision.Context,
		decision.Reasoning,
		decision.CreatedAt,
		decision.UpdatedAt,
		decision.ConfirmedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create decision: %w", err)
	}
	return nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Decision, error) {
	var decision Decision
	var title, urgency string
	var aiAnalysis sql.NullString
	var finalChoice sql.NullString

	query := `
		SELECT id, startup_id, title, description, category, urgency, status, context,
		       ai_analysis, final_choice, reasoning, created_at, updated_at, confirmed_at
		FROM decisions WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&decision.ID,
		&decision.StartupID,
		&title,
		&decision.Description,
		&decision.Category,
		&urgency,
		&decision.Status,
		&decision.Context,
		&aiAnalysis,
		&finalChoice,
		&decision.Reasoning,
		&decision.CreatedAt,
		&decision.UpdatedAt,
		&decision.ConfirmedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("decision not found")
		}
		return nil, fmt.Errorf("failed to get decision: %w", err)
	}

	// Set default values for fields not in database
	decision.Timeline = ""
	decision.Budget = 0.0
	decision.PreviousRiskScore = 0.0
	decision.ProjectedRiskScore = 0.0
	decision.RiskDelta = 0.0
	decision.Confidence = 0.8
	decision.Suggestions = []string{}

	return &decision, nil
}

func (r *repository) Update(ctx context.Context, decision *Decision) error {
	query := `
		UPDATE decisions 
		SET startup_id = $2, title = $3, description = $4, category = $5, urgency = $6, 
		    status = $7, context = $8, reasoning = $9, updated_at = $10, confirmed_at = $11
		WHERE id = $1
	`

	title := decision.Description // Use description as title
	urgency := "medium"           // Default urgency

	_, err := r.db.ExecContext(ctx, query,
		decision.ID,
		decision.StartupID,
		title,
		decision.Description,
		decision.Category,
		urgency,
		decision.Status,
		decision.Context,
		decision.Reasoning,
		decision.UpdatedAt,
		decision.ConfirmedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update decision: %w", err)
	}
	return nil
}

func (r *repository) GetByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error) {
	query := `
		SELECT id, startup_id, title, description, category, urgency, status, context,
		       ai_analysis, final_choice, reasoning, created_at, updated_at, confirmed_at
		FROM decisions WHERE startup_id = $1 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get decisions by startup ID: %w", err)
	}
	defer rows.Close()

	var decisions []*Decision
	for rows.Next() {
		var decision Decision
		var title, urgency string
		var aiAnalysis sql.NullString
		var finalChoice sql.NullString

		err := rows.Scan(
			&decision.ID,
			&decision.StartupID,
			&title,
			&decision.Description,
			&decision.Category,
			&urgency,
			&decision.Status,
			&decision.Context,
			&aiAnalysis,
			&finalChoice,
			&decision.Reasoning,
			&decision.CreatedAt,
			&decision.UpdatedAt,
			&decision.ConfirmedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan decision: %w", err)
		}

		// Set default values for fields not in database
		decision.Timeline = ""
		decision.Budget = 0.0
		decision.PreviousRiskScore = 0.0
		decision.ProjectedRiskScore = 0.0
		decision.RiskDelta = 0.0
		decision.Confidence = 0.8
		decision.Suggestions = []string{}

		decisions = append(decisions, &decision)
	}

	return decisions, nil
}

func (r *repository) GetSpeculativeByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error) {
	query := `
		SELECT id, startup_id, title, description, category, urgency, status, context,
		       ai_analysis, final_choice, reasoning, created_at, updated_at, confirmed_at
		FROM decisions WHERE startup_id = $1 AND status = $2 ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startupID, "speculative")
	if err != nil {
		return nil, fmt.Errorf("failed to get speculative decisions: %w", err)
	}
	defer rows.Close()

	var decisions []*Decision
	for rows.Next() {
		var decision Decision
		var title, urgency string
		var aiAnalysis sql.NullString
		var finalChoice sql.NullString

		err := rows.Scan(
			&decision.ID,
			&decision.StartupID,
			&title,
			&decision.Description,
			&decision.Category,
			&urgency,
			&decision.Status,
			&decision.Context,
			&aiAnalysis,
			&finalChoice,
			&decision.Reasoning,
			&decision.CreatedAt,
			&decision.UpdatedAt,
			&decision.ConfirmedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan decision: %w", err)
		}

		// Set default values for fields not in database
		decision.Timeline = ""
		decision.Budget = 0.0
		decision.PreviousRiskScore = 0.0
		decision.ProjectedRiskScore = 0.0
		decision.RiskDelta = 0.0
		decision.Confidence = 0.8
		decision.Suggestions = []string{}

		decisions = append(decisions, &decision)
	}

	return decisions, nil
}
