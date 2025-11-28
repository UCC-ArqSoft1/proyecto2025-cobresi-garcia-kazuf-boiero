package handlers

import "github.com/gin-gonic/gin"

// HealthHandler exposes endpoints for readiness/liveness.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router *gin.Engine) {
    router.GET("/api/health", h.GetStatus)
}

func (h *HealthHandler) GetStatus(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
}
