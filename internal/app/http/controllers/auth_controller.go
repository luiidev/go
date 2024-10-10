package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/http/requests"
	"github.com/luiidev/go/internal/app/models"
	"github.com/luiidev/go/pkg/logger"
	res "github.com/luiidev/go/pkg/response"
	"github.com/luiidev/go/pkg/validation"
	"gorm.io/gorm"
)

// Estructura de la respuesta del JWT
type JWTResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type AuthController struct {
	l   logger.Logger
	db  gorm.DB
	cfg config.Config
}

func NewAuthController(l logger.Logger, db gorm.DB, cfg config.Config) *AuthController {
	return &AuthController{l: l, db: db, cfg: cfg}
}

// Login: Validar las credenciales y generar el token JWT
func (c AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds requests.LoginRequest

	validator := validation.Make(r.Body, &creds)
	if validator.Fails() {
		validator.Response(w)
		return
	}

	// Buscar el usuario en la "base de datos"
	var user models.User
	result := c.db.First(&user, "email = ?", creds.Email)
	if result.Error != nil {
		invalidCredentials(w)
		return
	}

	if !user.CheckPasswordHash(creds.Password) {
		invalidCredentials(w)
		return
	}

	// Crear el token JWT
	token, err := c.createToken(user)
	if err != nil {
		res.JSON(w, res.H{"message": "Error generating token"}, http.StatusInternalServerError)
		return
	}

	// Responder con el token
	res.JSON(w, res.H{
		"message": "Login successful",
		"data": JWTResponse{
			Token: token,
			User:  user,
		},
	})
}

func (c AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var userRequest requests.StoreUserRequest

	validator := validation.Make(r.Body, &userRequest)
	if validator.Fails() {
		validator.Response(w)
		return
	}

	user := models.User{
		FirstName: userRequest.FirstName,
		LastName:  userRequest.LastName,
		Email:     userRequest.Email,
		Password:  userRequest.Password,
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

func (c AuthController) Me(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("sub").(float64)

	if !ok {
		res.JSON(w, res.H{"message": "User not found"}, http.StatusNotFound)
		return
	}

	var user models.User

	result := c.db.First(&user, userId)
	if result.Error != nil {
		res.JSON(w, res.H{"message": "User not found"}, http.StatusNotFound)
		return
	}

	res.JSON(w, res.H{
		"message": "User",
		"data":    user,
	})
}

func (c AuthController) MeUpdate(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("sub").(float64)

	if !ok {
		res.JSON(w, res.H{"message": "User not found"}, http.StatusNotFound)
		return
	}

	var userRequest requests.UpdateUserRequest
	validator := validation.Make(r.Body, &userRequest)
	if validator.Fails() {
		validator.Response(w)
		return
	}

	result := c.db.Model(models.User{}).
		Where("id = ?", userId).
		UpdateColumns(map[string]interface{}{
			"first_name": userRequest.FirstName,
			"last_name":  userRequest.LastName,
		})

	if result.Error != nil {
		res.JSON(w, res.H{"message": "Ocurrio un error"}, http.StatusInternalServerError)
		return
	}

	var user models.User

	result = c.db.First(&user, userId)
	if result.Error != nil {
		res.JSON(w, res.H{"message": "User not found"}, http.StatusNotFound)
		return
	}

	res.JSON(w, res.H{
		"message": "User updated",
		"data":    user,
	})
}

func invalidCredentials(w http.ResponseWriter) {
	res.JSON(w, res.H{"message": "Usuario y/o contraseña incorrecta"}, http.StatusUnauthorized)
}

// Crear el token JWT
func (c AuthController) createToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Duration(c.cfg.JWT.Expiration) * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.cfg.JWT.Secret))
}
