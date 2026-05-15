package identity

import (
	"net/http"

	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	svc *Service
}

func NewHandlers(svc *Service) *Handlers {
	return &Handlers{svc: svc}
}

// Submit godoc
// @Summary      Submit identity verification
// @Description  Citizen submits national ID and document paths. Re-submit resets status to pending.
// @Tags         identity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      SubmitInput  true  "Verification payload"
// @Success      200   {object}  IdentityVerification
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /identity/submit [post]
func (h *Handlers) Submit(c *gin.Context) {
	var in SubmitInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := middleware.GetAuthUser(c)
	out, err := h.svc.Submit(c.Request.Context(), user.ID, in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// MyStatus godoc
// @Summary      Get my identity verification status
// @Description  Returns the citizen's verification record, or user_id with empty status if never submitted.
// @Tags         identity
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  IdentityVerification
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /identity/status [get]
func (h *Handlers) MyStatus(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	out, err := h.svc.MyStatus(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// DecisionInput is an optional admin comment on approve/reject.
type DecisionInput struct {
	Comment string `json:"comment" example:"Documents verified"`
}

// Approve godoc
// @Summary      Approve identity verification
// @Description  Admin approves a verification by ID.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string         true  "Verification MongoDB ObjectID (hex)"
// @Param        body  body      DecisionInput  false  "Optional review comment"
// @Success      200   {object}  IdentityVerification
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /admin/identity/{id}/approve [put]
func (h *Handlers) Approve(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body DecisionInput
	_ = c.ShouldBindJSON(&body)
	admin := middleware.GetAuthUser(c)
	out, err := h.svc.Approve(c.Request.Context(), id, admin.ID, body.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

// Reject godoc
// @Summary      Reject identity verification
// @Description  Admin rejects a verification by ID.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string         true  "Verification MongoDB ObjectID (hex)"
// @Param        body  body      DecisionInput  false  "Optional review comment"
// @Success      200   {object}  IdentityVerification
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      403   {object}  map[string]string
// @Router       /admin/identity/{id}/reject [put]
func (h *Handlers) Reject(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body DecisionInput
	_ = c.ShouldBindJSON(&body)
	admin := middleware.GetAuthUser(c)
	out, err := h.svc.Reject(c.Request.Context(), id, admin.ID, body.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
