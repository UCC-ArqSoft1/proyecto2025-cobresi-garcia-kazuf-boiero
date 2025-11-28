package middlewares

import (
	"net/http"

	"github.com/alesio/gestion-actividades-deportivas/handlers"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware guards endpoints so only admin users may access them.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.APIError{
				Success: false,
				Error:   "Contexto de rol faltante",
				Code:    "FORBIDDEN",
			})
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, handlers.APIError{
				Success: false,
				Error:   "Acceso restringido a administradores",
				Code:    "FORBIDDEN",
			})
			return
		}

		c.Next()
	}
}
