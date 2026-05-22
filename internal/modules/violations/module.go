package violations

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
	v := rg.Group("/violations", jwt.RequireAuth())
	v.GET("", m.h.List)
	v.POST("", middleware.RequireRole(domain.RoleOfficer, domain.RoleAdmin), m.h.Create)
	v.PUT("/:id/status", middleware.RequireRole(domain.RoleOfficer, domain.RoleAdmin), m.h.UpdateStatus)
}
