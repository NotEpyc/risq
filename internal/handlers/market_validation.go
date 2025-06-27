package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"risq_backend/pkg/events"
	"risq_backend/pkg/external"
	"risq_backend/pkg/logger"

	"github.com/nats-io/nats.go"
)

// MarketValidationHandler handles startup onboarding events and triggers market validation
type MarketValidationHandler struct {
	marketDataService external.MarketDataService
	eventService      events.EventService
}

func NewMarketValidationHandler(marketDataService external.MarketDataService, eventService events.EventService) *MarketValidationHandler {
	return &MarketValidationHandler{
		marketDataService: marketDataService,
		eventService:      eventService,
	}
}

func (h *MarketValidationHandler) HandleStartupOnboarded(msg *nats.Msg) {
	logger.Info("Processing startup onboarded event for market validation")

	var event events.StartupOnboardedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal startup onboarded event: %v", err)
		return
	}

	// Extract startup details from the event
	industry, _ := event.StartupData["industry"].(string)
	sector, _ := event.StartupData["sector"].(string)
	targetMarket, _ := event.StartupData["target_market"].(string)
	businessModel, _ := event.StartupData["business_model"].(string)

	logger.Infof("Validating market for startup %s in %s sector", event.StartupID, sector)

	// Trigger market validation request
	marketValidationEvent := events.NewMarketValidationRequestedEvent(
		event.StartupID,
		industry,
		sector,
		targetMarket,
		businessModel,
	)

	if err := h.eventService.PublishMarketValidationRequested(marketValidationEvent); err != nil {
		logger.Errorf("Failed to publish market validation requested event: %v", err)
		return
	}

	logger.Infof("Market validation requested for startup %s", event.StartupID)
}

func (h *MarketValidationHandler) HandleMarketValidationRequested(msg *nats.Msg) {
	logger.Info("Processing market validation requested event")

	var event events.MarketValidationRequestedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal market validation requested event: %v", err)
		return
	}

	ctx := context.Background()

	// Get industry trends
	industryTrends, err := h.marketDataService.GetIndustryTrends(ctx, event.Industry)
	if err != nil {
		logger.Errorf("Failed to get industry trends: %v", err)
		industryTrends = &external.IndustryTrends{
			Industry:         event.Industry,
			GrowthRate:       10.0,
			CompetitionLevel: "moderate",
			Outlook:          "stable",
		}
	}

	// Get sector analysis
	sectorAnalysis, err := h.marketDataService.GetSectorAnalysis(ctx, event.Sector)
	if err != nil {
		logger.Errorf("Failed to get sector analysis: %v", err)
		sectorAnalysis = &external.SectorAnalysis{
			Sector:   event.Sector,
			IsActive: true,
			Activity: "stable",
		}
	}

	// Get market health
	marketHealth, err := h.marketDataService.GetMarketHealth(ctx, event.TargetMarket)
	if err != nil {
		logger.Errorf("Failed to get market health: %v", err)
		marketHealth = &external.MarketHealth{
			TargetMarket: event.TargetMarket,
			Health:       "fair",
		}
	}

	// Determine overall market health
	overallHealth := h.determineOverallMarketHealth(industryTrends, sectorAnalysis, marketHealth)

	// Create market validated event
	marketData := map[string]interface{}{
		"industry_trends": industryTrends,
		"market_health":   marketHealth,
	}

	sectorAnalysisData := map[string]interface{}{
		"sector_analysis": sectorAnalysis,
		"is_active":       sectorAnalysis.IsActive,
		"activity":        sectorAnalysis.Activity,
		"investment_flow": sectorAnalysis.InvestmentFlow,
	}

	recommendations := h.generateMarketRecommendations(industryTrends, sectorAnalysis, marketHealth)

	marketValidatedEvent := &events.MarketValidatedEvent{
		BaseEvent: events.BaseEvent{
			ID:      fmt.Sprintf("market-validated-%s", event.StartupID),
			Type:    "market.validated",
			Source:  "market-validation-service",
			Subject: events.SubjectMarketValidated,
		},
		StartupID:       event.StartupID,
		MarketData:      marketData,
		SectorAnalysis:  sectorAnalysisData,
		MarketHealth:    overallHealth,
		Recommendations: recommendations,
	}

	if err := h.eventService.PublishMarketValidated(marketValidatedEvent); err != nil {
		logger.Errorf("Failed to publish market validated event: %v", err)
		return
	}

	logger.Infof("Market validation completed for startup %s: Health=%s", event.StartupID, overallHealth)
}

func (h *MarketValidationHandler) determineOverallMarketHealth(
	industryTrends *external.IndustryTrends,
	sectorAnalysis *external.SectorAnalysis,
	marketHealth *external.MarketHealth,
) string {
	score := 0

	// Industry trends contribution (40%)
	if industryTrends.GrowthRate > 15 {
		score += 40
	} else if industryTrends.GrowthRate > 10 {
		score += 30
	} else if industryTrends.GrowthRate > 5 {
		score += 20
	} else {
		score += 10
	}

	// Sector activity contribution (35%)
	if sectorAnalysis.IsActive {
		switch sectorAnalysis.Activity {
		case "growing":
			score += 35
		case "stable":
			score += 25
		default:
			score += 10
		}
	} else {
		score += 5
	}

	// Market health contribution (25%)
	switch marketHealth.Health {
	case "excellent":
		score += 25
	case "good":
		score += 20
	case "fair":
		score += 15
	default:
		score += 10
	}

	// Determine overall health
	if score >= 80 {
		return "active"
	} else if score >= 60 {
		return "moderately_active"
	} else if score >= 40 {
		return "cautious"
	} else {
		return "inactive"
	}
}

func (h *MarketValidationHandler) generateMarketRecommendations(
	industryTrends *external.IndustryTrends,
	sectorAnalysis *external.SectorAnalysis,
	marketHealth *external.MarketHealth,
) []string {
	var recommendations []string

	// Industry-based recommendations
	if industryTrends.GrowthRate > 15 {
		recommendations = append(recommendations, "Leverage high industry growth rate for rapid scaling")
	}

	if industryTrends.CompetitionLevel == "high" {
		recommendations = append(recommendations, "Focus on differentiation due to high competition")
	}

	// Sector-based recommendations
	if sectorAnalysis.IsActive && sectorAnalysis.Activity == "growing" {
		recommendations = append(recommendations, "Take advantage of growing sector trends")
	}

	if sectorAnalysis.InvestmentFlow > 10 {
		recommendations = append(recommendations, "Strong investor interest - good timing for fundraising")
	}

	// Market health recommendations
	if marketHealth.SaturationLevel > 70 {
		recommendations = append(recommendations, "Market is saturated - consider niche targeting")
	}

	// Add general recommendations
	recommendations = append(recommendations, marketHealth.Recommendations...)

	return recommendations
}
