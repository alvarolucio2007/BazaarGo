package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/alvarolucio2007/bazaargo/internal/app"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"google.golang.org/genai"
)

type FAQDocument struct {
	Brand         string      `json:"brand"`
	Domain        string      `json:"domain"`
	SchemaVersion string      `json:"schema_version"`
	Locale        string      `json:"locale"`
	RAGMetadata   RAGMetadata `json:"rag_metadata"`
	FAQs          []FAQItem   `json:"faqs"`
}

type RAGMetadata struct {
	RecommendedChunkStrategy string   `json:"recommended_chunk_strategy"`
	FilterableFields         []string `json:"filterable_fields"`
}

type FAQItem struct {
	ID          string   `json:"id"`
	Category    string   `json:"category"`
	Question    string   `json:"question"`
	Answer      string   `json:"answer"`
	Keywords    []string `json:"keywords"`
	Priority    string   `json:"priority"`
	LastUpdated string   `json:"last_updated"`
}

func ParseFAQDocument(data []byte) (*FAQDocument, error) {
	var doc FAQDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal FAQ JSON: %w", err)
	}
	return &doc, nil
}

func (f FAQItem) EmbeddingText() string {
	return fmt.Sprintf(
		"Category: %s | Question: %s | Answer: %s | Keywords: %s",
		f.Category,
		f.Question,
		f.Answer,
		strings.Join(f.Keywords, ", "),
	)
}

// TODO: make this more agnostic for docker i guess
func main() {
	if err := app.LoadEnv(".env"); err != nil {
		log.Fatalf("main:error while trying to read env: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	text, err := os.ReadFile("data/knowledge/bazaargo.json")
	if err != nil {
		log.Fatalf("Error while reading JSON: %v", err)
	}
	faqDocuments, err := ParseFAQDocument(text)
	if err != nil {
		log.Fatalf("Error while parsing JSON: %v", err)
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("api key empty bro check yo .env")
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatalf("Error while creating Gemini client: %v", err)
	}

	dbURL := os.Getenv("DATABASE_LOCAL_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:root@localhost:5432/bazaargo?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Erro ao parsear DB URL: %v", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}
	defer pool.Close()

	query := "INSERT INTO documents (title, content, embedding) VALUES ($1, $2, $3)"
	dimensions := int32(384)

	_, err = pool.Exec(ctx, "TRUNCATE TABLE documents")
	if err != nil {
		log.Printf("WARNING: Wasn't able to truncate table %v", err)
	}

	for _, faq := range faqDocuments.FAQs {
		log.Printf("Creating embedding for: %s", faq.Question)

		textToEmbed := faq.EmbeddingText()
		contents := []*genai.Content{
			genai.NewContentFromText(textToEmbed, genai.RoleUser),
		}

		result, err := client.Models.EmbedContent(
			ctx,
			"gemini-embedding-001",
			contents,
			&genai.EmbedContentConfig{OutputDimensionality: &dimensions},
		)
		if err != nil {
			log.Printf("Gemini error for FAQ %s: %v", faq.ID, err)
			continue
		}

		vetor := result.Embeddings[0].Values

		_, err = pool.Exec(ctx, query, faq.Question, textToEmbed, pgvector.NewVector(vetor))
		if err != nil {
			log.Printf("DB insertion error for FAQ %s: %v", faq.ID, err)
			continue
		}
	}

	log.Println("Finished.")
}
