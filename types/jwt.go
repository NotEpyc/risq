package types

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWT Claims
type Claims struct {
	UserID    uuid.UUID `json:"user_id"`
	StartupID uuid.UUID `json:"startup_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	jwt.RegisteredClaims
}

// User types
type User struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	Role      string     `json:"role" db:"role"`
	StartupID *uuid.UUID `json:"startup_id,omitempty" db:"startup_id"`
	CreatedAt string     `json:"created_at" db:"created_at"`
	UpdatedAt string     `json:"updated_at" db:"updated_at"`
}

// Startup types
type FundingStage string

const (
	FundingStageIdea    FundingStage = "idea"
	FundingStagePreSeed FundingStage = "pre_seed"
	FundingStageSeed    FundingStage = "seed"
	FundingStageSeriesA FundingStage = "series_a"
	FundingStageSeriesB FundingStage = "series_b"
	FundingStageSeriesC FundingStage = "series_c"
	FundingStageIPO     FundingStage = "ipo"
)

type Startup struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	Name         string       `json:"name" db:"name"`
	Description  string       `json:"description" db:"description"`
	Industry     string       `json:"industry" db:"industry"`
	FundingStage FundingStage `json:"funding_stage" db:"funding_stage"`
	Location     string       `json:"location" db:"location"`
	FoundedDate  string       `json:"founded_date" db:"founded_date"`
	TeamSize     int          `json:"team_size" db:"team_size"`
	Website      string       `json:"website,omitempty" db:"website"`
	CreatedAt    string       `json:"created_at" db:"created_at"`
	UpdatedAt    string       `json:"updated_at" db:"updated_at"`
}
