package authuser

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientLoadsCurrentUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/userinfo/me" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token-1" {
			t.Fatalf("unexpected Authorization header %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode(UserSnapshot{
			ID:            "user-1",
			Email:         "demo@codelinks.dev",
			Name:          "Demo User",
			Status:        "active",
			EmailVerified: true,
			Licenses:      []string{"infra-link"},
		})
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}

	user, err := client.Me(context.Background(), "token-1")
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != "user-1" || user.Email != "demo@codelinks.dev" || !user.EmailVerified {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestClientLookupUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/userinfo/lookup" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		var request lookupRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if len(request.UserIDs) != 1 || request.UserIDs[0] != "user-1" {
			t.Fatalf("unexpected lookup request: %#v", request)
		}
		_ = json.NewEncoder(w).Encode([]UserCard{{ID: "user-1", Name: "Demo User"}})
	}))
	defer server.Close()

	client, err := NewClient(Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}

	cards, err := client.Lookup(context.Background(), "token-1", []string{"user-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != 1 || cards[0].Name != "Demo User" {
		t.Fatalf("unexpected cards: %#v", cards)
	}
}
