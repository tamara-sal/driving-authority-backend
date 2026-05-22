package platform

import (
	"net/http"

	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) ListUsers(c *gin.Context) {
	out, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) ListApplications(c *gin.Context) {
	out, err := h.svc.ListApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) ListActivity(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	out, err := h.svc.ListActivity(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) ListAuditLogs(c *gin.Context) {
	out, err := h.svc.ListAuditLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
