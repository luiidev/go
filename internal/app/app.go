package app

import (
	"fmt"

	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/http/server"
	"github.com/luiidev/go/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	db, err := gorm.Open(postgres.Open(cfg.DB.URL), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres: %w", err))
	}

	server.Run(cfg, l, db)
}
