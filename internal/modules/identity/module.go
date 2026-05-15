package identity

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
	identity := rg.Group("/identity", jwt.RequireAuth())
	identity.POST("/submit", middleware.RequireRole(domain.RoleCitizen), m.h.Submit)
	identity.GET("/status", middleware.RequireRole(domain.RoleCitizen), m.h.MyStatus)

	admin := rg.Group("/admin/identity", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.PUT("/:id/approve", m.h.Approve)
	admin.PUT("/:id/reject", m.h.Reject)
}
