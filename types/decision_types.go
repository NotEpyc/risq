package types

import (
	"time"

	"github.com/google/uuid"
)

// Decision-related types
type DecisionStatus string

const (
	DecisionStatusSpeculative DecisionStatus = "speculative"
	DecisionStatusConfirmed   DecisionStatus = "confirmed"
	DecisionStatusRejected    DecisionStatus = "rejected"
)

type DecisionCategory string

const (
	DecisionCategoryHiring     DecisionCategory = "hiring"
	DecisionCategoryFunding    DecisionCategory = "funding"
	DecisionCategoryProduct    DecisionCategory = "product"
	DecisionCategoryMarketing  DecisionCategory = "marketing"
	DecisionCategoryOperations DecisionCategory = "operations"
	DecisionCategoryStrategy   DecisionCategory = "strategy"
	DecisionCategoryLegal      DecisionCategory = "legal"
	DecisionCategoryOther      DecisionCategory = "other"
)

// Input for decision speculation
type DecisionInput struct {
	StartupID   uuid.UUID        `json:"startup_id" validate:"required"`
	Description string           `json:"description" validate:"required"`
	Category    DecisionCategory `json:"category" validate:"required"`
	Context     string           `json:"context,omitempty"`
	Timeline    string           `json:"timeline,omitempty"`
	Budget      float64          `json:"budget,omitempty"`
}

// Decision result from speculation
type DecisionResult struct {
	ID                 uuid.UUID        `json:"id" db:"id"`
	StartupID          uuid.UUID        `json:"startup_id" db:"startup_id"`
	Description        string           `json:"description" db:"description"`
	Category           DecisionCategory `json:"category" db:"category"`
	Context            string           `json:"context,omitempty" db:"context"`
	Timeline           string           `json:"timeline,omitempty" db:"timeline"`
	Budget             float64          `json:"budget,omitempty" db:"budget"`
	Status             DecisionStatus   `json:"status" db:"status"`
	PreviousRiskScore  float64          `json:"previous_risk_score" db:"previous_risk_score"`
	ProjectedRiskScore float64          `json:"projected_risk_score" db:"projected_risk_score"`
	RiskDelta          float64          `json:"risk_delta" db:"risk_delta"`
	Confidence         float64          `json:"confidence" db:"confidence"`
	Suggestions        []string         `json:"suggestions" db:"suggestions"`
	Reasoning          string           `json:"reasoning" db:"reasoning"`
	CreatedAt          time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at" db:"updated_at"`
	ConfirmedAt        *time.Time       `json:"confirmed_at,omitempty" db:"confirmed_at"`
}

// Decision response for API
type DecisionResponse struct {
	Decision    DecisionResult `json:"decision"`
	RiskScore   float64        `json:"risk_score"`
	RiskDelta   float64        `json:"risk_delta"`
	Suggestions []string       `json:"suggestions"`
	Reasoning   string         `json:"reasoning"`
	Confidence  float64        `json:"confidence"`
}

// Decision confirmation input
type DecisionConfirmInput struct {
	DecisionID uuid.UUID `json:"decision_id" validate:"required"`
	Notes      string    `json:"notes,omitempty"`
}
