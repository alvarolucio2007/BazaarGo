package gemini

import (
	"context"
	"fmt"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/ports"
	"google.golang.org/genai"
)

var _ ports.LLM = (*GeminiAdapter)(nil)

// GenerateAnswer implements [ports.LLM].
func (g *GeminiAdapter) GenerateAnswer(ctx context.Context, prompt string) (string, error) {
	contents := []*genai.Content{
		{Parts: []*genai.Part{{Text: prompt}}},
	}

	resp, err := g.client.Models.GenerateContent(ctx, "gemini-3.5-flash-lite", contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini: failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini: empty response from model")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}
