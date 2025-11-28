package handlers

import (
    "errors"
    "net/http"
    "strconv"
    "time"

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
    router.GET("/admin/activities", h.ListActivities)
    router.POST("/admin/activities", h.CreateActivity)
    router.PUT("/admin/activities/:id", h.UpdateActivity)
    router.DELETE("/admin/activities/:id", h.DeleteActivity)
}

type activityRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
    Category    string `json:"category" binding:"required"`
    DayOfWeek   int    `json:"day_of_week" binding:"required"`
    StartTime   string `json:"start_time" binding:"required"`
    EndTime     string `json:"end_time" binding:"required"`
    Capacity    int    `json:"capacity" binding:"required"`
    Instructor  string `json:"instructor" binding:"required"`
    ImageURL    string `json:"image_url"`
    IsActive    *bool  `json:"is_active"`
}

func (h *AdminActivitiesHandler) ListActivities(c *gin.Context) {
    var filter services.ActivityFilter
    filter.Query = c.Query("q")
    filter.Category = c.Query("category")

    if dayStr := c.Query("day"); dayStr != "" {
        day, err := strconv.Atoi(dayStr)
        if err != nil {
            respondError(c, http.StatusBadRequest, "day debe ser numerico", "VALIDATION_ERROR", "")
            return
        }
        if day < 0 || day > 6 {
            respondError(c, http.StatusBadRequest, "day debe estar entre 0 y 6", "VALIDATION_ERROR", "")
            return
        }
        filter.Day = &day
    }

    var isActiveFilter *bool
    if isActiveStr := c.Query("is_active"); isActiveStr != "" {
        value, err := strconv.ParseBool(isActiveStr)
        if err != nil {
            respondError(c, http.StatusBadRequest, "is_active debe ser booleano", "VALIDATION_ERROR", "")
            return
        }
        isActiveFilter = &value
    }

    activities, err := h.activityService.ListActivitiesAdmin(services.AdminActivityFilter{
        ActivityFilter: filter,
        IsActive:       isActiveFilter,
    })
    if err != nil {
        respondError(c, http.StatusInternalServerError, "No se pudieron listar las actividades", "INTERNAL_ERROR", err.Error())
        return
    }

    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Data:    activities,
    })
}

func (h *AdminActivitiesHandler) CreateActivity(c *gin.Context) {
    var req activityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, "Payload inválido", "VALIDATION_ERROR", err.Error())
        return
    }

    if err := validateActivityRequest(req); err != nil {
        respondError(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR", "")
        return
    }

    activity := models.Activity{
        Title:       req.Title,
        Description: req.Description,
        Category:    req.Category,
        DayOfWeek:   req.DayOfWeek,
        StartTime:   req.StartTime,
        EndTime:     req.EndTime,
        Capacity:    req.Capacity,
        Instructor:  req.Instructor,
        ImageURL:    req.ImageURL,
        IsActive:    true,
    }
    if req.IsActive != nil {
        activity.IsActive = *req.IsActive
    }

    if err := h.activityService.CreateActivity(&activity); err != nil {
        respondError(c, http.StatusInternalServerError, "No se pudo crear la actividad", "INTERNAL_ERROR", err.Error())
        return
    }

    c.JSON(http.StatusCreated, APIResponse{
        Success: true,
        Message: "Actividad creada",
        Data:    activity,
    })
}

func (h *AdminActivitiesHandler) UpdateActivity(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        respondError(c, http.StatusBadRequest, "ID de actividad invalido", "VALIDATION_ERROR", "")
        return
    }

    var req activityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondError(c, http.StatusBadRequest, "Payload inválido", "VALIDATION_ERROR", err.Error())
        return
    }

    if err := validateActivityRequest(req); err != nil {
        respondError(c, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR", "")
        return
    }

    activity, err := h.activityService.GetActivityByID(uint(id))
    if err != nil {
        if errors.Is(err, services.ErrActivityNotFound) {
            respondError(c, http.StatusNotFound, "Actividad no encontrada", "NOT_FOUND", "")
            return
        }
        respondError(c, http.StatusInternalServerError, "No se pudo obtener la actividad", "INTERNAL_ERROR", err.Error())
        return
    }

    activity.Title = req.Title
    activity.Description = req.Description
    activity.Category = req.Category
    activity.DayOfWeek = req.DayOfWeek
    activity.StartTime = req.StartTime
    activity.EndTime = req.EndTime
    activity.Capacity = req.Capacity
    activity.Instructor = req.Instructor
    activity.ImageURL = req.ImageURL
    if req.IsActive != nil {
        activity.IsActive = *req.IsActive
    }

    if err := h.activityService.UpdateActivity(activity); err != nil {
        respondError(c, http.StatusInternalServerError, "No se pudo actualizar la actividad", "INTERNAL_ERROR", err.Error())
        return
    }

    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Message: "Actividad actualizada",
        Data:    activity,
    })
}

func (h *AdminActivitiesHandler) DeleteActivity(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        respondError(c, http.StatusBadRequest, "ID de actividad invalido", "VALIDATION_ERROR", "")
        return
    }

    if err := h.activityService.DeleteActivity(uint(id)); err != nil {
        if errors.Is(err, services.ErrActivityNotFound) {
            respondError(c, http.StatusNotFound, "Actividad no encontrada", "NOT_FOUND", "")
            return
        }
        respondError(c, http.StatusInternalServerError, "No se pudo desactivar la actividad", "INTERNAL_ERROR", err.Error())
        return
    }

    c.JSON(http.StatusOK, APIResponse{
        Success: true,
        Message: "Actividad desactivada",
    })
}

func validateActivityRequest(req activityRequest) error {
    if req.DayOfWeek < 0 || req.DayOfWeek > 6 {
        return errors.New("day_of_week debe estar entre 0 y 6")
    }

    if req.Capacity <= 0 {
        return errors.New("capacity debe ser mayor a 0")
    }

    start, err := time.Parse("15:04", req.StartTime)
    if err != nil {
        return errors.New("start_time debe tener formato HH:MM")
    }

    end, err := time.Parse("15:04", req.EndTime)
    if err != nil {
        return errors.New("end_time debe tener formato HH:MM")
    }

    if !end.After(start) {
        return errors.New("end_time debe ser mayor a start_time")
    }

    return nil
}
