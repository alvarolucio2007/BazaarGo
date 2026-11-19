package ports

import (
	"context"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/domain"
)

type VectorStore interface {
	SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]domain.DocumentChunk, error)
}
