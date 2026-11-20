package api

import (
	"context"
	"log/slog"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/domain"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/service"
	"go.opentelemetry.io/otel"
)

var _ StrictServerInterface = (*ChatHandler)(nil)
var tracer = otel.Tracer("bazaargo-chatbot/internal/http")

type ChatHandler struct {
	svc    *service.ChatService
	logger *slog.Logger
}

func NewChatHandler(svc *service.ChatService, logger *slog.Logger) *ChatHandler {
	return &ChatHandler{svc: svc, logger: logger}
}

// Chat implements [StrictServerInterface].
func (c *ChatHandler) Chat(ctx context.Context, request ChatRequestObject) (ChatResponseObject, error) {
	ctx, span := tracer.Start(ctx, "ServerHandler.Chat")
	defer span.End()
	answer, err := c.svc.Ask(ctx, domain.ChatRequest{Question: request.Body.Question})
	if err != nil {
		return nil, err
	}
	return Chat200JSONResponse{
		Answer:  answer.Answer,
		Sources: answer.Sources,
	}, nil
}

// HealthCheck implements [StrictServerInterface].
func (c *ChatHandler) HealthCheck(ctx context.Context, request HealthCheckRequestObject) (HealthCheckResponseObject, error) {
	ctx, span := tracer.Start(ctx, "ServerHandler.HealthCheck")
	defer span.End()
	if err := c.svc.HealthCheck(ctx); err != nil {
		return HealthCheck500JSONResponse{
			PostgresHealth: false,
		}, nil
	}
	return HealthCheck200JSONResponse{
		PostgresHealth: true,
	}, nil
}
