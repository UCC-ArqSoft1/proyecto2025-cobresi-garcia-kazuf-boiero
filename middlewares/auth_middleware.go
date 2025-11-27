package middlewares

import (
    "net/http"
    "strings"

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
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            return
        }

        token := strings.TrimPrefix(header, "Bearer ")
        claims, err := m.authService.ValidateJWT(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        // TODO: Ensure token validation returns expiration errors and user existence checks.

        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
