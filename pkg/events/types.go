package events

import (
	"time"

	"github.com/google/uuid"
)

// Event base structure
type BaseEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Subject   string                 `json:"subject"`
	Timestamp time.Time              `json:"timestamp"`
	Version   string                 `json:"version"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Startup onboarding events

// StartupSubmittedEvent represents the initial startup submission
type StartupSubmittedEvent struct {
	BaseEvent
	StartupID uuid.UUID             `json:"startup_id"`
	UserID    uuid.UUID             `json:"user_id"`
	Data      StartupOnboardingData `json:"data"`
}

// StartupOnboardingData contains comprehensive startup information
type StartupOnboardingData struct {
	// Basic Information
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"required,min=10,max=1000"`
	Website     string `json:"website" validate:"url"`

	// Startup Details
	Industry     string    `json:"industry" validate:"required"`
	Sector       string    `json:"sector" validate:"required"` // logistics, edutech, fintech, etc.
	FundingStage string    `json:"funding_stage" validate:"required"`
	Location     string    `json:"location" validate:"required"`
	FoundedDate  time.Time `json:"founded_date" validate:"required"`
	TeamSize     int       `json:"team_size" validate:"required,min=1,max=10000"`

	// Business Model
	BusinessModel      string   `json:"business_model" validate:"required"` // B2B, B2C, B2B2C, etc.
	RevenueStreams     []string `json:"revenue_streams" validate:"required,min=1"`
	TargetMarket       string   `json:"target_market" validate:"required"`
	CompetitorAnalysis string   `json:"competitor_analysis" validate:"required,min=50,max=2000"`

	// Implementation Strategy
	ImplementationPlan  string   `json:"implementation_plan" validate:"required,min=100,max=5000"`
	TechnologyStack     []string `json:"technology_stack" validate:"required,min=1"`
	DevelopmentTimeline string   `json:"development_timeline" validate:"required"`
	GoToMarketStrategy  string   `json:"go_to_market_strategy" validate:"required,min=50,max=2000"`

	// Financial Information
	InitialInvestment  float64 `json:"initial_investment" validate:"min=0"`
	MonthlyBurnRate    float64 `json:"monthly_burn_rate" validate:"min=0"`
	ProjectedRevenue   float64 `json:"projected_revenue" validate:"min=0"`
	FundingRequirement float64 `json:"funding_requirement" validate:"min=0"`

	// Founder Information
	FounderDetails []FounderProfile `json:"founder_details" validate:"required,min=1"`
}

// FounderProfile contains founder's background information
type FounderProfile struct {
	Name             string           `json:"name" validate:"required,min=2,max=100"`
	Email            string           `json:"email" validate:"required,email"`
	Role             string           `json:"role" validate:"required"`
	LinkedInURL      string           `json:"linkedin_url" validate:"url"`
	Education        []string         `json:"education" validate:"required,min=1"`
	Experience       []WorkExperience `json:"experience" validate:"required,min=1"`
	Skills           []string         `json:"skills" validate:"required,min=3"`
	Achievements     []string         `json:"achievements"`
	PreviousStartups []string         `json:"previous_startups"`
}

// WorkExperience represents professional experience
type WorkExperience struct {
	Company     string     `json:"company" validate:"required"`
	Position    string     `json:"position" validate:"required"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"` // nil for current position
	Description string     `json:"description" validate:"required,min=20,max=500"`
	Industry    string     `json:"industry" validate:"required"`
}

// StartupValidatedEvent represents successful validation
type StartupValidatedEvent struct {
	BaseEvent
	StartupID         uuid.UUID         `json:"startup_id"`
	UserID            uuid.UUID         `json:"user_id"`
	ValidationResults ValidationResults `json:"validation_results"`
}

// ValidationResults contains validation outcomes
type ValidationResults struct {
	IsValid         bool               `json:"is_valid"`
	ValidationScore float64            `json:"validation_score"` // 0-100
	PassedChecks    []string           `json:"passed_checks"`
	FailedChecks    []string           `json:"failed_checks"`
	Warnings        []string           `json:"warnings"`
	RequiredActions []string           `json:"required_actions"`
	QualityMetrics  map[string]float64 `json:"quality_metrics"`
}

