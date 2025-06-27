package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"risq_backend/internal/contextmem"
	"risq_backend/pkg/events"
	"risq_backend/pkg/logger"

	"github.com/nats-io/nats.go"
)

// ContextStorageHandler handles context storage events for RAG
type ContextStorageHandler struct {
	contextMemService contextmem.Service
	eventService      events.EventService
}

func NewContextStorageHandler(contextMemService contextmem.Service, eventService events.EventService) *ContextStorageHandler {
	return &ContextStorageHandler{
		contextMemService: contextMemService,
		eventService:      eventService,
	}
}

// HandleRiskAnalysisCompleted stores risk analysis results in RAG context
func (h *ContextStorageHandler) HandleRiskAnalysisCompleted(msg *nats.Msg) {
	logger.Info("Processing risk analysis completed event for context storage")

	var event events.RiskAnalysisCompletedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal risk analysis completed event: %v", err)
		return
	}

	ctx := context.Background()

	// Build comprehensive context content
	contextContent := h.buildContextContent(event)

	// Create metadata for the context
	metadata := map[string]interface{}{
		"type":               "risk_analysis",
		"startup_id":         event.StartupID.String(),
		"risk_score":         event.RiskScore,
		"risk_level":         event.RiskLevel,
		"analysis_timestamp": event.BaseEvent.Timestamp,
		"source":             "comprehensive_risk_analysis",
		"version":            "1.0",
	}

	// Store the context
	logger.Infof("Storing risk analysis context for startup %s", event.StartupID)

	if err := h.contextMemService.StoreStartupContext(ctx, event.StartupID, contextContent, metadata); err != nil {
		logger.Errorf("Failed to store risk analysis context: %v", err)
		return
	}

	// Create context stored event (for potential future use)
	// contextStoredEvent := &events.ContextStoreRequestedEvent{
	// 	BaseEvent: events.BaseEvent{
	// 		ID:        fmt.Sprintf("context-stored-%s", event.StartupID),
	// 		Type:      "context.stored",
	// 		Source:    "context-storage-service",
	// 		Subject:   events.SubjectContextStored,
	// 	},
	// 	StartupID:   event.StartupID,
	// 	ContentType: "risk_analysis",
	// 	Content:     contextContent,
	// 	Metadata:    metadata,
	// }

	// Note: We could publish this event if we need downstream processing
	// For now, we'll just log the completion
	logger.Infof("Risk analysis context stored successfully for startup %s - Risk Score: %.2f, Level: %s",
		event.StartupID, event.RiskScore, event.RiskLevel)

	// Optional: Store additional context fragments for different aspects
	h.storeDetailedContextFragments(ctx, event)
}

// HandleContextStoreRequested processes explicit context storage requests
func (h *ContextStorageHandler) HandleContextStoreRequested(msg *nats.Msg) {
	logger.Info("Processing context store requested event")

	var event events.ContextStoreRequestedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal context store requested event: %v", err)
		return
	}

	ctx := context.Background()

	logger.Infof("Storing context for startup %s - Type: %s", event.StartupID, event.ContentType)

	if err := h.contextMemService.StoreStartupContext(ctx, event.StartupID, event.Content, event.Metadata); err != nil {
		logger.Errorf("Failed to store context: %v", err)
		return
	}

	logger.Infof("Context stored successfully for startup %s", event.StartupID)
}

