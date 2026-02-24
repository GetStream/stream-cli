package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type HTTPClient struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

func NewHTTPClient(apiKey, apiSecret, baseURL string) *HTTPClient {
	return &HTTPClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (h *HTTPClient) createServerToken() (string, error) {
	claims := jwt.MapClaims{
		"server": true,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.apiSecret))
}

func (h *HTTPClient) DoRequest(ctx context.Context, method, path string, body interface{}) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/%s?api_key=%s", h.baseURL, path, h.apiKey)

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token, err := h.createServerToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth token: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("Stream-Auth-Type", "jwt")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

func DoTypedRequest[T any](h *HTTPClient, ctx context.Context, method, path string, body interface{}) (*T, error) {
	raw, err := h.DoRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &result, nil
}
