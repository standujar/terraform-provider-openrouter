package client

import "time"

type ApiKeyInfo struct {
	ID            string     `json:"hash"`
	Key           string     `json:"key,omitempty"`
	Label         string     `json:"label,omitempty"`
	Name          string     `json:"name"`
	IsProvisioner bool       `json:"is_provisioner,omitempty"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Limit         *float64   `json:"limit,omitempty"`
	LimitMinutes  *int       `json:"limit_minutes,omitempty"`
	Usage         float64    `json:"usage"`
	IsDisabled    bool       `json:"disabled"`
}

type ListApiKeysRequest struct {
	IncludeDisabled bool `url:"include_disabled,omitempty"`
	Offset          int  `url:"offset,omitempty"`
	Limit           int  `url:"limit,omitempty"`
}

type ListApiKeysResponse struct {
	Data []ApiKeyInfo `json:"data"`
}

type CreateApiKeyRequest struct {
	Name         string   `json:"name"`
	Limit        *float64 `json:"limit,omitempty"`
	LimitMinutes *int     `json:"limit_minutes,omitempty"`
}

type CreateApiKeyResponse struct {
	Data ApiKeyInfo `json:"data"`
	Key  string     `json:"key"`
}

type GetApiKeyResponse struct {
	Data ApiKeyInfo `json:"data"`
}

type UpdateApiKeyRequest struct {
	Name       *string     `json:"name,omitempty"`
	Limit      *float64    `json:"limit,omitempty"`
	IsDisabled *bool       `json:"disabled,omitempty"`
	BYOK       *BYOKConfig `json:"byok,omitempty"`
}

type BYOKConfig struct {
	Provider string  `json:"provider"`
	APIKey   string  `json:"api_key"`
	BaseURL  *string `json:"base_url,omitempty"`
}

type UpdateApiKeyResponse struct {
	Data ApiKeyInfo `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}