package controllers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/models"
	"github.com/luiidev/go/pkg/logger"
	"github.com/luiidev/go/pkg/validation"
	"gorm.io/gorm"
)

type Credentials struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

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
	var creds Credentials
	errors, err := validation.DecodeAndValidate(r.Body, &creds)
	if err != nil {
		JsonResponse(w, Response{Message: err.Error(), Errors: errors}, http.StatusUnprocessableEntity)
		return
	}

	// Buscar el usuario en la "base de datos"
	var user models.User
	result := c.db.First(&user, "email = ?", creds.Email)
	if result.Error != nil {
		invalidCredentials(w)
		return
	}

	if user.CheckPasswordHash(creds.Password) {
		invalidCredentials(w)
		return
	}

	// Crear el token JWT
	token, err := c.createToken(user)
	if err != nil {
		JsonResponse(w, Response{Message: "Error generating token"}, http.StatusInternalServerError)
		return
	}

	// Responder con el token
	JsonResponse(w, Response{Message: "Login successful", Data: JWTResponse{Token: token, User: user}})
}

func invalidCredentials(w http.ResponseWriter) {
	JsonResponse(w, Response{Message: "Usuario y/o contraseña incorrecta"}, http.StatusUnauthorized)
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
