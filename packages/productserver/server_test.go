package productserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthExposesProductID(t *testing.T) {
	handler := Handler(Config{ProductID: "infra-link"}, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != "{\"productId\":\"infra-link\",\"status\":\"ok\"}\n" &&
		recorder.Body.String() != "{\"status\":\"ok\",\"productId\":\"infra-link\"}\n" {
		t.Fatalf("unexpected body: %s", recorder.Body.String())
	}
}
