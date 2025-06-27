package decision

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
		INSERT INTO decisions (id, startup_id, description, category, context, timeline, budget, status, 
		                      previous_risk_score, projected_risk_score, risk_delta, confidence, 
		                      suggestions, reasoning, created_at, updated_at, confirmed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	// Generate UUID if not provided
	if decision.ID == uuid.Nil {
		decision.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		decision.ID,
		decision.StartupID,
		decision.Description,
		decision.Category,
		decision.Context,
		decision.Timeline,
		decision.Budget,
		decision.Status,
		decision.PreviousRiskScore,
		decision.ProjectedRiskScore,
		decision.RiskDelta,
		decision.Confidence,
		pq.Array(decision.Suggestions),
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
	query := `
		SELECT id, startup_id, description, category, context, timeline, budget, status,
		       previous_risk_score, projected_risk_score, risk_delta, confidence,
		       suggestions, reasoning, created_at, updated_at, confirmed_at
		FROM decisions WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&decision.ID,
		&decision.StartupID,
		&decision.Description,
		&decision.Category,
		&decision.Context,
		&decision.Timeline,
		&decision.Budget,
		&decision.Status,
		&decision.PreviousRiskScore,
		&decision.ProjectedRiskScore,
		&decision.RiskDelta,
		&decision.Confidence,
		pq.Array(&decision.Suggestions),
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
	return &decision, nil
}

func (r *repository) Update(ctx context.Context, decision *Decision) error {
	query := `
		UPDATE decisions 
		SET startup_id = $2, description = $3, category = $4, context = $5, timeline = $6, budget = $7, 
		    status = $8, previous_risk_score = $9, projected_risk_score = $10, risk_delta = $11, 
		    confidence = $12, suggestions = $13, reasoning = $14, updated_at = $15, confirmed_at = $16
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		decision.ID,
		decision.StartupID,
		decision.Description,
		decision.Category,
		decision.Context,
		decision.Timeline,
		decision.Budget,
		decision.Status,
		decision.PreviousRiskScore,
		decision.ProjectedRiskScore,
		decision.RiskDelta,
		decision.Confidence,
		pq.Array(decision.Suggestions),
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
		SELECT id, startup_id, description, category, context, timeline, budget, status,
		       previous_risk_score, projected_risk_score, risk_delta, confidence,
		       suggestions, reasoning, created_at, updated_at, confirmed_at
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
		err := rows.Scan(
			&decision.ID,
			&decision.StartupID,
			&decision.Description,
			&decision.Category,
			&decision.Context,
			&decision.Timeline,
			&decision.Budget,
			&decision.Status,
			&decision.PreviousRiskScore,
			&decision.ProjectedRiskScore,
			&decision.RiskDelta,
			&decision.Confidence,
			pq.Array(&decision.Suggestions),
			&decision.Reasoning,
			&decision.CreatedAt,
			&decision.UpdatedAt,
			&decision.ConfirmedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan decision: %w", err)
		}
		decisions = append(decisions, &decision)
	}

	return decisions, nil
}

func (r *repository) GetSpeculativeByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error) {
	query := `
		SELECT id, startup_id, description, category, context, timeline, budget, status,
		       previous_risk_score, projected_risk_score, risk_delta, confidence,
		       suggestions, reasoning, created_at, updated_at, confirmed_at
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
		err := rows.Scan(
			&decision.ID,
			&decision.StartupID,
			&decision.Description,
			&decision.Category,
			&decision.Context,
			&decision.Timeline,
			&decision.Budget,
			&decision.Status,
			&decision.PreviousRiskScore,
			&decision.ProjectedRiskScore,
			&decision.RiskDelta,
			&decision.Confidence,
			pq.Array(&decision.Suggestions),
			&decision.Reasoning,
			&decision.CreatedAt,
			&decision.UpdatedAt,
			&decision.ConfirmedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan decision: %w", err)
		}
		decisions = append(decisions, &decision)
	}

	return decisions, nil
}
