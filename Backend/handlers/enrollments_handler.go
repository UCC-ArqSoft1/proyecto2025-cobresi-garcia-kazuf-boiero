package handlers

import (
	"net/http"
	"strconv"

	"github.com/alesio/gestion-actividades-deportivas/services"
	"github.com/gin-gonic/gin"
)

// EnrollmentsHandler handles actions for enrolled members.
type EnrollmentsHandler struct {
	enrollmentService services.EnrollmentService
}

type myActivityDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	DayOfWeek   int    `json:"day_of_week"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	Instructor  string `json:"instructor"`
}

func NewEnrollmentsHandler(enrollmentService services.EnrollmentService) *EnrollmentsHandler {
	return &EnrollmentsHandler{enrollmentService: enrollmentService}
}

func (h *EnrollmentsHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/activities/:id/enroll", h.EnrollInActivity)
	router.DELETE("/activities/:id/enroll", h.UnenrollFromActivity)
	router.GET("/me/activities", h.ListMyActivities)
}

func (h *EnrollmentsHandler) EnrollInActivity(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	activityIDParam := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Success: false,
			Error:   "ID de actividad invalido",
			Code:    "VALIDATION_ERROR",
		})
		return
	}

	enrollment, err := h.enrollmentService.EnrollUserInActivity(userID, uint(activityID))
	if err != nil {
		switch err {
		case services.ErrActivityNotFound:
			c.JSON(http.StatusNotFound, APIError{
				Success: false,
				Error:   "Actividad no encontrada",
				Code:    "ACTIVITY_NOT_FOUND",
			})
		case services.ErrActivityInactive:
			c.JSON(http.StatusBadRequest, APIError{
				Success: false,
				Error:   "La actividad no esta activa",
				Code:    "ACTIVITY_INACTIVE",
			})
		case services.ErrAlreadyEnrolled:
			c.JSON(http.StatusConflict, APIError{
				Success: false,
				Error:   "Ya estas inscripto en esta actividad",
				Code:    "ALREADY_ENROLLED",
			})
		case services.ErrNoCapacity:
			c.JSON(http.StatusConflict, APIError{
				Success: false,
				Error:   "La actividad no tiene cupos disponibles",
				Code:    "NO_CAPACITY",
			})
		case services.ErrScheduleConflict:
			c.JSON(http.StatusConflict, APIError{
				Success: false,
				Error:   "La actividad se solapa en dia y horario con otra inscripcion activa",
				Code:    "SCHEDULE_CONFLICT",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIError{
				Success: false,
				Error:   "No se pudo completar la inscripcion",
				Code:    "INTERNAL_ERROR",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Inscripcion exitosa",
		Data:    enrollment,
	})
}

func (h *EnrollmentsHandler) ListMyActivities(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	enrollments, err := h.enrollmentService.GetUserEnrollments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIError{
			Success: false,
			Error:   "No se pudieron obtener las actividades",
			Code:    "INTERNAL_ERROR",
			Details: err.Error(),
		})
		return
	}

	activities := make([]myActivityDTO, 0, len(enrollments))
	for _, enrollment := range enrollments {
		activity := enrollment.Activity
		activities = append(activities, myActivityDTO{
			ID:          activity.ID,
			Title:       activity.Title,
			Description: activity.Description,
			Category:    activity.Category,
			DayOfWeek:   activity.DayOfWeek,
			StartTime:   activity.StartTime,
			EndTime:     activity.EndTime,
			Instructor:  activity.Instructor,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    activities,
	})
}

func (h *EnrollmentsHandler) UnenrollFromActivity(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	activityIDParam := c.Param("id")
	activityID, err := strconv.ParseUint(activityIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{
			Success: false,
			Error:   "ID de actividad invalido",
			Code:    "VALIDATION_ERROR",
		})
		return
	}

	if err := h.enrollmentService.UnenrollUserFromActivity(userID, uint(activityID)); err != nil {
		switch err {
		case services.ErrEnrollmentNotFound:
			c.JSON(http.StatusNotFound, APIError{
				Success: false,
				Error:   "No estabas inscripto en esta actividad",
				Code:    "ENROLLMENT_NOT_FOUND",
			})
		default:
			c.JSON(http.StatusInternalServerError, APIError{
				Success: false,
				Error:   "No pudimos procesar la baja",
				Code:    "INTERNAL_ERROR",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Te desinscribiste de la actividad",
	})
}

func getUserIDFromContext(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIError{
			Success: false,
			Error:   "Token faltante",
			Code:    "UNAUTHORIZED",
		})
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, APIError{
			Success: false,
			Error:   "Contexto de usuario invalido",
			Code:    "INTERNAL_ERROR",
		})
		return 0, false
	}

	return userID, true
}
