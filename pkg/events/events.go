package events

import (
	"encoding/json"
	"fmt"
	"time"

	"risq_backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// Event subjects (topics)
const (
	SubjectStartupOnboarding     = "startup.onboarded"
	SubjectMarketValidation      = "market.validation.requested"
	SubjectMarketValidated       = "market.validated"
	SubjectRiskAnalysisRequested = "risk.analysis.requested"
	SubjectRiskAnalysisCompleted = "risk.analysis.completed"
	SubjectContextStoreRequested = "context.store.requested"
	SubjectContextStored         = "context.stored"
)

// Event types are now defined in types.go - using those comprehensive definitions

// EventService interface
type EventService interface {
	Connect() error
	Close() error
	PublishStartupOnboarded(event *StartupOnboardedEvent) error
	PublishMarketValidationRequested(event *MarketValidationRequestedEvent) error
	PublishMarketValidated(event *MarketValidatedEvent) error
	PublishRiskAnalysisRequested(event *RiskAnalysisRequestedEvent) error
	PublishRiskAnalysisCompleted(event *RiskAnalysisCompletedEvent) error
	PublishContextStoreRequested(event *ContextStoreRequestedEvent) error
	Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error)
}

// NATS implementation
type natsEventService struct {
	conn *nats.Conn
	url  string
}

func NewNATSEventService(url string) EventService {
	return &natsEventService{
		url: url,
	}
}

func (s *natsEventService) Connect() error {
	logger.Info("Connecting to NATS...")

	conn, err := nats.Connect(s.url,
		nats.ReconnectWait(2*time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Errorf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	s.conn = conn
	logger.Infof("Connected to NATS at %s", s.url)
	return nil
}

func (s *natsEventService) Close() error {
	if s.conn != nil {
		s.conn.Close()
		logger.Info("NATS connection closed")
	}
	return nil
}

func (s *natsEventService) publish(subject string, event interface{}) error {
	if s.conn == nil {
		return fmt.Errorf("NATS connection not established")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := s.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", subject, err)
	}

	logger.Infof("Published event to subject: %s", subject)
	return nil
}

func (s *natsEventService) PublishStartupOnboarded(event *StartupOnboardedEvent) error {
	return s.publish(SubjectStartupOnboarding, event)
}

func (s *natsEventService) PublishMarketValidationRequested(event *MarketValidationRequestedEvent) error {
	return s.publish(SubjectMarketValidation, event)
}

func (s *natsEventService) PublishMarketValidated(event *MarketValidatedEvent) error {
	return s.publish(SubjectMarketValidated, event)
}

func (s *natsEventService) PublishRiskAnalysisRequested(event *RiskAnalysisRequestedEvent) error {
	return s.publish(SubjectRiskAnalysisRequested, event)
}

func (s *natsEventService) PublishRiskAnalysisCompleted(event *RiskAnalysisCompletedEvent) error {
	return s.publish(SubjectRiskAnalysisCompleted, event)
}

func (s *natsEventService) PublishContextStoreRequested(event *ContextStoreRequestedEvent) error {
	return s.publish(SubjectContextStoreRequested, event)
}

func (s *natsEventService) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("NATS connection not established")
	}

	sub, err := s.conn.Subscribe(subject, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}

	logger.Infof("Subscribed to subject: %s", subject)
	return sub, nil
}

// Helper functions to create events
func NewStartupOnboardedEvent(startupID, userID uuid.UUID, startupData, founderCV, businessPlan map[string]interface{}) *StartupOnboardedEvent {
	return &StartupOnboardedEvent{
		BaseEvent: BaseEvent{
			ID:        uuid.New().String(),
			Type:      "startup.onboarded",
			Source:    "startup-service",
			Subject:   SubjectStartupOnboarding,
			Timestamp: time.Now(),
		},
		StartupID:    startupID,
		UserID:       userID,
		StartupData:  startupData,
		FounderCV:    founderCV,
		BusinessPlan: businessPlan,
	}
}

func NewMarketValidationRequestedEvent(startupID uuid.UUID, industry, sector, targetMarket, businessModel string) *MarketValidationRequestedEvent {
	return &MarketValidationRequestedEvent{
		BaseEvent: BaseEvent{
			ID:        uuid.New().String(),
			Type:      "market.validation.requested",
			Source:    "startup-service",
			Subject:   SubjectMarketValidation,
			Timestamp: time.Now(),
		},
		StartupID:     startupID,
		Industry:      industry,
		Sector:        sector,
		TargetMarket:  targetMarket,
		BusinessModel: businessModel,
	}
}

func NewRiskAnalysisRequestedEvent(startupID uuid.UUID, startupData, founderCV, businessPlan, marketData, sectorAnalysis map[string]interface{}) *RiskAnalysisRequestedEvent {
	return &RiskAnalysisRequestedEvent{
		BaseEvent: BaseEvent{
			ID:        uuid.New().String(),
			Type:      "risk.analysis.requested",
			Source:    "market-service",
			Subject:   SubjectRiskAnalysisRequested,
			Timestamp: time.Now(),
		},
		StartupID:      startupID,
		StartupData:    startupData,
		FounderCV:      founderCV,
		BusinessPlan:   businessPlan,
		MarketData:     marketData,
		SectorAnalysis: sectorAnalysis,
	}
}
