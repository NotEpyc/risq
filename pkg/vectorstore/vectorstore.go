package vectorstore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"risq_backend/pkg/cache"
	"risq_backend/pkg/logger"

	"github.com/google/uuid"
)

type Vector struct {
	ID        string                 `json:"id"`
	Data      []float32              `json:"data"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

type SearchResult struct {
	Vector *Vector `json:"vector"`
	Score  float64 `json:"score"`
}

type Store interface {
	Store(ctx context.Context, vector *Vector) error
	Search(ctx context.Context, queryVector []float32, topK int, filters map[string]interface{}) ([]*SearchResult, error)
	Delete(ctx context.Context, id string) error
}

// Redis-based vector store implementation (simplified)
// In production, use Redis Vector or Pinecone
type redisVectorStore struct {
	redis *cache.Redis
}

func NewRedisVectorStore(redis *cache.Redis) Store {
	return &redisVectorStore{redis: redis}
}

func (r *redisVectorStore) Store(ctx context.Context, vector *Vector) error {
	logger.Debugf("Storing vector: %s", vector.ID)

	// Serialize vector
	vectorData, err := json.Marshal(vector)
	if err != nil {
		return fmt.Errorf("failed to marshal vector: %w", err)
	}

	// Store in Redis with expiration
	key := fmt.Sprintf("vector:%s", vector.ID)
	if err := r.redis.Set(ctx, key, string(vectorData), 7*24*time.Hour); err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	// Add to searchable index (simplified - in production use Redis Vector)
	indexKey := "vector_index"
	if err := r.redis.Client().SAdd(ctx, indexKey, vector.ID).Err(); err != nil {
		logger.Warnf("Failed to add to vector index: %v", err)
	}

	logger.Debugf("Successfully stored vector: %s", vector.ID)
	return nil
}

func (r *redisVectorStore) Search(ctx context.Context, queryVector []float32, topK int, filters map[string]interface{}) ([]*SearchResult, error) {
	logger.Debugf("Searching vectors with topK: %d", topK)

	// Simplified search implementation
	// In production, this would use proper vector similarity search

	// Get all vector IDs from index
	indexKey := "vector_index"
	vectorIDs, err := r.redis.Client().SMembers(ctx, indexKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get vector index: %w", err)
	}

	results := make([]*SearchResult, 0)

	// Load and compare vectors (simplified implementation)
	for _, id := range vectorIDs {
		vectorKey := fmt.Sprintf("vector:%s", id)
		vectorData, err := r.redis.Get(ctx, vectorKey)
		if err != nil {
			continue
		}

		var vector Vector
		if err := json.Unmarshal([]byte(vectorData), &vector); err != nil {
			continue
		}

		// Apply filters
		if !r.matchesFilters(&vector, filters) {
			continue
		}

		// Calculate similarity (simplified cosine similarity)
		score := r.cosineSimilarity(queryVector, vector.Data)

		results = append(results, &SearchResult{
			Vector: &vector,
			Score:  score,
		})
	}

	// Sort by score (descending) and limit
	r.sortResultsByScore(results)
	if len(results) > topK {
		results = results[:topK]
	}

	logger.Debugf("Found %d similar vectors", len(results))
	return results, nil
}

func (r *redisVectorStore) Delete(ctx context.Context, id string) error {
	logger.Debugf("Deleting vector: %s", id)

	// Remove from storage
	vectorKey := fmt.Sprintf("vector:%s", id)
	if err := r.redis.Del(ctx, vectorKey); err != nil {
		return fmt.Errorf("failed to delete vector: %w", err)
	}

	// Remove from index
	indexKey := "vector_index"
	if err := r.redis.Client().SRem(ctx, indexKey, id).Err(); err != nil {
		logger.Warnf("Failed to remove from vector index: %v", err)
	}

	logger.Debugf("Successfully deleted vector: %s", id)
	return nil
}

func (r *redisVectorStore) matchesFilters(vector *Vector, filters map[string]interface{}) bool {
	if filters == nil || len(filters) == 0 {
		return true
	}

	for key, value := range filters {
		if metadataValue, exists := vector.Metadata[key]; !exists || metadataValue != value {
			return false
		}
	}

	return true
}

func (r *redisVectorStore) cosineSimilarity(vec1, vec2 []float32) float64 {
	if len(vec1) != len(vec2) {
		return 0.0
	}

	var dotProduct, norm1, norm2 float64

	for i := 0; i < len(vec1); i++ {
		dotProduct += float64(vec1[i] * vec2[i])
		norm1 += float64(vec1[i] * vec1[i])
		norm2 += float64(vec2[i] * vec2[i])
	}

	if norm1 == 0.0 || norm2 == 0.0 {
		return 0.0
	}

	return dotProduct / (norm1 * norm2)
}

func (r *redisVectorStore) sortResultsByScore(results []*SearchResult) {
	// Simple bubble sort by score (descending)
	n := len(results)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
}

// Helper functions to create vectors
func NewVector(id string, data []float32, metadata map[string]interface{}) *Vector {
	return &Vector{
		ID:        id,
		Data:      data,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
}

func NewDecisionVector(decisionID uuid.UUID, startupID uuid.UUID, embedding []float32, decisionText string) *Vector {
	return &Vector{
		ID:   fmt.Sprintf("decision:%s", decisionID.String()),
		Data: embedding,
		Metadata: map[string]interface{}{
			"type":        "decision",
			"decision_id": decisionID.String(),
			"startup_id":  startupID.String(),
			"content":     decisionText,
		},
		Timestamp: time.Now(),
	}
}

func NewStartupVector(startupID uuid.UUID, embedding []float32, startupInfo string) *Vector {
	return &Vector{
		ID:   fmt.Sprintf("startup:%s", startupID.String()),
		Data: embedding,
		Metadata: map[string]interface{}{
			"type":       "startup",
			"startup_id": startupID.String(),
			"content":    startupInfo,
		},
		Timestamp: time.Now(),
	}
}
