package risk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile *RiskProfile) error
	GetCurrentProfile(ctx context.Context, startupID uuid.UUID) (*RiskProfile, error)
	UpdateProfile(ctx context.Context, profile *RiskProfile) error
	CreateEvolution(ctx context.Context, evolution *RiskEvolution) error
	GetEvolutionHistory(ctx context.Context, startupID uuid.UUID, limit int) ([]*RiskEvolution, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateProfile(ctx context.Context, profile *RiskProfile) error {
	query := `
		INSERT INTO risk_profiles (id, startup_id, overall_risk, market_risk, technical_risk, financial_risk, team_risk, regulatory_risk, risk_factors, recommendations, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// Generate UUID if not provided
	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}

	// Calculate overall risk level from score
	overallRisk := "medium"
	if profile.Score >= 75 {
		overallRisk = "high"
	} else if profile.Score <= 25 {
		overallRisk = "low"
	}

	// Default individual risk scores to 0.5 (medium)
	marketRisk := 0.5
	technicalRisk := 0.5
	financialRisk := 0.5
	teamRisk := 0.5
	regulatoryRisk := 0.5

	_, err := r.db.ExecContext(ctx, query,
		profile.ID,
		profile.StartupID,
		overallRisk,
		marketRisk,
		technicalRisk,
		financialRisk,
		teamRisk,
		regulatoryRisk,
		pq.Array(profile.Factors),
		pq.Array(profile.Suggestions),
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create risk profile: %w", err)
	}
	return nil
}

func (r *repository) GetCurrentProfile(ctx context.Context, startupID uuid.UUID) (*RiskProfile, error) {
	var profile RiskProfile
	var overallRisk string
	var marketRisk, technicalRisk, financialRisk, teamRisk, regulatoryRisk float64

	query := `
		SELECT id, startup_id, overall_risk, market_risk, technical_risk, financial_risk, team_risk, regulatory_risk, risk_factors, recommendations, created_at, updated_at
		FROM risk_profiles WHERE startup_id = $1 ORDER BY created_at DESC LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, startupID).Scan(
		&profile.ID,
		&profile.StartupID,
		&overallRisk,
		&marketRisk,
		&technicalRisk,
		&financialRisk,
		&teamRisk,
		&regulatoryRisk,
		(*pq.StringArray)(&profile.Factors),
		(*pq.StringArray)(&profile.Suggestions),
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("risk profile not found")
		}
		return nil, fmt.Errorf("failed to get risk profile: %w", err)
	}

	// Calculate overall score from individual risk components
	avgRisk := (marketRisk + technicalRisk + financialRisk + teamRisk + regulatoryRisk) / 5
	profile.Score = avgRisk * 100 // Convert to 0-100 scale
	profile.Level = RiskLevel(overallRisk)
	profile.Confidence = 0.8 // Default confidence

	return &profile, nil
}

func (r *repository) UpdateProfile(ctx context.Context, profile *RiskProfile) error {
	query := `
		UPDATE risk_profiles 
		SET startup_id = $2, overall_risk = $3, market_risk = $4, technical_risk = $5, 
		    financial_risk = $6, team_risk = $7, regulatory_risk = $8, risk_factors = $9, 
		    recommendations = $10, updated_at = $11
		WHERE id = $1
	`

	// Calculate overall risk level from score
	overallRisk := "medium"
	if profile.Score >= 75 {
		overallRisk = "high"
	} else if profile.Score <= 25 {
		overallRisk = "low"
	}

	// Default individual risk scores to 0.5 (medium)
	marketRisk := 0.5
	technicalRisk := 0.5
	financialRisk := 0.5
	teamRisk := 0.5
	regulatoryRisk := 0.5

	_, err := r.db.ExecContext(ctx, query,
		profile.ID,
		profile.StartupID,
		overallRisk,
		marketRisk,
		technicalRisk,
		financialRisk,
		teamRisk,
		regulatoryRisk,
		pq.Array(profile.Factors),
		pq.Array(profile.Suggestions),
		profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update risk profile: %w", err)
	}
	return nil
}

func (r *repository) CreateEvolution(ctx context.Context, evolution *RiskEvolution) error {
	query := `
		INSERT INTO risk_evolutions (id, startup_id, score, level, trigger, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	// Generate UUID if not provided
	if evolution.ID == uuid.Nil {
		evolution.ID = uuid.New()
	}

	_, err := r.db.ExecContext(ctx, query,
		evolution.ID,
		evolution.StartupID,
		evolution.Score,
		evolution.Level,
		evolution.Trigger,
		evolution.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create risk evolution: %w", err)
	}
	return nil
}

func (r *repository) GetEvolutionHistory(ctx context.Context, startupID uuid.UUID, limit int) ([]*RiskEvolution, error) {
	baseQuery := `
		SELECT id, startup_id, score, level, trigger, created_at
		FROM risk_evolutions WHERE startup_id = $1 ORDER BY created_at DESC
	`

	var query string
	var args []interface{}

	if limit > 0 {
		query = baseQuery + " LIMIT $2"
		args = []interface{}{startupID, limit}
	} else {
		query = baseQuery
		args = []interface{}{startupID}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get risk evolution history: %w", err)
	}
	defer rows.Close()

	var evolutions []*RiskEvolution
	for rows.Next() {
		var evolution RiskEvolution
		err := rows.Scan(
			&evolution.ID,
			&evolution.StartupID,
			&evolution.Score,
			&evolution.Level,
			&evolution.Trigger,
			&evolution.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan risk evolution: %w", err)
		}
		evolutions = append(evolutions, &evolution)
	}

	return evolutions, nil
}
