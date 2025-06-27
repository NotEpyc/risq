package decision

import (
	"context"
	"fmt"
	"time"

	"risq_backend/internal/contextmem"
	"risq_backend/internal/llm"
	"risq_backend/internal/risk"
	"risq_backend/pkg/logger"
	"risq_backend/types"

	"github.com/google/uuid"
)

type Service interface {
	Speculate(ctx context.Context, input types.DecisionInput) (*DecisionResponse, error)
	Confirm(ctx context.Context, input DecisionConfirmInput) (*Decision, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Decision, error)
	GetByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error)
}

type service struct {
	repo        Repository
	llmService  llm.Service
	contextMem  contextmem.Service
	riskService risk.Service
}

func NewService(repo Repository, llmService llm.Service, contextMem contextmem.Service, riskService risk.Service) Service {
	return &service{
		repo:        repo,
		llmService:  llmService,
		contextMem:  contextMem,
		riskService: riskService,
	}
}

func (s *service) Speculate(ctx context.Context, input types.DecisionInput) (*DecisionResponse, error) {
	logger.Infof("Starting speculation for decision: %s", input.Description)

	// Get current risk score
	currentRisk, err := s.riskService.GetCurrentRisk(ctx, input.StartupID)
	if err != nil {
		logger.Errorf("Failed to get current risk: %v", err)
		return nil, fmt.Errorf("failed to get current risk: %w", err)
	}

	// Fetch relevant context from vector memory
	contextChunks, err := s.contextMem.FetchRelevantContext(ctx, input.StartupID, input.Description)
	if err != nil {
		logger.Warnf("Failed to fetch context: %v", err)
		contextChunks = []string{} // Continue without context
	}

	// Generate LLM-based speculation
	speculation, err := s.llmService.SpeculateDecision(ctx, input, contextChunks, currentRisk.Score)
	if err != nil {
		logger.Errorf("Failed to generate speculation: %v", err)
		return nil, fmt.Errorf("failed to generate speculation: %w", err)
	}

	// Create decision record
	decision := &Decision{
		ID:                 uuid.New(),
		StartupID:          input.StartupID,
		Description:        input.Description,
		Category:           input.Category,
		Context:            input.Context,
		Timeline:           input.Timeline,
		Budget:             input.Budget,
		Status:             types.DecisionStatusSpeculative,
		PreviousRiskScore:  currentRisk.Score,
		ProjectedRiskScore: speculation.ProjectedRiskScore,
		RiskDelta:          speculation.ProjectedRiskScore - currentRisk.Score,
		Confidence:         speculation.Confidence,
		Suggestions:        speculation.Suggestions,
		Reasoning:          speculation.Reasoning,
	}

	if err := s.repo.Create(ctx, decision); err != nil {
		logger.Errorf("Failed to create decision: %v", err)
		return nil, err
	}

	response := &DecisionResponse{
		Decision:    *decision,
		RiskScore:   speculation.ProjectedRiskScore,
		RiskDelta:   decision.RiskDelta,
		Suggestions: speculation.Suggestions,
		Reasoning:   speculation.Reasoning,
		Confidence:  speculation.Confidence,
	}

	logger.Infof("Successfully created speculation: %s", decision.ID)
	return response, nil
}

func (s *service) Confirm(ctx context.Context, input DecisionConfirmInput) (*Decision, error) {
	logger.Infof("Confirming decision: %s", input.DecisionID)

	// Get the speculative decision
	decision, err := s.repo.GetByID(ctx, input.DecisionID)
	if err != nil {
		return nil, err
	}

	if decision.Status != types.DecisionStatusSpeculative {
		return nil, fmt.Errorf("decision is not in speculative state")
	}

	// Update decision status
	now := time.Now()
	decision.Status = types.DecisionStatusConfirmed
	decision.ConfirmedAt = &now

	if err := s.repo.Update(ctx, decision); err != nil {
		logger.Errorf("Failed to update decision: %v", err)
		return nil, err
	}

	// Update context memory with confirmed decision
	if err := s.contextMem.StoreDecision(ctx, decision); err != nil {
		logger.Warnf("Failed to store decision in context memory: %v", err)
	}

	// Update risk score
	if err := s.riskService.UpdateRiskFromDecision(ctx, decision); err != nil {
		logger.Warnf("Failed to update risk score: %v", err)
	}

	logger.Infof("Successfully confirmed decision: %s", decision.ID)
	return decision, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Decision, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByStartupID(ctx context.Context, startupID uuid.UUID) ([]*Decision, error) {
	return s.repo.GetByStartupID(ctx, startupID)
}
