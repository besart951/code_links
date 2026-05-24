package main

import (
	"context"
	"log"

	"github.com/besart951/code-links/apps/auth-service/backend/internal/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
