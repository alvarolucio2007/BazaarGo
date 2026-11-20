package gemini

import (
	"go.opentelemetry.io/otel"
	"google.golang.org/genai"
)

type GeminiAdapter struct {
	client *genai.Client
}

// NewGeminiAdapter creates a new gemini connection with [*genai.Client].
func NewGeminiAdapter(client *genai.Client) *GeminiAdapter {
	return &GeminiAdapter{client: client}
}

var tracer = otel.Tracer("bazaargo-chat/internal/gemini")
