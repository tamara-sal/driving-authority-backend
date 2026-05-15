package http

import (
	"net/http"

	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

// Health godoc
// @Summary      Health check
// @Description  Returns OK when the API is running.
// @Tags         system
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{OK: true})
}

// Me godoc
// @Summary      Current user profile
// @Description  Returns the authenticated user's id, email, and role.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  MeResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /me [get]
func Me(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	c.JSON(http.StatusOK, MeResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Role:  string(user.Role),
	})
}

// AdminPing godoc
// @Summary      Admin RBAC smoke test
// @Description  Example admin-only endpoint; returns admin:true when the caller has the admin role.
// @Tags         admin
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  AdminPingResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Router       /admin/ping [get]
func AdminPing(c *gin.Context) {
	c.JSON(http.StatusOK, AdminPingResponse{Admin: true})
}
