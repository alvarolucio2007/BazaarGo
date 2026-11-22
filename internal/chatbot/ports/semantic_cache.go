package ports

import "context"

type SemanticCache interface {
	GetSimilar(ctx context.Context, embedding []float32, threshold float64) (string, bool, error)
	Save(ctx context.Context, question string, embedding []float32, answer string) error
}
