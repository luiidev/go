package server

import (
	"net/http"

	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/http/controllers"
	"github.com/luiidev/go/internal/app/http/middleware"
	"github.com/luiidev/go/pkg/logger"
	"gorm.io/gorm"
)

func Router(l *logger.Logger, db *gorm.DB, cfg *config.Config) *http.ServeMux {
	router := http.NewServeMux()
	authMiddleware := &middleware.AuthMiddleware{Cfg: *cfg, Db: *db}

	exampleController := controllers.NewExampleController(*l, *db)
	router.HandleFunc("GET /helloworld", exampleController.Helloworld)

	userController := controllers.NewUserController(*l, *db)
	router.HandleFunc("GET /users", authMiddleware.Handle(userController.Index))
	router.HandleFunc("GET /users/{id}", authMiddleware.Handle(userController.Show))
	router.HandleFunc("POST /users", authMiddleware.Handle(userController.Store))

	authController := controllers.NewAuthController(*l, *db, *cfg)
	router.HandleFunc("POST /login", authController.Login)
	router.HandleFunc("POST /register", authController.Register)
	router.HandleFunc("GET /me", authMiddleware.Handle(authController.Me))
	router.HandleFunc("POST /me", authMiddleware.Handle(authController.MeUpdate))

	return router
}