// MarketDataFetchedEvent represents market data retrieval
type MarketDataFetchedEvent struct {
	BaseEvent
	StartupID  uuid.UUID  `json:"startup_id"`
	UserID     uuid.UUID  `json:"user_id"`
	MarketData MarketData `json:"market_data"`
}

// MarketData contains market analysis information
type MarketData struct {
	Industry         string           `json:"industry"`
	Sector           string           `json:"sector"`
	MarketSize       MarketSizeData   `json:"market_size"`
	GrowthRate       float64          `json:"growth_rate"`       // annual percentage
	MarketStatus     string           `json:"market_status"`     // active, declining, emerging, mature
	CompetitionLevel string           `json:"competition_level"` // low, medium, high
	KeyTrends        []string         `json:"key_trends"`
	Opportunities    []string         `json:"opportunities"`
	Threats          []string         `json:"threats"`
	RegulationLevel  string           `json:"regulation_level"` // low, medium, high
	BarriersToEntry  []string         `json:"barriers_to_entry"`
	NewsAnalysis     NewsAnalysisData `json:"news_analysis"`
	LastUpdated      time.Time        `json:"last_updated"`
}

// MarketSizeData contains market size information
type MarketSizeData struct {
	TAM  float64 `json:"tam"`  // Total Addressable Market
	SAM  float64 `json:"sam"`  // Serviceable Addressable Market
	SOM  float64 `json:"som"`  // Serviceable Obtainable Market
	Unit string  `json:"unit"` // USD, users, etc.
}

// NewsAnalysisData contains sentiment analysis from news
type NewsAnalysisData struct {
	SentimentScore   float64  `json:"sentiment_score"` // -1 to 1
	PositiveKeywords []string `json:"positive_keywords"`
	NegativeKeywords []string `json:"negative_keywords"`
	RecentNews       []string `json:"recent_news"`
	InvestmentNews   []string `json:"investment_news"`
}

// RiskAnalysisRequestedEvent represents risk analysis request - compatible with events.go interface
type RiskAnalysisRequestedEvent struct {
	BaseEvent
	StartupID      uuid.UUID              `json:"startup_id"`
	StartupData    map[string]interface{} `json:"startup_data"`
	FounderCV      map[string]interface{} `json:"founder_cv"`
	BusinessPlan   map[string]interface{} `json:"business_plan"`
	MarketData     map[string]interface{} `json:"market_data"`
	SectorAnalysis map[string]interface{} `json:"sector_analysis"`
}

// RiskAnalysisCompletedEvent represents completed risk analysis - compatible with events.go interface
type RiskAnalysisCompletedEvent struct {
	BaseEvent
	StartupID        uuid.UUID              `json:"startup_id"`
	RiskScore        float64                `json:"risk_score"`
	RiskLevel        string                 `json:"risk_level"`
	Strengths        []string               `json:"strengths"`
	Weaknesses       []string               `json:"weaknesses"`
	Recommendations  []string               `json:"recommendations"`
	DetailedAnalysis map[string]interface{} `json:"detailed_analysis"`
}

// Structured versions for internal processing (renamed to avoid conflicts)
type RiskAnalysisRequestedEventStructured struct {
	BaseEvent
	StartupID    uuid.UUID             `json:"startup_id"`
	UserID       uuid.UUID             `json:"user_id"`
	StartupData  StartupOnboardingData `json:"startup_data"`
	MarketData   MarketData            `json:"market_data"`
	AnalysisType string                `json:"analysis_type"` // initial, update, deep_dive
}

type RiskAnalysisCompletedEventStructured struct {
	BaseEvent
	StartupID    uuid.UUID    `json:"startup_id"`
	UserID       uuid.UUID    `json:"user_id"`
	RiskAnalysis RiskAnalysis `json:"risk_analysis"`
}

// RiskAnalysis contains comprehensive risk assessment
type RiskAnalysis struct {
	OverallRiskScore float64            `json:"overall_risk_score"` // 0-100
	Confidence       float64            `json:"confidence"`         // 0-1
	RiskCategories   map[string]float64 `json:"risk_categories"`
	Strengths        []StrengthWeakness `json:"strengths"`
	Weaknesses       []StrengthWeakness `json:"weaknesses"`
	Opportunities    []StrengthWeakness `json:"opportunities"`
	Threats          []StrengthWeakness `json:"threats"`
	Recommendations  []Recommendation   `json:"recommendations"`
	CompetitorRisks  []CompetitorRisk   `json:"competitor_risks"`
	FinancialRisks   []FinancialRisk    `json:"financial_risks"`
	Summary          string             `json:"summary"`
	CreatedAt        time.Time          `json:"created_at"`
}

