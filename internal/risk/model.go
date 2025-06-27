package risk

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"
	RiskLevelMedium   RiskLevel = "medium"
	RiskLevelHigh     RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

type RiskProfile struct {
	ID          uuid.UUID      `json:"id"`
	StartupID   uuid.UUID      `json:"startup_id"`
	Score       float64        `json:"score"`
	Level       RiskLevel      `json:"level"`
	Confidence  float64        `json:"confidence"`
	Factors     pq.StringArray `json:"factors"`
	Suggestions pq.StringArray `json:"suggestions"`
	Reasoning   string         `json:"reasoning"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (RiskProfile) TableName() string {
	return "risk_profiles"
}

// Risk evolution for timeline
type RiskEvolution struct {
	ID        uuid.UUID `json:"id"`
	StartupID uuid.UUID `json:"startup_id"`
	Score     float64   `json:"score"`
	Level     RiskLevel `json:"level"`
	Trigger   string    `json:"trigger"` // what caused this change
	CreatedAt time.Time `json:"created_at"`
}

func (RiskEvolution) TableName() string {
	return "risk_evolutions"
}

func DetermineRiskLevel(score float64) RiskLevel {
	switch {
	case score < 25:
		return RiskLevelLow
	case score < 50:
		return RiskLevelMedium
	case score < 75:
		return RiskLevelHigh
	default:
		return RiskLevelCritical
	}
}
