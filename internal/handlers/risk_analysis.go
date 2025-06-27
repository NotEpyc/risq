package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"risq_backend/internal/llm"
	"risq_backend/internal/risk"
	"risq_backend/pkg/events"
	"risq_backend/pkg/logger"

	"github.com/nats-io/nats.go"
)

// RiskAnalysisHandler handles risk analysis events
type RiskAnalysisHandler struct {
	riskService  risk.Service
	llmService   llm.Service
	eventService events.EventService
}

func NewRiskAnalysisHandler(riskService risk.Service, llmService llm.Service, eventService events.EventService) *RiskAnalysisHandler {
	return &RiskAnalysisHandler{
		riskService:  riskService,
		llmService:   llmService,
		eventService: eventService,
	}
}

// HandleMarketValidated processes market validation completion and triggers risk analysis
func (h *RiskAnalysisHandler) HandleMarketValidated(msg *nats.Msg) {
	logger.Info("Processing market validated event for risk analysis")

	var event events.MarketValidatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal market validated event: %v", err)
		return
	}

	logger.Infof("Starting risk analysis for startup %s with market health: %s", event.StartupID, event.MarketHealth)

	// Create risk analysis request event
	riskAnalysisEvent := &events.RiskAnalysisRequestedEvent{
		BaseEvent: events.BaseEvent{
			ID:      fmt.Sprintf("risk-analysis-%s", event.StartupID),
			Type:    "risk.analysis.requested",
			Source:  "market-validation-service",
			Subject: events.SubjectRiskAnalysisRequested,
		},
		StartupID:      event.StartupID,
		StartupData:    map[string]interface{}{}, // Will be populated from startup data
		FounderCV:      map[string]interface{}{},
		BusinessPlan:   map[string]interface{}{},
		MarketData:     event.MarketData,
		SectorAnalysis: event.SectorAnalysis,
	}

	if err := h.eventService.PublishRiskAnalysisRequested(riskAnalysisEvent); err != nil {
		logger.Errorf("Failed to publish risk analysis requested event: %v", err)
		return
	}

	logger.Infof("Risk analysis requested for startup %s", event.StartupID)
}

// HandleRiskAnalysisRequested processes risk analysis requests
func (h *RiskAnalysisHandler) HandleRiskAnalysisRequested(msg *nats.Msg) {
	logger.Info("Processing risk analysis requested event")

	var event events.RiskAnalysisRequestedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal risk analysis requested event: %v", err)
		return
	}

	ctx := context.Background()

	// Get current startup info (this would typically come from the startup service)
	startupInfo := h.buildStartupAnalysisPrompt(event)

	// Generate comprehensive risk analysis using LLM
	logger.Infof("Generating comprehensive risk analysis for startup %s", event.StartupID)

	// Use the existing LLM service for initial risk analysis
	initialRisk, err := h.llmService.GenerateInitialRiskProfile(ctx, startupInfo)
	if err != nil {
		logger.Errorf("Failed to generate risk analysis: %v", err)
		return
	}

	// Create comprehensive risk analysis
	riskAnalysis := h.enhanceRiskAnalysisWithMarketData(initialRisk, event)

	// Create risk analysis completed event
	riskCompletedEvent := &events.RiskAnalysisCompletedEvent{
		BaseEvent: events.BaseEvent{
			ID:      fmt.Sprintf("risk-completed-%s", event.StartupID),
			Type:    "risk.analysis.completed",
			Source:  "risk-analysis-service",
			Subject: events.SubjectRiskAnalysisCompleted,
		},
		StartupID:        event.StartupID,
		RiskScore:        riskAnalysis.RiskScore,
		RiskLevel:        h.determineRiskLevel(riskAnalysis.RiskScore),
		Strengths:        riskAnalysis.Strengths,
		Weaknesses:       riskAnalysis.Weaknesses,
		Recommendations:  riskAnalysis.Recommendations,
		DetailedAnalysis: riskAnalysis.DetailedAnalysis,
	}

	// Publish risk analysis completed event
	if err := h.eventService.PublishRiskAnalysisCompleted(riskCompletedEvent); err != nil {
		logger.Errorf("Failed to publish risk analysis completed event: %v", err)
		return
	}

	logger.Infof("Risk analysis completed for startup %s with score: %.2f", event.StartupID, riskAnalysis.RiskScore)
}

