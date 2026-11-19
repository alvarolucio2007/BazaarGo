package ports

import "context"

type LLM interface {
	GenerateAnswer(ctx context.Context, prompt string) (string, error)
}
