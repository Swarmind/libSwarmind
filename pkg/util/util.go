package util

import (
	"github.com/Swarmind/libswarmind/pkg/api"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
)

func MessageContentToMessages(messageContent ...llms.MessageContent) []api.Message {
	messages := []api.Message{}

	for _, msg := range messageContent {
		messages = append(messages, api.Message{
			ID:      uuid.NewString(),
			Message: &msg,
		})
	}

	return messages
}