// RiskAnalysisResult represents comprehensive risk analysis results
type RiskAnalysisResult struct {
	RiskScore        float64                `json:"risk_score"`
	Strengths        []string               `json:"strengths"`
	Weaknesses       []string               `json:"weaknesses"`
	Recommendations  []string               `json:"recommendations"`
	DetailedAnalysis map[string]interface{} `json:"detailed_analysis"`
}

func (h *RiskAnalysisHandler) buildStartupAnalysisPrompt(event events.RiskAnalysisRequestedEvent) string {
	// Build comprehensive analysis prompt including market data
	prompt := fmt.Sprintf(`
COMPREHENSIVE STARTUP RISK ANALYSIS REQUEST

Startup ID: %s

MARKET DATA ANALYSIS:
- Market Health: %v
- Sector Analysis: %v
- Market Data: %v

FOUNDER & BUSINESS INFORMATION:
- Founder CV: %v
- Business Plan: %v
- Startup Data: %v

Please provide a comprehensive risk assessment considering:
1. Market conditions and competitiveness
2. Founder experience and team composition
3. Business model viability
4. Financial projections and funding requirements
5. Regulatory and operational risks
6. Technology and execution risks

Focus on actionable insights and specific recommendations.
`,
		event.StartupID,
		event.MarketData,
		event.SectorAnalysis,
		event.MarketData,
		event.FounderCV,
		event.BusinessPlan,
		event.StartupData,
	)

	return prompt
}

func (h *RiskAnalysisHandler) enhanceRiskAnalysisWithMarketData(initialRisk *llm.InitialRiskResult, event events.RiskAnalysisRequestedEvent) *RiskAnalysisResult {
	// Extract market health from market data
	marketHealth := "fair" // default
	if marketData, ok := event.MarketData["market_health"]; ok {
		if health, ok := marketData.(string); ok {
			marketHealth = health
		}
	}

	// Adjust risk score based on market conditions
	adjustedScore := initialRisk.RiskScore
	switch marketHealth {
	case "excellent":
		adjustedScore *= 0.85 // Reduce risk by 15%
	case "good":
		adjustedScore *= 0.92 // Reduce risk by 8%
	case "declining":
		adjustedScore *= 1.15 // Increase risk by 15%
	case "poor":
		adjustedScore *= 1.25 // Increase risk by 25%
	}

	// Ensure score stays within bounds
	if adjustedScore > 100 {
		adjustedScore = 100
	}
	if adjustedScore < 0 {
		adjustedScore = 0
	}

	// Enhanced strengths and weaknesses
	strengths := initialRisk.Factors[:len(initialRisk.Factors)/2] // First half as strengths
	if len(strengths) == 0 {
		strengths = []string{"Market analysis completed", "Initial assessment positive"}
	}

	weaknesses := initialRisk.Factors[len(initialRisk.Factors)/2:] // Second half as weaknesses
	if len(weaknesses) == 0 {
		weaknesses = []string{"Market conditions uncertain", "Need more validation"}
	}

	// Enhanced recommendations
	recommendations := initialRisk.Suggestions
	if marketHealth == "declining" || marketHealth == "poor" {
		recommendations = append(recommendations,
			"Consider pivoting to a more stable market segment",
			"Develop stronger competitive differentiation",
			"Build strategic partnerships to mitigate market risks")
	}

	return &RiskAnalysisResult{
		RiskScore:       adjustedScore,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
		DetailedAnalysis: map[string]interface{}{
			"market_health":   marketHealth,
			"initial_score":   initialRisk.RiskScore,
			"adjusted_score":  adjustedScore,
			"market_factor":   fmt.Sprintf("Market conditions: %s", marketHealth),
			"analysis_source": "AI-powered comprehensive assessment",
			"confidence":      0.85,
			"reasoning":       initialRisk.Reasoning,
		},
	}
}

func (h *RiskAnalysisHandler) determineRiskLevel(score float64) string {
	switch {
	case score >= 80:
		return "very_high"
	case score >= 60:
		return "high"
	case score >= 40:
		return "medium"
	case score >= 20:
		return "low"
	default:
		return "very_low"
	}
}
