package productserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/besart951/code-links/packages/productauth"
)

type Config struct {
	Port           string
	ProductID      string
	JWKSURL        string
	Issuer         string
	Audience       string
	AllowedOrigins []string
}

type RouteRegistrar func(*http.ServeMux, RouteTools)

type RouteTools struct {
	Config        Config
	Authenticated func(http.Handler) http.Handler
}

func Main(defaultProductID string, defaultPort string) {
	config := LoadConfig(defaultProductID, defaultPort)
	if err := ListenAndServe(context.Background(), config); err != nil {
		log.Fatal(err)
	}
}

func LoadConfig(defaultProductID string, defaultPort string) Config {
	return Config{
		Port:           env("PORT", defaultPort),
		ProductID:      env("PRODUCT_ID", defaultProductID),
		JWKSURL:        env("AUTH_JWKS_URL", "http://localhost:8080/.well-known/jwks.json"),
		Issuer:         env("JWT_ISSUER", "http://auth.codelinks.localhost"),
		Audience:       env("JWT_AUDIENCE", "codelinks-products"),
		AllowedOrigins: splitCSV(os.Getenv("PRODUCT_ALLOWED_ORIGINS")),
	}
}

func ListenAndServe(ctx context.Context, config Config) error {
	validator, err := connectValidator(ctx, config)
	if err != nil {
		return err
	}

	log.Printf("%s backend listening on :%s", config.ProductID, config.Port)
	return http.ListenAndServe(":"+config.Port, Handler(config, validator))
}

func Handler(config Config, validator *productauth.RemoteValidator, registrars ...RouteRegistrar) http.Handler {
	mux := http.NewServeMux()
	tools := RouteTools{
		Config: config,
		Authenticated: func(next http.Handler) http.Handler {
			if validator == nil {
				return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth validator unavailable"})
				})
			}
			return validator.Middleware(next)
		},
	}

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "productId": config.ProductID})
	})
	mux.Handle("GET /api/me", tools.Authenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := productauth.CurrentUserFromContext(r.Context())
		writeJSON(w, http.StatusOK, map[string]any{
			"userId":        user.ID,
			"email":         user.Email,
			"name":          user.Name,
			"status":        user.Status,
			"emailVerified": user.EmailVerified,
			"licenses":      user.Licenses,
			"productId":     config.ProductID,
		})
	})))
	for _, registrar := range registrars {
		registrar(mux, tools)
	}

	return withCORS(mux, config.AllowedOrigins)
}

func connectValidator(ctx context.Context, config Config) (*productauth.RemoteValidator, error) {
	var lastErr error
	for attempt := 0; attempt < 10; attempt++ {
		validator, err := productauth.NewRemoteValidator(ctx, productauth.RemoteValidatorConfig{
			JWKSURL:   config.JWKSURL,
			Issuer:    config.Issuer,
			Audience:  config.Audience,
			ProductID: config.ProductID,
		})
		if err == nil {
			return validator, nil
		}
		lastErr = err
		time.Sleep(time.Second)
	}

	return nil, lastErr
}

func withCORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if origin != "" {
			allowed[origin] = true
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowed[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Headers", "authorization, content-type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}
