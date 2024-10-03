package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/luiidev/go/config"
	"github.com/luiidev/go/internal/app/http/controllers"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	Db  gorm.DB
	Cfg config.Config
}

// Middleware para proteger rutas (verifica el JWT)
func (m AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el valor del encabezado Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			controllers.JsonResponse(w, controllers.Response{Message: "Authorization header missing"}, http.StatusUnauthorized)
			return
		}

		// Verificar que el encabezado empiece con "Bearer "
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			controllers.JsonResponse(w, controllers.Response{Message: "Invalid token prefix, expected 'Bearer'"}, http.StatusUnauthorized)
			return
		}

		// Extraer el token eliminando el prefijo "Bearer "
		tokenStr := strings.TrimPrefix(authHeader, bearerPrefix)

		// Parsear el token y validarlo
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return []byte(m.Cfg.JWT.Secret), nil
		})

		// Validar si hubo un error durante el parsing o el token es inválido
		if err != nil || !token.Valid {
			controllers.JsonResponse(w, controllers.Response{Message: "Invalid token"}, http.StatusUnauthorized)
			return
		}

		// Extraer las reclamaciones (claims) del token JWT
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Aquí se asume que el token tiene una reclamación "id"
			sub := claims["sub"]

			// Agregar el user al contexto de la solicitud
			ctx := context.WithValue(r.Context(), "sub", sub)

			// Pasar la solicitud con el nuevo contexto al siguiente manejador
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			controllers.JsonResponse(w, controllers.Response{Message: "Invalid token claims"}, http.StatusUnauthorized)
		}
	}
}
