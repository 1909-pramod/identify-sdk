// Package identify provides a client for the Identify auth service.
package identify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client is an HTTP client for the Identify auth service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new Client pointing at baseURL (e.g. "http://localhost:8080").
func New(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Error is returned when the server responds with a non-2xx status code.
type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("identify: HTTP %d: %s", e.StatusCode, e.Message)
}

// Identity is the resolved user record returned by the server.
type Identity struct {
	ID                   int       `json:"id"`
	PrimaryIdentifier    string    `json:"primary_identifier"`
	IdentifierGroupName  string    `json:"identifier_group_name,omitempty"`
	IdentifierGroupType  string    `json:"identifier_group_type,omitempty"`
	UserType             string    `json:"user_type,omitempty"`
	ExpiryAt             *string   `json:"expiry_at,omitempty"`
	Roles                []string  `json:"roles,omitempty"`
	IsSuperAdmin         bool      `json:"is_super_admin,omitempty"`
}

// TokenResponse is returned when a credential is exchanged for a JWT.
type TokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

// ValidateResponse is returned when a JWT is validated.
type ValidateResponse struct {
	Valid      bool      `json:"valid"`
	Identity   Identity  `json:"identity"`
	ExpiresAt  string    `json:"expires_at,omitempty"`
}

// LoginWithPassword exchanges a username and password for a JWT.
func (c *Client) LoginWithPassword(ctx context.Context, username, password string) (*TokenResponse, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return c.post(ctx, "/auth", "Basic "+encoded)
}

// LoginWithGoogle exchanges a Google ID token for an internal JWT.
func (c *Client) LoginWithGoogle(ctx context.Context, googleIDToken string) (*TokenResponse, error) {
	return c.post(ctx, "/auth", "Bearer "+googleIDToken)
}

// RefreshToken exchanges an existing JWT for a fresh one.
func (c *Client) RefreshToken(ctx context.Context, token string) (*TokenResponse, error) {
	return c.post(ctx, "/auth", "Bearer "+token)
}

// ValidateToken validates a JWT and returns the resolved identity.
func (c *Client) ValidateToken(ctx context.Context, token string) (*ValidateResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/token", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var out ValidateResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMe returns the identity for an authenticated token.
func (c *Client) GetMe(ctx context.Context, token string) (*Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/me", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var out Identity
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) post(ctx context.Context, path, authHeader string) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)

	var out TokenResponse
	if err := c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("identify: failed to decode response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := body["error"].(string)
		if msg == "" {
			msg = "request failed"
		}
		return &Error{StatusCode: resp.StatusCode, Message: msg}
	}

	// Re-encode and decode into the typed output struct.
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
