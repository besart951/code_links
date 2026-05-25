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

func TestCORSAllowsOnlyConfiguredOriginsWithoutCredentials(t *testing.T) {
	handler := Handler(Config{ProductID: "infra-link", AllowedOrigins: []string{"https://app.example.com"}}, nil)

	allowed := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	request.Header.Set("Origin", "https://app.example.com")
	handler.ServeHTTP(allowed, request)

	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Fatalf("expected allowed origin, got %q", got)
	}
	if got := allowed.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("expected no credentials header, got %q", got)
	}

	disallowed := httptest.NewRecorder()
	disallowedRequest := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	disallowedRequest.Header.Set("Origin", "https://evil.example.com")
	handler.ServeHTTP(disallowed, disallowedRequest)
	if got := disallowed.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected disallowed origin to get no ACAO, got %q", got)
	}
}

func TestCORSOmitsHeadersWhenOriginMissing(t *testing.T) {
	handler := Handler(Config{ProductID: "infra-link", AllowedOrigins: []string{"https://app.example.com"}}, nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/health", nil))

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no ACAO for missing origin, got %q", got)
	}
}
