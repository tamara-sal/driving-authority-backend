package vehicles

import (
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	h *Handlers
}

func NewModule(db *mongo.Database) *Module {
	repo := NewRepo(db)
	svc := NewService(repo)
	return &Module{h: NewHandlers(svc)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwt *middleware.JWT) {
	v := rg.Group("/vehicles", jwt.RequireAuth(), middleware.RequireRole(domain.RoleCitizen))
	v.POST("", m.h.Create)
	v.GET("/me", m.h.MyVehicles)
	v.POST("/:id/transfer", m.h.Transfer)

	admin := rg.Group("/admin/transfer", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.PUT("/:id/approve", m.h.ApproveTransfer)
}
