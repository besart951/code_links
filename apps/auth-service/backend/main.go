package main

import (
	"context"
	"log"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/app"
	"github.com/besart951/code-links/apps/auth-service/backend/internal/processlog"
)

func main() {
	closeLog, err := processlog.ConfigureFromEnv()
	if err != nil {
		log.Fatalf("fatal: configure process log: %v", err)
	}
	defer closeLog()

	if err := app.Run(context.Background()); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}
