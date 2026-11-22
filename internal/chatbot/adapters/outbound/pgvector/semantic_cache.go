package pgvector

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/ports"
	"github.com/pgvector/pgvector-go"
	"go.opentelemetry.io/otel/codes"
)

var _ ports.SemanticCache = (*SemanticCache)(nil)

type SemanticCache struct {
	db *sql.DB
}

func NewSemanticCache(db *sql.DB) *SemanticCache {
	return &SemanticCache{db: db}
}

// GetSimilar implements [ports.SemanticCache].
func (s *SemanticCache) GetSimilar(ctx context.Context, embedding []float32, threshold float64) (string, bool, error) {
	ctx, span := tracer.Start(ctx, "SemanticCache.GetSimilar")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, PGVectorQueryTimeout)
	defer cancel()
	var answer string
	var similarity float64
	query := `
			SELECT answer, (1 - (question_embedding <=> $1)) AS similarity
			FROM semantic_cache
			WHERE (1 - (question_embedding <=> $1)) >= $2
			ORDER BY question_embedding <=> $1
			LIMIT 1;
		`
	err := s.db.QueryRowContext(ctx, query, pgvector.NewVector(embedding), threshold).Scan(&answer, &similarity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, nil // Cache miss
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return "", false, fmt.Errorf("postgres: error getting similar cache: %w", err)
	}
	return answer, true, nil
}

// Save implements [ports.SemanticCache].
func (s *SemanticCache) Save(ctx context.Context, question string, embedding []float32, answer string) error {
	ctx, span := tracer.Start(ctx, "SemanticCache.Save")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, PGVectorQueryTimeout)
	defer cancel()
	query := `INSERT INTO semantic_cache (question_embedding,question_text,answer) VALUES ($1,$2,$3)`
	if err := s.db.QueryRowContext(ctx, query, pgvector.NewVector(embedding), question, answer).Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}