func (h *ContextStorageHandler) buildContextContent(event events.RiskAnalysisCompletedEvent) string {
	// Build comprehensive context content that will be useful for future RAG queries
	content := fmt.Sprintf(`
COMPREHENSIVE RISK ANALYSIS RESULTS

Startup ID: %s
Overall Risk Score: %.2f/100
Risk Level: %s

STRENGTHS IDENTIFIED:
%s

WEAKNESSES IDENTIFIED:
%s

KEY RECOMMENDATIONS:
%s

DETAILED ANALYSIS INSIGHTS:
%s

This analysis was completed on %s and represents a comprehensive assessment 
of market conditions, founder capabilities, business model viability, 
and operational risks. The risk score of %.2f indicates a %s risk level 
requiring specific attention to the identified weaknesses and implementation 
of the provided recommendations.

CONTEXT FOR FUTURE DECISIONS:
- Current risk profile established at score %.2f
- Market conditions assessed and factored into analysis
- Founder experience and team composition evaluated
- Business model strengths and vulnerabilities identified
- Specific mitigation strategies provided

This context should be referenced for future business decisions, 
funding discussions, and strategic planning initiatives.
`,
		event.StartupID,
		event.RiskScore,
		event.RiskLevel,
		h.formatListItems(event.Strengths),
		h.formatListItems(event.Weaknesses),
		h.formatListItems(event.Recommendations),
		h.formatDetailedAnalysis(event.DetailedAnalysis),
		event.BaseEvent.Timestamp.Format("2006-01-02 15:04:05"),
		event.RiskScore,
		event.RiskLevel,
		event.RiskScore,
	)

	return content
}

func (h *ContextStorageHandler) formatListItems(items []string) string {
	if len(items) == 0 {
		return "None specified"
	}

	result := ""
	for i, item := range items {
		result += fmt.Sprintf("  %d. %s\n", i+1, item)
	}
	return result
}

func (h *ContextStorageHandler) formatDetailedAnalysis(analysis map[string]interface{}) string {
	if len(analysis) == 0 {
		return "No detailed analysis available"
	}

	result := ""
	for key, value := range analysis {
		result += fmt.Sprintf("  - %s: %v\n", key, value)
	}
	return result
}

func (h *ContextStorageHandler) storeDetailedContextFragments(ctx context.Context, event events.RiskAnalysisCompletedEvent) {
	// Store specific context fragments for different aspects

	// 1. Store strengths as a separate context
	if len(event.Strengths) > 0 {
		strengthsContent := fmt.Sprintf("STARTUP STRENGTHS (Risk Score: %.2f):\n%s",
			event.RiskScore, h.formatListItems(event.Strengths))

		strengthsMetadata := map[string]interface{}{
			"type":       "strengths",
			"startup_id": event.StartupID.String(),
			"risk_score": event.RiskScore,
			"fragment":   "strengths_analysis",
			"timestamp":  event.BaseEvent.Timestamp,
		}

		if err := h.contextMemService.StoreStartupContext(ctx, event.StartupID, strengthsContent, strengthsMetadata); err != nil {
			logger.Warnf("Failed to store strengths context: %v", err)
		}
	}

	// 2. Store weaknesses as a separate context
	if len(event.Weaknesses) > 0 {
		weaknessesContent := fmt.Sprintf("STARTUP WEAKNESSES (Risk Score: %.2f):\n%s",
			event.RiskScore, h.formatListItems(event.Weaknesses))

		weaknessesMetadata := map[string]interface{}{
			"type":       "weaknesses",
			"startup_id": event.StartupID.String(),
			"risk_score": event.RiskScore,
			"fragment":   "weaknesses_analysis",
			"timestamp":  event.BaseEvent.Timestamp,
		}

		if err := h.contextMemService.StoreStartupContext(ctx, event.StartupID, weaknessesContent, weaknessesMetadata); err != nil {
			logger.Warnf("Failed to store weaknesses context: %v", err)
		}
	}

	// 3. Store recommendations as actionable context
	if len(event.Recommendations) > 0 {
		recommendationsContent := fmt.Sprintf("ACTIONABLE RECOMMENDATIONS (Risk Score: %.2f):\n%s",
			event.RiskScore, h.formatListItems(event.Recommendations))

		recommendationsMetadata := map[string]interface{}{
			"type":       "recommendations",
			"startup_id": event.StartupID.String(),
			"risk_score": event.RiskScore,
			"fragment":   "recommendations",
			"actionable": true,
			"timestamp":  event.BaseEvent.Timestamp,
		}

		if err := h.contextMemService.StoreStartupContext(ctx, event.StartupID, recommendationsContent, recommendationsMetadata); err != nil {
			logger.Warnf("Failed to store recommendations context: %v", err)
		}
	}

	logger.Infof("Stored %d detailed context fragments for startup %s",
		3, event.StartupID) // strengths, weaknesses, recommendations
}
