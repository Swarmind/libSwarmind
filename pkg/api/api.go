package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tmc/langchaingo/llms"
)

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID        string               `json:"id"`
	CreatedAt time.Time            `json:"created_at,omitempty"`
	UpdatedAt time.Time            `json:"updated_at,omitempty"`
	Message   *llms.MessageContent `json:"message"`
}

type SwarmindAPI struct {
	Namespace string
	Token     string
	URL       string
	Client    *http.Client
}

func (s SwarmindAPI) GetHeads(ctx context.Context) ([]string, error) {
	requestURL, err := url.JoinPath(s.URL, "heads")
	if err != nil {
		return nil, fmt.Errorf("failed to process request URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d status code response: %s", resp.StatusCode, string(body))
	}

	var heads []string
	if err := json.NewDecoder(resp.Body).Decode(&heads); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return heads, nil
}

func (s SwarmindAPI) GetChats(ctx context.Context) ([]string, error) {
	requestURL, err := url.JoinPath(s.URL, "store", "chats", s.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to process request URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d status code response: %s", resp.StatusCode, string(body))
	}

	var chats []string
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return chats, nil
}

func (s SwarmindAPI) GetHistory(ctx context.Context, chat string) ([]Message, error) {
	requestURL, err := url.JoinPath(s.URL, "store", "history", s.Namespace, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to process request URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d status code response: %s", resp.StatusCode, string(body))
	}

	var messages []Message
	if err := json.NewDecoder(resp.Body).Decode(&messages); err != nil {
		return nil, fmt.Errorf("decoding response body: %w", err)
	}

	return messages, nil
}

func (s SwarmindAPI) UpdateHistory(ctx context.Context, chat string, messages ...Message) error {
	requestURL, err := url.JoinPath(s.URL, "store", "history", s.Namespace, chat)
	if err != nil {
		return fmt.Errorf("failed to process request URL: %w", err)
	}
	jsonBytes, err := json.Marshal(Messages{Messages: messages})
	if err != nil {
		return fmt.Errorf("marshaling messages: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%d status code response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s SwarmindAPI) DropHistory(ctx context.Context, chat string) error {
	requestURL, err := url.JoinPath(s.URL, "store", "history", s.Namespace, chat)
	if err != nil {
		return fmt.Errorf("failed to process request URL: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	resp, err := s.Client.Do(req)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%d status code response: %s", resp.StatusCode, string(body))
	}

	return nil
}
