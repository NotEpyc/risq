package app

import (
	"context"
	"fmt"

	"risq_backend/config"
	"risq_backend/internal/contextmem"
	"risq_backend/internal/handlers"
	"risq_backend/internal/llm"
	"risq_backend/internal/risk"
	"risq_backend/pkg/events"
	"risq_backend/pkg/external"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// EventManager manages all event subscriptions and handlers
type EventManager struct {
	eventService            events.EventService
	marketValidationHandler *handlers.MarketValidationHandler
	riskAnalysisHandler     *handlers.RiskAnalysisHandler
	contextStorageHandler   *handlers.ContextStorageHandler
	subscriptions           []*nats.Subscription
}

// NewEventManager creates a new event manager with all handlers
func NewEventManager(
	cfg *config.Config,
	riskService risk.Service,
	llmService llm.Service,
	contextMemService contextmem.Service,
	marketDataService external.MarketDataService,
) *EventManager {
	// Create NATS event service
	eventService := events.NewNATSEventService(cfg.NATS.URL)

	// Create handlers
	marketValidationHandler := handlers.NewMarketValidationHandler(marketDataService, eventService)
	riskAnalysisHandler := handlers.NewRiskAnalysisHandler(riskService, llmService, eventService)
	contextStorageHandler := handlers.NewContextStorageHandler(contextMemService, eventService)

	return &EventManager{
		eventService:            eventService,
		marketValidationHandler: marketValidationHandler,
		riskAnalysisHandler:     riskAnalysisHandler,
		contextStorageHandler:   contextStorageHandler,
		subscriptions:           make([]*nats.Subscription, 0),
	}
}

// Start connects to NATS and sets up all event subscriptions
func (em *EventManager) Start(ctx context.Context) error {
	logger.Info("Starting Event Manager...")

	// Connect to NATS
	if err := em.eventService.Connect(); err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Setup all event subscriptions
	if err := em.setupSubscriptions(); err != nil {
		return fmt.Errorf("failed to setup subscriptions: %w", err)
	}

	logger.Info("Event Manager started successfully with all subscriptions active")
	return nil
}

// Stop closes all subscriptions and connections
func (em *EventManager) Stop() error {
	logger.Info("Stopping Event Manager...")

	// Unsubscribe from all subscriptions
	for _, sub := range em.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			logger.Warnf("Failed to unsubscribe: %v", err)
		}
	}

	// Close event service connection
	if err := em.eventService.Close(); err != nil {
		logger.Warnf("Failed to close event service: %v", err)
	}

	logger.Info("Event Manager stopped")
	return nil
}

// GetEventService returns the event service for publishing events
func (em *EventManager) GetEventService() events.EventService {
	return em.eventService
}

func (em *EventManager) setupSubscriptions() error {
	logger.Info("Setting up event subscriptions...")

	// 1. Startup Onboarded -> Market Validation
	logger.Info("Subscribing to startup onboarded events...")
	sub1, err := em.eventService.Subscribe(events.SubjectStartupOnboarding, em.marketValidationHandler.HandleStartupOnboarded)
	if err != nil {
		return fmt.Errorf("failed to subscribe to startup onboarded events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub1)

	// 2. Market Validation Requested -> Market Analysis
	logger.Info("Subscribing to market validation requested events...")
	sub2, err := em.eventService.Subscribe(events.SubjectMarketValidation, em.marketValidationHandler.HandleMarketValidationRequested)
	if err != nil {
		return fmt.Errorf("failed to subscribe to market validation requested events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub2)

	// 3. Market Validated -> Risk Analysis
	logger.Info("Subscribing to market validated events...")
	sub3, err := em.eventService.Subscribe(events.SubjectMarketValidated, em.riskAnalysisHandler.HandleMarketValidated)
	if err != nil {
		return fmt.Errorf("failed to subscribe to market validated events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub3)

	// 4. Risk Analysis Requested -> Risk Processing
	logger.Info("Subscribing to risk analysis requested events...")
	sub4, err := em.eventService.Subscribe(events.SubjectRiskAnalysisRequested, em.riskAnalysisHandler.HandleRiskAnalysisRequested)
	if err != nil {
		return fmt.Errorf("failed to subscribe to risk analysis requested events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub4)

	// 5. Risk Analysis Completed -> Context Storage
	logger.Info("Subscribing to risk analysis completed events...")
	sub5, err := em.eventService.Subscribe(events.SubjectRiskAnalysisCompleted, em.contextStorageHandler.HandleRiskAnalysisCompleted)
	if err != nil {
		return fmt.Errorf("failed to subscribe to risk analysis completed events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub5)

	// 6. Context Store Requested -> Context Storage
	logger.Info("Subscribing to context store requested events...")
	sub6, err := em.eventService.Subscribe(events.SubjectContextStoreRequested, em.contextStorageHandler.HandleContextStoreRequested)
	if err != nil {
		return fmt.Errorf("failed to subscribe to context store requested events: %w", err)
	}
	em.subscriptions = append(em.subscriptions, sub6)

	logger.Infof("Successfully set up %d event subscriptions", len(em.subscriptions))
	logger.Info("Event flow: Startup Onboarding → Market Validation → Risk Analysis → Context Storage")

	return nil
}

// PublishStartupOnboarded publishes a startup onboarded event (used by controllers)
func (em *EventManager) PublishStartupOnboarded(startupID, userID uuid.UUID, startupData, founderCV, businessPlan map[string]interface{}) error {
	event := events.NewStartupOnboardedEvent(startupID, userID, startupData, founderCV, businessPlan)
	return em.eventService.PublishStartupOnboarded(event)
}
