package identify_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1909-pramod/identify-sdk/go/identify"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*identify.Client, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	return identify.New(srv.URL), srv.Close
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func TestLoginWithPassword(t *testing.T) {
	client, close := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"token": "test-token", "token_type": "Bearer", "expires_in": 3600,
		})
	})
	defer close()

	resp, err := client.LoginWithPassword(context.Background(), "user@example.com", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "test-token" {
		t.Errorf("got token %q, want %q", resp.Token, "test-token")
	}
}

func TestValidateToken(t *testing.T) {
	client, close := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"valid": true,
			"identity": map[string]any{
				"id": 1, "primary_identifier": "user@example.com",
			},
		})
	})
	defer close()

	resp, err := client.ValidateToken(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Valid {
		t.Error("expected valid=true")
	}
	if resp.Identity.PrimaryIdentifier != "user@example.com" {
		t.Errorf("unexpected identifier: %s", resp.Identity.PrimaryIdentifier)
	}
}

func TestGetMe(t *testing.T) {
	client, close := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"id": 1, "primary_identifier": "user@example.com",
		})
	})
	defer close()

	identity, err := client.GetMe(context.Background(), "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.PrimaryIdentifier != "user@example.com" {
		t.Errorf("unexpected identifier: %s", identity.PrimaryIdentifier)
	}
}

func TestErrorResponse(t *testing.T) {
	client, close := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
	})
	defer close()

	_, err := client.LoginWithPassword(context.Background(), "user@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var idErr *identify.Error
	if !errors.As(err, &idErr) {
		t.Fatalf("expected *identify.Error, got %T", err)
	}
	if idErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("got status %d, want 401", idErr.StatusCode)
	}
}
