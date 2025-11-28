package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alesio/gestion-actividades-deportivas/handlers"
	"github.com/alesio/gestion-actividades-deportivas/services"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates tokens and hydrates the context with user data.
type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handlers.APIError{
				Success: false,
				Error:   "Token faltante",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := m.authService.ValidateJWT(token)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Token inválido"
			code := "UNAUTHORIZED"

			switch {
			case errors.Is(err, services.ErrTokenExpired):
				message = "Token expirado"
			case errors.Is(err, services.ErrInvalidToken):
				// keep defaults
			default:
				status = http.StatusInternalServerError
				message = "No se pudo validar el token"
				code = "INTERNAL_ERROR"
			}

			c.AbortWithStatusJSON(status, handlers.APIError{
				Success: false,
				Error:   message,
				Code:    code,
				Details: err.Error(),
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
