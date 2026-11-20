package gemini

import (
	"context"
	"fmt"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/ports"
	"google.golang.org/genai"
)

var _ ports.Embedder = (*GeminiAdapter)(nil)

// GenerateEmbedding implements [ports.Embedder].
func (g *GeminiAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	ctx, span := tracer.Start(ctx, "GeminiAdapter.GenerateEmbedding")
	defer span.End()
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: text}}},
	}
	dimensions := int32(384)
	config := &genai.EmbedContentConfig{
		OutputDimensionality: &dimensions,
	}
	resp, err := g.client.Models.EmbedContent(ctx, "gemini-embedding-001", contents, config)
	if err != nil {
		return nil, fmt.Errorf("gemini: failed to generate embedding: %w", err)
	}
	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("gemini: no embeddings returned")
	}
	return resp.Embeddings[0].Values, nil
}
