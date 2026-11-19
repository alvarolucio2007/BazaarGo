package ports

import "context"

type Embedder interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}
