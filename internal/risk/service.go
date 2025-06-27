package risk

import (
	"context"

	"risq_backend/internal/llm"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type Service interface {
	GetCurrentRisk(ctx context.Context, startupID uuid.UUID) (*RiskProfile, error)
	CreateInitialProfile(ctx context.Context, startupID uuid.UUID, startupInfo string) (*RiskProfile, error)
	UpdateRiskFromDecision(ctx context.Context, decision interface{}) error
	GetEvolutionHistory(ctx context.Context, startupID uuid.UUID, limit int) ([]*RiskEvolution, error)
}

type service struct {
	repo       Repository
	llmService llm.Service
}

func NewService(repo Repository, llmService llm.Service) Service {
	return &service{
		repo:       repo,
		llmService: llmService,
	}
}

func (s *service) GetCurrentRisk(ctx context.Context, startupID uuid.UUID) (*RiskProfile, error) {
	logger.Debugf("Getting current risk for startup: %s", startupID)

	profile, err := s.repo.GetCurrentProfile(ctx, startupID)
	if err != nil {
		logger.Errorf("Failed to get current risk profile: %v", err)
		return nil, err
	}

	return profile, nil
}

func (s *service) CreateInitialProfile(ctx context.Context, startupID uuid.UUID, startupInfo string) (*RiskProfile, error) {
	logger.Infof("Creating initial risk profile for startup: %s", startupID)

	// Generate initial risk assessment using LLM
	result, err := s.llmService.GenerateInitialRiskProfile(ctx, startupInfo)
	if err != nil {
		logger.Errorf("Failed to generate initial risk profile: %v", err)
		return nil, err
	}

	// Create risk profile
	profile := &RiskProfile{
		ID:          uuid.New(),
		StartupID:   startupID,
		Score:       result.RiskScore,
		Level:       DetermineRiskLevel(result.RiskScore),
		Confidence:  0.8, // Initial confidence
		Factors:     result.Factors,
		Suggestions: result.Suggestions,
		Reasoning:   result.Reasoning,
	}

	if err := s.repo.CreateProfile(ctx, profile); err != nil {
		logger.Errorf("Failed to create risk profile: %v", err)
		return nil, err
	}

	// Create evolution record
	evolution := &RiskEvolution{
		ID:        uuid.New(),
		StartupID: startupID,
		Score:     result.RiskScore,
		Level:     profile.Level,
		Trigger:   "Initial profile creation",
	}

	if err := s.repo.CreateEvolution(ctx, evolution); err != nil {
		logger.Warnf("Failed to create risk evolution: %v", err)
	}

	logger.Infof("Successfully created initial risk profile with score: %.2f", result.RiskScore)
	return profile, nil
}

func (s *service) UpdateRiskFromDecision(ctx context.Context, decision interface{}) error {
	logger.Infof("Updating risk from decision")

	// This would need proper type assertion and logic
	// For now, we'll implement a basic version
	logger.Warnf("UpdateRiskFromDecision not fully implemented yet")
	return nil
}

func (s *service) GetEvolutionHistory(ctx context.Context, startupID uuid.UUID, limit int) ([]*RiskEvolution, error) {
	logger.Debugf("Getting risk evolution history for startup: %s", startupID)

	history, err := s.repo.GetEvolutionHistory(ctx, startupID, limit)
	if err != nil {
		logger.Errorf("Failed to get risk evolution history: %v", err)
		return nil, err
	}

	return history, nil
}
