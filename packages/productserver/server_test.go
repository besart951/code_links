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

func TestHandlerAcceptsProductRoutes(t *testing.T) {
	handler := Handler(Config{ProductID: "planer-link"}, nil, func(mux *http.ServeMux, tools RouteTools) {
		mux.HandleFunc("GET /api/product-state", func(w http.ResponseWriter, _ *http.Request) {
			writeJSON(w, http.StatusOK, map[string]string{"productId": tools.Config.ProductID})
		})
	})
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/product-state", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != "{\"productId\":\"planer-link\"}\n" {
		t.Fatalf("unexpected body: %s", recorder.Body.String())
	}
}
