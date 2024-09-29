package server

import (
	"fmt"
	"net"

	"net/http"

	"github.com/luiidev/go/config"
	"github.com/luiidev/go/pkg/logger"
	"gorm.io/gorm"
)

// Run creates objects via constructors.
func Run(cfg *config.Config, l *logger.Logger, db *gorm.DB) {
	router := Router(l, db, cfg)

	httpServer := http.Server{
		Addr:    net.JoinHostPort("", cfg.HTTP.Port),
		Handler: router,
	}

	fmt.Printf("Server is running on port %v\n", cfg.HTTP.Port)
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
