package handlers

import (
    "net/http"
    "strconv"

    "github.com/alesio/gestion-actividades-deportivas/services"
    "github.com/gin-gonic/gin"
)

// EnrollmentsHandler handles actions for enrolled members.
type EnrollmentsHandler struct {
    enrollmentService *services.EnrollmentService
}

func NewEnrollmentsHandler(enrollmentService *services.EnrollmentService) *EnrollmentsHandler {
    return &EnrollmentsHandler{enrollmentService: enrollmentService}
}

func (h *EnrollmentsHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/activities/:id/enroll", h.EnrollInActivity)
    router.GET("/me/activities", h.ListMyActivities)
}

func (h *EnrollmentsHandler) EnrollInActivity(c *gin.Context) {
    userID, ok := c.Get("userID")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
        return
    }
    userIDUint, ok := userID.(uint)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
        return
    }

    activityIDParam := c.Param("id")
    activityID, err := strconv.Atoi(activityIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
        return
    }

    enrollment, err := h.enrollmentService.EnrollUserInActivity(userIDUint, uint(activityID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, enrollment)
}

func (h *EnrollmentsHandler) ListMyActivities(c *gin.Context) {
    userID, ok := c.Get("userID")
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user context"})
        return
    }
    userIDUint, ok := userID.(uint)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
        return
    }

    enrollments, err := h.enrollmentService.GetEnrollmentsByUser(userIDUint)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch enrollments"})
        return
    }

    c.JSON(http.StatusOK, enrollments)
}
