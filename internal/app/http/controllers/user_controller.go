package controllers

import (
	"errors"
	"net/http"

	"github.com/luiidev/go/internal/app/models"
	"github.com/luiidev/go/pkg/logger"
	res "github.com/luiidev/go/pkg/response"
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

	res.JSON(w, res.H{
		"message": "Users",
		"data":    users,
	})
}

func (c UserController) Show(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")

	user := models.User{}

	result := c.db.First(&user, userId)
	if result.Error != nil {
		res.JSON(w, res.H{"message": "Usuario no encontrado"}, http.StatusNotFound)
		return
	}

	res.JSON(w, res.H{"message": "User", "data": user})
}

func (c UserController) Store(w http.ResponseWriter, r *http.Request) {
	var user models.User

	validator := validation.Make(r.Body, &user)
	if validator.Fails() {
		validator.Response(w)
		return
	}

	result := c.db.Create(&user)

	if result.Error != nil {
		c.l.Error(result.Error)

		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			res.JSON(w, res.H{"message": "El email ya esta en uso"}, http.StatusUnprocessableEntity)
			return
		}

		res.JSON(w, res.H{"message": "Ocurrio un error"}, http.StatusInternalServerError)
		return
	}

	res.JSON(w, res.H{
		"message": "User created",
		"data":    user,
	})
}
