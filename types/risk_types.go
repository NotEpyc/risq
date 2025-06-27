package types

import (
	"time"

	"github.com/google/uuid"
)

// Risk-related types
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

type RiskScore struct {
	Score       float64   `json:"score" db:"score"`
	Level       RiskLevel `json:"level" db:"level"`
	Confidence  float64   `json:"confidence" db:"confidence"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
}

type RiskFactor struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	Impact      float64   `json:"impact" db:"impact"`
	Probability float64   `json:"probability" db:"probability"`
	Weight      float64   `json:"weight" db:"weight"`
}

type RiskProfile struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	StartupID   uuid.UUID    `json:"startup_id" db:"startup_id"`
	Score       RiskScore    `json:"score"`
	Factors     []RiskFactor `json:"factors"`
	Suggestions []string     `json:"suggestions"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}

// Risk evolution for timeline
type RiskEvolution struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StartupID uuid.UUID `json:"startup_id" db:"startup_id"`
	Score     float64   `json:"score" db:"score"`
	Level     RiskLevel `json:"level" db:"level"`
	Trigger   string    `json:"trigger" db:"trigger"` // what caused this change
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
