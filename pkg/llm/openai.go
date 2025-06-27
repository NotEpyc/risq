package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client *openai.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	client := openai.NewClient(apiKey)
	return &Client{
		client: client,
		apiKey: apiKey,
	}
}

type ChatRequest struct {
	Messages    []openai.ChatCompletionMessage
	Model       string
	MaxTokens   int
	Temperature float32
}

type ChatResponse struct {
	Content string
	Usage   openai.Usage
}

func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = openai.GPT4oMini
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 1000
	}
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    req.Messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from OpenAI")
	}

	return &ChatResponse{
		Content: resp.Choices[0].Message.Content,
		Usage:   resp.Usage,
	}, nil
}

type EmbeddingRequest struct {
	Input []string
	Model string
}

type EmbeddingResponse struct {
	Embeddings [][]float32
	Usage      openai.Usage
}

func (c *Client) Embedding(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if req.Model == "" {
		req.Model = "text-embedding-ada-002"
	}

	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: req.Input,
		Model: openai.EmbeddingModel(req.Model),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return &EmbeddingResponse{
		Embeddings: embeddings,
		Usage:      resp.Usage,
	}, nil
}
