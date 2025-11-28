package handlers

import (
	"net/http"

	"github.com/alesio/gestion-actividades-deportivas/models"
	"github.com/alesio/gestion-actividades-deportivas/security"
	"github.com/alesio/gestion-actividades-deportivas/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{authService: authService, userService: userService}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/auth/login", h.Login)
	router.POST("/auth/register", h.Register)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Payload inválido", "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.authService.Authenticate(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			respondError(c, http.StatusUnauthorized, "Credenciales inválidas", "UNAUTHORIZED", "")
			return
		}
		respondError(c, http.StatusInternalServerError, "No se pudo autenticar", "INTERNAL_ERROR", err.Error())
		return
	}

	token, err := h.authService.GenerateJWT(user)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "No se pudo generar el token", "INTERNAL_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "Login exitoso",
		Data: gin.H{
			"token": token,
			"user":  toUserResponse(user),
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Payload inválido", "VALIDATION_ERROR", err.Error())
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Email, req.Password, "socio")
	if err != nil {
		switch err {
		case services.ErrEmailAlreadyExists:
			respondError(c, http.StatusConflict, "El email ya está registrado", "VALIDATION_ERROR", "")
		case security.ErrEmptyPassword:
			respondError(c, http.StatusBadRequest, "La contraseña es obligatoria", "VALIDATION_ERROR", "")
		default:
			respondError(c, http.StatusInternalServerError, "No se pudo crear el usuario", "INTERNAL_ERROR", err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: "Registro exitoso",
		Data:    toUserResponse(user),
	})
}

func toUserResponse(user *models.User) userResponse {
	return userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
}

func respondError(c *gin.Context, status int, message, code, details string) {
	c.JSON(status, APIError{
		Success: false,
		Error:   message,
		Code:    code,
		Details: details,
	})
}
