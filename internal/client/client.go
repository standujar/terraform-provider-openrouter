package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://openrouter.ai/api/v1"
	defaultTimeout = 30 * time.Second
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string, baseURL *string) *Client {
	url := defaultBaseURL
	if baseURL != nil && *baseURL != "" {
		url = *baseURL
	}

	return &Client{
		baseURL: url,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *Client) GetCurrentApiKey(ctx context.Context) (*ApiKeyInfo, error) {
	var resp GetApiKeyResponse
	if err := c.doRequest(ctx, "GET", "/key", nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) ListApiKeys(ctx context.Context, params *ListApiKeysRequest) ([]ApiKeyInfo, error) {
	path := "/keys"
	if params != nil {
		v := url.Values{}
		if params.IncludeDisabled {
			v.Set("include_disabled", "true")
		}
		if params.Offset > 0 {
			v.Set("offset", fmt.Sprintf("%d", params.Offset))
		}
		if params.Limit > 0 {
			v.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if len(v) > 0 {
			path = path + "?" + v.Encode()
		}
	}

	var resp ListApiKeysResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) CreateApiKey(ctx context.Context, req *CreateApiKeyRequest) (*ApiKeyInfo, error) {
	var resp CreateApiKeyResponse
	if err := c.doRequest(ctx, "POST", "/keys", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) GetApiKey(ctx context.Context, hash string) (*ApiKeyInfo, error) {
	var resp GetApiKeyResponse
	if err := c.doRequest(ctx, "GET", "/keys/"+hash, nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) UpdateApiKey(ctx context.Context, hash string, req *UpdateApiKeyRequest) (*ApiKeyInfo, error) {
	var resp UpdateApiKeyResponse
	if err := c.doRequest(ctx, "PATCH", "/keys/"+hash, req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) DeleteApiKey(ctx context.Context, hash string) error {
	return c.doRequest(ctx, "DELETE", "/keys/"+hash, nil, nil)
}