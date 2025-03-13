package main

import (
	"log"

	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("main - main - Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
