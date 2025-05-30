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

func MessagesToMessageContent(messages ...api.Message) []llms.MessageContent {
	messageContents := []llms.MessageContent{}

	for _, msgContent := range messages {
		messageContents = append(messageContents, *msgContent.Message)
	}

	return messageContents
}
