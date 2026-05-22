package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/besart951/code-links/packages/productauth"
)

const defaultProductID = "loka-link"
const defaultPort = "8093"

func main() {
	config := loadConfig()
	validator, err := connectValidator(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "productId": config.ProductID})
	})
	mux.Handle("GET /api/me", validator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, _ := productauth.ClaimsFromContext(r.Context())
		writeJSON(w, http.StatusOK, map[string]any{
			"userId":    claims.Subject,
			"email":     claims.Email,
			"name":      claims.Name,
			"productId": config.ProductID,
		})
	})))

	log.Printf("%s backend listening on :%s", config.ProductID, config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, withCORS(mux)))
}

type appConfig struct {
	Port      string
	ProductID string
	JWKSURL   string
	Issuer    string
	Audience  string
}

func loadConfig() appConfig {
	return appConfig{
		Port:      env("PORT", defaultPort),
		ProductID: env("PRODUCT_ID", defaultProductID),
		JWKSURL:   env("AUTH_JWKS_URL", "http://localhost:8080/.well-known/jwks.json"),
		Issuer:    env("JWT_ISSUER", "http://auth.codelinks.localhost"),
		Audience:  env("JWT_AUDIENCE", "codelinks-products"),
	}
}

func connectValidator(ctx context.Context, config appConfig) (*productauth.RemoteValidator, error) {
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

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
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
