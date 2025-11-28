package handlers

import (
    "net/http"
    "strconv"

    "github.com/alesio/gestion-actividades-deportivas/services"
    "github.com/gin-gonic/gin"
)

// ActivitiesHandler serves public activities endpoints.
type ActivitiesHandler struct {
    activityService *services.ActivityService
}

func NewActivitiesHandler(activityService *services.ActivityService) *ActivitiesHandler {
    return &ActivitiesHandler{activityService: activityService}
}

func (h *ActivitiesHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.GET("/activities", h.ListActivities)
    router.GET("/activities/:id", h.GetActivity)
}

func (h *ActivitiesHandler) ListActivities(c *gin.Context) {
    var filter services.ActivityFilter
    filter.Query = c.Query("q")
    filter.Category = c.Query("category")

    if dayStr := c.Query("day"); dayStr != "" {
        if dayInt, err := strconv.Atoi(dayStr); err == nil {
            filter.Day = &dayInt
        }
    }

    activities, err := h.activityService.ListActivities(filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list activities"})
        return
    }

    c.JSON(http.StatusOK, activities)
}

func (h *ActivitiesHandler) GetActivity(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
        return
    }

    activity, err := h.activityService.GetActivityByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
        return
    }

    c.JSON(http.StatusOK, activity)
}
