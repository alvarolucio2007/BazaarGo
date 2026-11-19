package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/domain"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/ports"
)

// ChatService implements the use cases related to the chatbot.
type ChatService struct {
	embedder ports.Embedder
	store    ports.VectorStore
	llm      ports.LLM
}

// NewChatService creates a new instance of ChatService.
func NewChatService(embedder ports.Embedder, vectorStore ports.VectorStore, llm ports.LLM) *ChatService {
	return &ChatService{embedder: embedder, store: vectorStore, llm: llm}
}

func (s *ChatService) Ask(ctx context.Context, req domain.ChatRequest) (domain.ChatResponse, error) {
	if err := req.Validate(); err != nil {
		return domain.ChatResponse{}, err
	}
	vector, err := s.embedder.GenerateEmbedding(ctx, req.Question)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("error while generating embedding: %w", err)
	}
	chunks, err := s.store.SearchSimilar(ctx, vector, 3)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("error while searching context: %w", err)
	}
	prompt := s.buildPrompt(req.Question, chunks)
	answer, err := s.llm.GenerateAnswer(ctx, prompt)
	if err != nil {
		return domain.ChatResponse{}, fmt.Errorf("error while generating response: %w", err)
	}
	var sources []string
	for _, c := range chunks {
		sources = append(sources, c.Title)
	}

	return domain.ChatResponse{
		Answer:  answer,
		Sources: sources,
	}, nil
}

func (s *ChatService) buildPrompt(question string, chunks []domain.DocumentChunk) string {
	var contexto strings.Builder
	for _, c := range chunks {
		fmt.Fprintf(&contexto, "- %s\n", c.Content)
	}

	return fmt.Sprintf(`Você é o assistente virtual do BazaarGo.
Responda à pergunta do usuário usando EXCLUSIVAMENTE o contexto abaixo.
Se a resposta não estiver no contexto, diga "Desculpe, não encontrei essa informação na base de conhecimento".

CONTEXTO:
%s

PERGUNTA DO USUÁRIO: %s`, contexto.String(), question) //prompt is in portuguese, because our userbase is portuguese (whos gonna be reading this anyway???)
}
