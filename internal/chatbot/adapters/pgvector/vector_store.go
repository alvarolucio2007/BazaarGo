package pgvector

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/domain"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/ports"
	"github.com/pgvector/pgvector-go"
)

var _ ports.VectorStore = (*VectorStore)(nil)

type VectorStore struct {
	db     *sql.DB
	logger *slog.Logger
}

// SearchSimilar implements [ports.VectorStore].
// It performs a cosine similarity search in the PostgreSQL pgvector database.
// It retrieves the most relevant document chunks based on the provided embedding.
//
// Returns a slice of [domain.DocumentChunk] ordered by similarity, or an error if the query fails.
func (v *VectorStore) SearchSimilar(ctx context.Context, embedding []float32, limit int) ([]domain.DocumentChunk, error) {
	ctx, cancel := context.WithTimeout(ctx, PGVectorQueryTimeout)
	defer cancel()
	query := "SELECT title, content FROM documents ORDER BY embedding <=> $1 LIMIT $2;"
	rows, err := v.db.QueryContext(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("postgres: search similar: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			v.logger.Error("postgres: rows.close search similar", "error", err)
		}
	}()
	documents := make([]domain.DocumentChunk, 0, limit)
	for rows.Next() {
		var document domain.DocumentChunk
		if err := rows.Scan(&document.Title, &document.Content); err != nil {
			return nil, fmt.Errorf("postgres: scan documents: %w", err)
		}
		documents = append(documents, document)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: iterate documents: %w", err)
	}
	return documents, nil
}
