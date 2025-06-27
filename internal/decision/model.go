package decision

import (
	"risq_backend/types"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Decision struct {
	ID                 uuid.UUID              `json:"id"`
	StartupID          uuid.UUID              `json:"startup_id"`
	Description        string                 `json:"description"`
	Category           types.DecisionCategory `json:"category"`
	Context            string                 `json:"context,omitempty"`
	Timeline           string                 `json:"timeline,omitempty"`
	Budget             float64                `json:"budget,omitempty"`
	Status             types.DecisionStatus   `json:"status"`
	PreviousRiskScore  float64                `json:"previous_risk_score"`
	ProjectedRiskScore float64                `json:"projected_risk_score"`
	RiskDelta          float64                `json:"risk_delta"`
	Confidence         float64                `json:"confidence"`
	Suggestions        pq.StringArray         `json:"suggestions"`
	Reasoning          string                 `json:"reasoning"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	ConfirmedAt        *time.Time             `json:"confirmed_at,omitempty"`
}

func (Decision) TableName() string {
	return "decisions"
}

// Decision response for API
type DecisionResponse struct {
	Decision    Decision `json:"decision"`
	RiskScore   float64  `json:"risk_score"`
	RiskDelta   float64  `json:"risk_delta"`
	Suggestions []string `json:"suggestions"`
	Reasoning   string   `json:"reasoning"`
	Confidence  float64  `json:"confidence"`
}

// Decision confirmation input
type DecisionConfirmInput struct {
	DecisionID uuid.UUID `json:"decision_id" validate:"required"`
	Notes      string    `json:"notes,omitempty"`
}
