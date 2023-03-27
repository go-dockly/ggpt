package entity

import (
	"github.com/go-dockly/ggpt/pkg/gpt"
	"github.com/sashabaranov/go-openai"
)

type Session struct {
	CurrentChat *openai.ChatCompletionRequest            `json:"current"`
	Chats       map[string]*openai.ChatCompletionRequest `json:"chats"`
	Persona     string                                   `json:"persona"`
	Personas    map[string]*Persona                      `json:"personas"`
	Cost        map[string]*gpt.CostGaugeMeta            `json:"cost"`
}
