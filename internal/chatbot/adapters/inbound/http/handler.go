package api

import (
	"context"

	"github.com/alvarolucio2007/bazaargo/internal/chatbot/domain"
	"github.com/alvarolucio2007/bazaargo/internal/chatbot/service"
)

var _ StrictServerInterface = (*ChatHandler)(nil)

type ChatHandler struct {
	svc *service.ChatService
}

// Chat implements [StrictServerInterface].
func (c *ChatHandler) Chat(ctx context.Context, request ChatRequestObject) (ChatResponseObject, error) {
	answer, err := c.svc.Ask(ctx, domain.ChatRequest{Question: request.Body.Question})
	if err != nil {
		return nil, err
	}
	return Chat200JSONResponse{
		Answer:  answer.Answer,
		Sources: answer.Sources,
	}, nil
}
