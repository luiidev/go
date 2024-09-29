package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/models"
	"github.com/luiidev/go/pkg/logger"
	"github.com/luiidev/go/pkg/validation"
	"gorm.io/gorm"
)

var jwtKey = []byte("my_secret_key")

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
	token, err := c.createToken(user.ID)
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
func (c AuthController) createToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.cfg.App.JwtSecret))
}

// Middleware para proteger rutas (verifica el JWT)
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Ruta protegida
func ProtectedRoute(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "You are authenticated"})
}
