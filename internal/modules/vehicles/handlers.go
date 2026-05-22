package vehicles

import (
	"net/http"
	"strings"

	"driving-authority-backend/internal/http/middleware"
	"driving-authority-backend/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	svc   *Service
	users *auth.UserRepo
}

func NewHandlers(svc *Service, users *auth.UserRepo) *Handlers {
	return &Handlers{svc: svc, users: users}
}

func (h *Handlers) Create(c *gin.Context) {
	var in CreateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := middleware.GetAuthUser(c)
	out, err := h.svc.Create(c.Request.Context(), user.ID, in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *Handlers) MyVehicles(c *gin.Context) {
	user := middleware.GetAuthUser(c)
	out, err := h.svc.MyVehicles(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) Transfer(c *gin.Context) {
	vehicleID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle id"})
		return
	}
	var in TransferInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(in.BuyerID) == "" && strings.TrimSpace(in.BuyerEmail) != "" {
		buyer, err := h.users.FindByEmail(c.Request.Context(), strings.ToLower(strings.TrimSpace(in.BuyerEmail)))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusBadRequest, gin.H{"error": "buyer not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		in.BuyerID = buyer.ID.Hex()
	}
	user := middleware.GetAuthUser(c)
	out, err := h.svc.RequestTransfer(c.Request.Context(), vehicleID, user.ID, in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *Handlers) ListAll(c *gin.Context) {
	out, err := h.svc.AllVehicles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handlers) ApproveTransfer(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	admin := middleware.GetAuthUser(c)
	out, err := h.svc.ApproveTransfer(c.Request.Context(), id, admin.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
