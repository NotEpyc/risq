package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"risq_backend/pkg/llm"
	"risq_backend/pkg/logger"
	"risq_backend/types"

	"github.com/sashabaranov/go-openai"
)

type SpeculationResult struct {
	ProjectedRiskScore float64  `json:"projected_risk_score"`
	Confidence         float64  `json:"confidence"`
	Suggestions        []string `json:"suggestions"`
	Reasoning          string   `json:"reasoning"`
}

type Service interface {
	SpeculateDecision(ctx context.Context, input types.DecisionInput, contextChunks []string, currentRisk float64) (*SpeculationResult, error)
	GenerateInitialRiskProfile(ctx context.Context, startupInfo string) (*InitialRiskResult, error)
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

type InitialRiskResult struct {
	RiskScore   float64  `json:"risk_score"`
	Factors     []string `json:"factors"`
	Suggestions []string `json:"suggestions"`
	Reasoning   string   `json:"reasoning"`
}

type service struct {
	llmClient *llm.Client
}

func NewService(llmClient *llm.Client) Service {
	return &service{llmClient: llmClient}
}

func (s *service) SpeculateDecision(ctx context.Context, input types.DecisionInput, contextChunks []string, currentRisk float64) (*SpeculationResult, error) {
	logger.Infof("Generating speculation for decision: %s", input.Description)

	// Build context from chunks
	contextStr := strings.Join(contextChunks, "\n")
	if contextStr == "" {
		contextStr = "No relevant historical context available."
	}

	// Build prompt for speculation
	prompt := fmt.Sprintf(`You are an expert startup risk analyst. A startup founder is considering the following decision:

DECISION: %s
CATEGORY: %s
CONTEXT: %s
TIMELINE: %s
BUDGET: $%.2f

CURRENT RISK SCORE: %.2f (0-100 scale, where 0 is no risk and 100 is extremely high risk)

HISTORICAL CONTEXT:
%s

Please analyze this decision and provide:
1. Projected risk score (0-100) after implementing this decision
2. Confidence level (0-1) in your assessment
3. 3-5 specific risk mitigation suggestions
4. Detailed reasoning for your assessment

Format your response as JSON:
{
  "projected_risk_score": <number>,
  "confidence": <number>,
  "suggestions": ["suggestion1", "suggestion2", ...],
  "reasoning": "detailed explanation"
}`, input.Description, input.Category, input.Context, input.Timeline, input.Budget, currentRisk, contextStr)

	// Call LLM
	response, err := s.llmClient.Chat(ctx, llm.ChatRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert startup risk analyst. Always respond with valid JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Model:       openai.GPT4oMini,
		MaxTokens:   1000,
		Temperature: 0.3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	// Parse JSON response
	var result SpeculationResult
	if err := json.Unmarshal([]byte(response.Content), &result); err != nil {
		logger.Warnf("Failed to parse LLM JSON response: %v", err)
		// Fallback to basic parsing
		result = SpeculationResult{
			ProjectedRiskScore: currentRisk + 5, // Simple fallback
			Confidence:         0.5,
			Suggestions:        []string{"Consider conducting more research", "Seek expert advice", "Start with a small pilot"},
			Reasoning:          "Unable to parse detailed analysis from LLM response",
		}
	}

	// Validate and bound the results
	if result.ProjectedRiskScore < 0 {
		result.ProjectedRiskScore = 0
	}
	if result.ProjectedRiskScore > 100 {
		result.ProjectedRiskScore = 100
	}
	if result.Confidence < 0 {
		result.Confidence = 0
	}
	if result.Confidence > 1 {
		result.Confidence = 1
	}

	logger.Infof("Generated speculation - Risk: %.2f, Confidence: %.2f", result.ProjectedRiskScore, result.Confidence)
	return &result, nil
}

func (s *service) GenerateInitialRiskProfile(ctx context.Context, startupInfo string) (*InitialRiskResult, error) {
	logger.Infof("Generating initial risk profile for startup")

	prompt := fmt.Sprintf(`You are an expert startup risk analyst. Please analyze the following startup information and provide an initial risk assessment:

STARTUP INFORMATION:
%s

Please provide:
1. Overall risk score (0-100 scale)
2. Key risk factors identified
3. Specific risk mitigation suggestions
4. Detailed reasoning for the assessment

Format your response as JSON:
{
  "risk_score": <number>,
  "factors": ["factor1", "factor2", ...],
  "suggestions": ["suggestion1", "suggestion2", ...],
  "reasoning": "detailed explanation"
}`, startupInfo)

	response, err := s.llmClient.Chat(ctx, llm.ChatRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert startup risk analyst. Always respond with valid JSON.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Model:       openai.GPT4oMini,
		MaxTokens:   1000,
		Temperature: 0.3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM response: %w", err)
	}

	var result InitialRiskResult
	if err := json.Unmarshal([]byte(response.Content), &result); err != nil {
		logger.Warnf("Failed to parse LLM JSON response: %v", err)
		// Fallback
		result = InitialRiskResult{
			RiskScore:   50, // Medium risk as default
			Factors:     []string{"Market uncertainty", "Team inexperience", "Funding challenges"},
			Suggestions: []string{"Conduct market research", "Build team expertise", "Develop funding strategy"},
			Reasoning:   "Unable to parse detailed analysis from LLM response",
		}
	}

	// Validate and bound the results
	if result.RiskScore < 0 {
		result.RiskScore = 0
	}
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}

	logger.Infof("Generated initial risk profile - Score: %.2f", result.RiskScore)
	return &result, nil
}

func (s *service) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	logger.Debugf("Generating embedding for text length: %d", len(text))

	response, err := s.llmClient.Embedding(ctx, llm.EmbeddingRequest{
		Input: []string{text},
		Model: "text-embedding-ada-002",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	if len(response.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return response.Embeddings[0], nil
}