// StrengthWeakness represents SWOT analysis items
type StrengthWeakness struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"`   // low, medium, high
	Priority    int    `json:"priority"` // 1-5
	Category    string `json:"category"`
}

// Recommendation represents actionable advice
type Recommendation struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Priority       string   `json:"priority"` // low, medium, high, critical
	Timeline       string   `json:"timeline"` // immediate, short_term, long_term
	Effort         string   `json:"effort"`   // low, medium, high
	ExpectedImpact string   `json:"expected_impact"`
	Resources      []string `json:"resources"`
	KPIs           []string `json:"kpis"`
	Implementation []string `json:"implementation_steps"`
}

// CompetitorRisk represents competitive threats
type CompetitorRisk struct {
	CompetitorName string   `json:"competitor_name"`
	ThreatLevel    string   `json:"threat_level"` // low, medium, high
	Advantages     []string `json:"advantages"`
	Weaknesses     []string `json:"weaknesses"`
	MarketShare    float64  `json:"market_share"`
}

// FinancialRisk represents financial concerns
type FinancialRisk struct {
	RiskType    string   `json:"risk_type"`
	Probability float64  `json:"probability"` // 0-1
	Impact      float64  `json:"impact"`      // monetary impact
	Mitigation  []string `json:"mitigation_strategies"`
	Timeline    string   `json:"timeline"`
}

// ContextStorageRequestedEvent represents context storage request
type ContextStorageRequestedEvent struct {
	BaseEvent
	StartupID    uuid.UUID             `json:"startup_id"`
	UserID       uuid.UUID             `json:"user_id"`
	RiskAnalysis RiskAnalysis          `json:"risk_analysis"`
	StartupData  StartupOnboardingData `json:"startup_data"`
	MarketData   MarketData            `json:"market_data"`
}

// StartupOnboardingCompletedEvent represents successful onboarding
type StartupOnboardingCompletedEvent struct {
	BaseEvent
	StartupID uuid.UUID `json:"startup_id"`
	UserID    uuid.UUID `json:"user_id"`
	RiskScore float64   `json:"risk_score"`
	Summary   string    `json:"summary"`
}

// StartupOnboardedEvent represents the initial startup submission (alias for StartupSubmittedEvent)
type StartupOnboardedEvent struct {
	BaseEvent
	StartupID    uuid.UUID              `json:"startup_id"`
	UserID       uuid.UUID              `json:"user_id"`
	StartupData  map[string]interface{} `json:"startup_data"`
	FounderCV    map[string]interface{} `json:"founder_cv"`
	BusinessPlan map[string]interface{} `json:"business_plan"`
}

// MarketValidationRequestedEvent represents market validation request
type MarketValidationRequestedEvent struct {
	BaseEvent
	StartupID     uuid.UUID `json:"startup_id"`
	Industry      string    `json:"industry"`
	Sector        string    `json:"sector"`
	TargetMarket  string    `json:"target_market"`
	BusinessModel string    `json:"business_model"`
}

// MarketValidatedEvent represents market validation completion
type MarketValidatedEvent struct {
	BaseEvent
	StartupID       uuid.UUID              `json:"startup_id"`
	MarketData      map[string]interface{} `json:"market_data"`
	SectorAnalysis  map[string]interface{} `json:"sector_analysis"`
	MarketHealth    string                 `json:"market_health"` // "active", "inactive", "declining", "growing"
	Recommendations []string               `json:"recommendations"`
}

// ContextStoreRequestedEvent represents context storage request
type ContextStoreRequestedEvent struct {
	BaseEvent
	StartupID   uuid.UUID              `json:"startup_id"`
	ContentType string                 `json:"content_type"` // "startup_profile", "risk_analysis", "market_data"
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Utility function to create base event
func NewBaseEvent(source string) BaseEvent {
	return BaseEvent{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Source:    source,
		Version:   "1.0",
	}
}
