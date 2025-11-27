package middlewares

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// AdminMiddleware guards endpoints so only admin users may access them.
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        roleValue, exists := c.Get("role")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing role"})
            return
        }

        role, ok := roleValue.(string)
        if !ok || role != "admin" {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
            return
        }

        c.Next()
    }
}
