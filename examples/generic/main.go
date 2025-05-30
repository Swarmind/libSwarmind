package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Swarmind/libswarmind/pkg/api"
	"github.com/Swarmind/libswarmind/pkg/util"
	"github.com/tmc/langchaingo/llms"
)

func main() {
	swarmindAPI := api.SwarmindAPI{
		Namespace: "test",
		Token:     os.Getenv("LIBSWARMIND_TOKEN"),
		URL:       os.Getenv("LIBSWARMIND_API_URL"),
		Client:    &http.Client{},
	}
	ctx := context.Background()

	heads, err := swarmindAPI.GetHeads(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Swarmind API Heads: %s\n\n", heads)

	chatName := "testing"
	historyMessageContent := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeAI, "Message 1"),
		llms.TextParts(llms.ChatMessageTypeHuman, "Message 2"),
		llms.TextParts(llms.ChatMessageTypeAI, "Message 3"),
		llms.TextParts(llms.ChatMessageTypeHuman, "Message 4"),
	}
	historyMessages := util.MessageContentToMessages(historyMessageContent...)
	toEdit := historyMessages[1]
	toDelete := historyMessages[3]

	err = swarmindAPI.UpdateHistory(ctx, chatName, historyMessages...)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pushed new messages")

	chats, err := swarmindAPI.GetChats(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Chat list: %s\n\n", chats)

	messages, err := swarmindAPI.GetHistory(ctx, chatName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Messages stored:")
	for _, msg := range messages {
		fmt.Printf(
			"ID: %s\nCreatedAt: %s, UpdatedAt: %s\nContent: %s\n\n",
			msg.ID,
			msg.CreatedAt.String(),
			msg.UpdatedAt.String(),
			*msg.Message,
		)
	}

	toEdit.Message.Parts[0] = llms.TextPart("EDITED MESSAGE")
	toDelete.Message = nil
	err = swarmindAPI.UpdateHistory(ctx, chatName, toEdit, toDelete)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated messages (edit + delete)\n\n")

	messages, err = swarmindAPI.GetHistory(ctx, chatName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Changed messages:")
	for _, msg := range messages {
		fmt.Printf(
			"ID: %s\nCreatedAt: %s, UpdatedAt: %s\nContent: %s\n\n",
			msg.ID,
			msg.CreatedAt.String(),
			msg.UpdatedAt.String(),
			*msg.Message,
		)
	}

	err = swarmindAPI.DropHistory(ctx, chatName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Dropped chat")

	chats, err = swarmindAPI.GetChats(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Chat list: %s\n", chats)
}
