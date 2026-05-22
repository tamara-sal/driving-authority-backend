package analytics

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
	admin := rg.Group("/admin/analytics", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.GET("/overview", m.h.Overview)
	admin.GET("/revenue", m.h.Revenue)
	admin.GET("/exams", m.h.Exams)
	admin.GET("/trends", m.h.Trends)
}
