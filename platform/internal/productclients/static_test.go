package productclients

import (
	"context"
	"errors"
	"testing"
)

func TestStaticAuthenticator(t *testing.T) {
	authenticator := StaticAuthenticator{Tokens: map[string]string{"infra_link": "secret"}}
	client, err := authenticator.Authenticate(context.Background(), "infra_link", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if client.Audience != "infra_link" {
		t.Fatalf("unexpected audience %s", client.Audience)
	}
	if _, err := authenticator.Authenticate(context.Background(), "infra_link", "wrong"); !errors.Is(err, ErrInvalidProductClient) {
		t.Fatalf("expected invalid product client, got %v", err)
	}
}
