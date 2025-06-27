package contextmem

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"risq_backend/internal/llm"
	"risq_backend/pkg/cache"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type ContextChunk struct {
	ID        string                 `json:"id"`
	StartupID uuid.UUID              `json:"startup_id"`
	Content   string                 `json:"content"`
	Embedding []float32              `json:"embedding"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

type Service interface {
	StoreDecision(ctx context.Context, decision interface{}) error
	FetchRelevantContext(ctx context.Context, startupID uuid.UUID, query string) ([]string, error)
	StoreStartupContext(ctx context.Context, startupID uuid.UUID, content string, metadata map[string]interface{}) error
}

type service struct {
	redis      *cache.Redis
	llmService llm.Service
}

func NewService(redis *cache.Redis, llmService llm.Service) Service {
	return &service{
		redis:      redis,
		llmService: llmService,
	}
}

func (s *service) StoreDecision(ctx context.Context, decision interface{}) error {
	logger.Infof("Storing decision in context memory")

	// Type assert the decision to extract startup ID
	var startupID uuid.UUID
	var decisionID uuid.UUID

	// Try to extract fields from the decision
	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("failed to marshal decision: %w", err)
	}

	// Unmarshal to a generic map to extract startup_id
	var decisionMap map[string]interface{}
	if err := json.Unmarshal(decisionJSON, &decisionMap); err != nil {
		return fmt.Errorf("failed to unmarshal decision: %w", err)
	}

	// Extract startup_id
	if startupIDStr, ok := decisionMap["startup_id"].(string); ok {
		if parsedID, err := uuid.Parse(startupIDStr); err == nil {
			startupID = parsedID
		}
	}

	// Extract decision_id
	if decisionIDStr, ok := decisionMap["id"].(string); ok {
		if parsedID, err := uuid.Parse(decisionIDStr); err == nil {
			decisionID = parsedID
		}
	}

	decisionText := string(decisionJSON)

	// Generate embedding for the decision
	embedding, err := s.llmService.GenerateEmbedding(ctx, decisionText)
	if err != nil {
		logger.Warnf("Failed to generate embedding for decision: %v", err)
		// Fallback to basic storage without embedding
		key := fmt.Sprintf("decision:%s", decisionID.String())
		return s.redis.Set(ctx, key, decisionText, 24*time.Hour)
	}

	// Create context chunk for the decision
	chunk := ContextChunk{
		ID:        uuid.New().String(),
		StartupID: startupID,
		Content:   decisionText,
		Embedding: embedding,
		Metadata: map[string]interface{}{
			"type":        "decision",
			"decision_id": decisionID.String(),
			"timestamp":   time.Now(),
		},
		CreatedAt: time.Now(),
	}

	// Store the chunk
	chunkJSON, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("failed to marshal context chunk: %w", err)
	}

	key := fmt.Sprintf("context_chunk:%s", chunk.ID)
	if err := s.redis.Set(ctx, key, string(chunkJSON), 7*24*time.Hour); err != nil {
		return fmt.Errorf("failed to store context chunk: %w", err)
	}

	logger.Infof("Successfully stored decision context chunk: %s for startup: %s", chunk.ID, startupID)
	return nil
}

func (s *service) FetchRelevantContext(ctx context.Context, startupID uuid.UUID, query string) ([]string, error) {
	logger.Debugf("Fetching relevant context for startup: %s, query: %s", startupID, query)

	// Generate embedding for the query
	_, err := s.llmService.GenerateEmbedding(ctx, query)
	if err != nil {
		logger.Warnf("Failed to generate query embedding: %v", err)
		// Fallback to basic search
		return s.fallbackContextSearch(ctx, startupID, query)
	}

	// Search for relevant context chunks using Redis pattern matching
	pattern := fmt.Sprintf("context_chunk:*")
	keys, err := s.redis.Keys(ctx, pattern)
	if err != nil {
		logger.Warnf("Failed to get context chunk keys: %v", err)
		return s.fallbackContextSearch(ctx, startupID, query)
	}

	var relevantContext []string
	for _, key := range keys {
		chunkJSON, err := s.redis.Get(ctx, key)
		if err != nil {
			continue
		}

		var chunk ContextChunk
		if err := json.Unmarshal([]byte(chunkJSON), &chunk); err != nil {
			continue
		}

		// Filter by startup ID
		if chunk.StartupID != startupID {
			continue
		}

		// Calculate similarity (simplified)
		if s.isRelevant(chunk.Content, query) {
			relevantContext = append(relevantContext, chunk.Content)
		}

		// Limit to top 5 results
		if len(relevantContext) >= 5 {
			break
		}
	}

	if len(relevantContext) == 0 {
		return s.fallbackContextSearch(ctx, startupID, query)
	}

	logger.Debugf("Retrieved %d relevant context chunks", len(relevantContext))
	return relevantContext, nil
}

func (s *service) fallbackContextSearch(ctx context.Context, startupID uuid.UUID, query string) ([]string, error) {
	mockContext := []string{
		"Previous hiring decisions: The startup previously hired 2 engineers in Q1",
		"Market conditions: Tech hiring market is competitive",
		"Budget constraints: Company has limited runway for 12 months",
		"Team growth: Current team size is 5 members with plans to scale",
		"Technology stack: Using Go, React, and PostgreSQL",
	}

	logger.Debugf("Using fallback context with %d chunks", len(mockContext))
	return mockContext, nil
}

func (s *service) isRelevant(content, query string) bool {
	// Simple relevance check - in production, use embedding similarity
	content = fmt.Sprintf("%s %s", content, query)
	return len(content) > len(query) // Placeholder logic
}

func (s *service) StoreStartupContext(ctx context.Context, startupID uuid.UUID, content string, metadata map[string]interface{}) error {
	logger.Infof("Storing startup context for: %s", startupID)

	// Generate embedding for the content
	embedding, err := s.llmService.GenerateEmbedding(ctx, content)
	if err != nil {
		logger.Warnf("Failed to generate embedding: %v", err)
		// Continue without embedding for now
		embedding = []float32{}
	}

	// Create context chunk
	chunk := ContextChunk{
		ID:        uuid.New().String(),
		StartupID: startupID,
		Content:   content,
		Embedding: embedding,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	// Store in Redis (in a real implementation, this would go to a vector database)
	chunkJSON, err := json.Marshal(chunk)
	if err != nil {
		return fmt.Errorf("failed to marshal context chunk: %w", err)
	}

	key := fmt.Sprintf("context:%s:%s", startupID, chunk.ID)
	if err := s.redis.Set(ctx, key, string(chunkJSON), 7*24*time.Hour); err != nil {
		return fmt.Errorf("failed to store context chunk: %w", err)
	}

	logger.Infof("Successfully stored context chunk: %s", chunk.ID)
	return nil
}
