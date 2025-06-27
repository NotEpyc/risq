package startup

import (
	"time"

	"github.com/google/uuid"
)

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
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Industry     string       `json:"industry"`
	FundingStage FundingStage `json:"funding_stage"`
	Location     string       `json:"location"`
	FoundedDate  time.Time    `json:"founded_date"`
	TeamSize     int          `json:"team_size"`
	Website      string       `json:"website,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

func (Startup) TableName() string {
	return "startups"
}

type StartupOnboardingInput struct {
	// Basic startup information
	Name         string       `json:"name" validate:"required"`
	Description  string       `json:"description" validate:"required"`
	Industry     string       `json:"industry" validate:"required"`
	Sector       string       `json:"sector" validate:"required"` // e.g., "edutech", "logistics", "fintech"
	FundingStage FundingStage `json:"funding_stage" validate:"required"`
	Location     string       `json:"location" validate:"required"`
	FoundedDate  string       `json:"founded_date" validate:"required"`
	TeamSize     int          `json:"team_size" validate:"required"`
	Website      string       `json:"website,omitempty"`

	// Business model information
	BusinessModel      string   `json:"business_model" validate:"required"`
	RevenueStreams     []string `json:"revenue_streams" validate:"required"`
	TargetMarket       string   `json:"target_market" validate:"required"`
	CompetitorAnalysis string   `json:"competitor_analysis"`

	// Implementation and strategy (simplified to strings)
	ImplementationPlan  string   `json:"implementation_plan" validate:"required"`
	TechnologyStack     []string `json:"technology_stack"`
	DevelopmentTimeline string   `json:"development_timeline"`
	GoToMarketStrategy  string   `json:"go_to_market_strategy"`

	// Financial information
	InitialInvestment  float64 `json:"initial_investment"`
	MonthlyBurnRate    float64 `json:"monthly_burn_rate"`
	ProjectedRevenue   float64 `json:"projected_revenue"`
	FundingRequirement float64 `json:"funding_requirement"`

	// Founder information (simplified)
	FounderDetails []FounderProfile `json:"founder_details" validate:"required"`
}

// Simplified founder profile for easier JSON handling
type FounderProfile struct {
	Name             string           `json:"name" validate:"required"`
	Email            string           `json:"email" validate:"required,email"`
	Role             string           `json:"role" validate:"required"`
	LinkedInURL      string           `json:"linkedin_url"`
	Education        []string         `json:"education"` // Simplified to string array
	Experience       []WorkExperience `json:"experience"`
	Skills           []string         `json:"skills"`
	Achievements     []string         `json:"achievements"`
	PreviousStartups []string         `json:"previous_startups"` // Simplified to string array
}

// Internal complex types for detailed analysis (renamed to avoid conflicts)
type DetailedFounderCV struct {
	Name             string            `json:"name" validate:"required"`
	Email            string            `json:"email" validate:"required,email"`
	Education        []Education       `json:"education"`
	WorkExperience   []WorkExperience  `json:"work_experience"`
	Skills           []string          `json:"skills"`
	Achievements     []string          `json:"achievements"`
	LinkedInURL      string            `json:"linkedin_url"`
	PreviousStartups []PreviousStartup `json:"previous_startups"`
}

type Education struct {
	Institution    string  `json:"institution"`
	Degree         string  `json:"degree"`
	Field          string  `json:"field"`
	GraduationYear int     `json:"graduation_year"`
	GPA            float64 `json:"gpa,omitempty"`
}

type WorkExperience struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
}

type PreviousStartup struct {
	Name        string `json:"name"`
	Industry    string `json:"industry"`
	Role        string `json:"role"`
	Outcome     string `json:"outcome"` // "successful_exit", "failed", "ongoing"
	Description string `json:"description"`
}

type DetailedBusinessPlan struct {
	ExecutiveSummary     string               `json:"executive_summary" validate:"required"`
	ProblemStatement     string               `json:"problem_statement" validate:"required"`
	SolutionDescription  string               `json:"solution_description" validate:"required"`
	MarketAnalysis       MarketAnalysis       `json:"market_analysis" validate:"required"`
	FinancialProjections FinancialProjections `json:"financial_projections" validate:"required"`
	FundingRequirements  FundingRequirements  `json:"funding_requirements"`
	RiskAssessment       []string             `json:"risk_assessment"`
}

type MarketAnalysis struct {
	TotalAddressableMarket int64        `json:"total_addressable_market"`
	ServiceableMarket      int64        `json:"serviceable_market"`
	TargetCustomers        []string     `json:"target_customers"`
	MarketTrends           []string     `json:"market_trends"`
	CompetitiveLandscape   []Competitor `json:"competitive_landscape"`
}

type Competitor struct {
	Name        string   `json:"name"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	MarketShare float64  `json:"market_share"`
}

type FinancialProjections struct {
	Year1Revenue    int64 `json:"year1_revenue"`
	Year2Revenue    int64 `json:"year2_revenue"`
	Year3Revenue    int64 `json:"year3_revenue"`
	BreakEvenMonth  int   `json:"break_even_month"`
	InitialCosts    int64 `json:"initial_costs"`
	MonthlyBurnRate int64 `json:"monthly_burn_rate"`
}

type FundingRequirements struct {
	TotalFunding   int64    `json:"total_funding"`
	FundingPurpose []string `json:"funding_purpose"`
	Timeline       string   `json:"timeline"`
	InvestorType   string   `json:"investor_type"` // "angel", "vc", "seed", "series_a"
}

type DetailedImplementationPlan struct {
	Phase1         Phase          `json:"phase1" validate:"required"`
	Phase2         Phase          `json:"phase2"`
	Phase3         Phase          `json:"phase3"`
	TechnicalStack TechnicalStack `json:"technical_stack"`
	Timeline       string         `json:"timeline" validate:"required"`
	Milestones     []Milestone    `json:"milestones"`
}

type Phase struct {
	Name         string   `json:"name"`
	Duration     string   `json:"duration"`
	Objectives   []string `json:"objectives"`
	Deliverables []string `json:"deliverables"`
	Resources    []string `json:"resources"`
}

type TechnicalStack struct {
	Frontend       []string `json:"frontend"`
	Backend        []string `json:"backend"`
	Database       []string `json:"database"`
	Infrastructure []string `json:"infrastructure"`
	ThirdPartyAPIs []string `json:"third_party_apis"`
}

type Milestone struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Success     string `json:"success_criteria"`
}
