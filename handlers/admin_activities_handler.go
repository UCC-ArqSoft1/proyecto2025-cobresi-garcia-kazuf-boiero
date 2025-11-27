package handlers

import (
    "net/http"
    "strconv"

    "github.com/alesio/gestion-actividades-deportivas/models"
    "github.com/alesio/gestion-actividades-deportivas/services"
    "github.com/gin-gonic/gin"
)

// AdminActivitiesHandler exposes admin-only endpoints for managing activities.
type AdminActivitiesHandler struct {
    activityService *services.ActivityService
}

func NewAdminActivitiesHandler(activityService *services.ActivityService) *AdminActivitiesHandler {
    return &AdminActivitiesHandler{activityService: activityService}
}

func (h *AdminActivitiesHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/admin/activities", h.CreateActivity)
    router.PUT("/admin/activities/:id", h.UpdateActivity)
    router.DELETE("/admin/activities/:id", h.DeleteActivity)
}

func (h *AdminActivitiesHandler) CreateActivity(c *gin.Context) {
    var activity models.Activity
    if err := c.ShouldBindJSON(&activity); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.activityService.CreateActivity(&activity); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create activity"})
        return
    }

    c.JSON(http.StatusCreated, activity)
}

func (h *AdminActivitiesHandler) UpdateActivity(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
        return
    }

    var payload models.Activity
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    payload.ID = uint(id)

    if err := h.activityService.UpdateActivity(&payload); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update activity"})
        return
    }

    c.JSON(http.StatusOK, payload)
}

func (h *AdminActivitiesHandler) DeleteActivity(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
        return
    }

    if err := h.activityService.DeleteActivity(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete activity"})
        return
    }

    c.Status(http.StatusNoContent)
}
