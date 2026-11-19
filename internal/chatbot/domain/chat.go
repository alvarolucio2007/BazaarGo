package domain

import (
	"strings"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message represents one single message in the conversation.
type Message struct {
	Role      Role
	Content   string
	Timestamp time.Time
}

// ChatRequest is what the user sends to the API.
type ChatRequest struct {
	Question string
}

// ChatResponse is what the API responds back.
type ChatResponse struct {
	Answer  string
	Sources []string
}

// Validate checks whether a chat request is valid or not.
//
// It currently only checks whether it's empty and returns an [ErrEmptyQuestion] if so.
func (c *ChatRequest) Validate() error {
	if strings.TrimSpace(c.Question) == "" {
		return ErrEmptyQuestion
	}
	return nil
}
