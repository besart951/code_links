package httpapi

import (
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
}

func TestWriteAdminErrorMapping(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeAdminError(recorder, &adminsvc.Error{Kind: adminsvc.KindNotFound, Message: "user not found"})

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected not found status, got %d", recorder.Code)
	}
}
