package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	adminsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/admin"
	authsvc "github.com/besart951/code-links/apps/auth-service/backend/internal/auth"
)

func TestWriteAuthErrorMapping(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeAuthError(recorder, &authsvc.Error{Kind: authsvc.KindConflict, Message: "email already exists"})

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected conflict status, got %d", recorder.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["code"] != string(authsvc.KindConflict) || body["error"] != "email already exists" {
		t.Fatalf("unexpected auth error body: %#v", body)
	}
}

func TestWriteAdminErrorMapping(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeAdminError(recorder, &adminsvc.Error{Kind: adminsvc.KindNotFound, Message: "user not found"})

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected not found status, got %d", recorder.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["code"] != string(adminsvc.KindNotFound) || body["error"] != "user not found" {
		t.Fatalf("unexpected admin error body: %#v", body)
	}
}
