package server

import (
	"net/http"

	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/http/controllers"
	"github.com/luiidev/go/pkg/logger"
	"gorm.io/gorm"
)

func Router(l *logger.Logger, db *gorm.DB, cfg *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	exampleController := controllers.NewExampleController(*l, *db)
	router.HandleFunc("GET /helloworld", exampleController.Helloworld)

	userController := controllers.NewUserController(*l, *db)
	router.HandleFunc("GET /users", userController.Index)
	router.HandleFunc("GET /users/{id}", userController.Show)
	router.HandleFunc("POST /users", userController.Store)

	authController := controllers.NewAuthController(*l, *db, *cfg)
	router.HandleFunc("POST /login", authController.Login)

	return router
}
