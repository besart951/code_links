package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	config := loadConfig()

	signer, err := newTokenSigner(config)
	if err != nil {
		log.Fatal(err)
	}

	store, cleanup, err := openStore(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	server := newServer(config, store, signer)
	httpServer := &http.Server{
		Addr:              ":" + config.Port,
		Handler:           server.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("auth-service listening on :%s", config.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
