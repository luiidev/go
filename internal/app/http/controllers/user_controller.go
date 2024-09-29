package controllers

import (
	"net/http"

	"github.com/luiidev/go/internal/app/models"
	"github.com/luiidev/go/pkg/logger"
	"github.com/luiidev/go/pkg/validation"
	"gorm.io/gorm"
)

type UserController struct {
	l  logger.Logger
	db gorm.DB
}

func NewUserController(l logger.Logger, db gorm.DB) *UserController {
	return &UserController{l: l, db: db}
}

func (c UserController) Index(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}
	c.db.Limit(10).Find(&users)

	JsonResponse(w, Response{Message: "Users", Data: users})
}

func (c UserController) Show(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")

	user := models.User{}

	result := c.db.First(&user, userId)
	if result.Error != nil {
		JsonResponse(w, Response{Message: "Usuario no encontrado"}, http.StatusNotFound)
		return
	}

	JsonResponse(w, Response{Message: "User", Data: user})
}

func (c UserController) Store(w http.ResponseWriter, r *http.Request) {
	var user models.User

	errors, err := validation.DecodeAndValidate(r.Body, &user)
	if err != nil {
		JsonResponse(w, Response{Message: err.Error(), Errors: errors}, http.StatusUnprocessableEntity)
		return
	}

	result := c.db.Create(&user)

	if result.Error != nil {
		c.l.Error(result.Error)
		JsonResponse(w, Response{Message: "Ocurrio un error"}, http.StatusInternalServerError)
		return
	}

	JsonResponse(w, Response{Message: "User created", Data: user})
}
