package payments

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
	pay := rg.Group("/payments", jwt.RequireAuth())
	pay.POST("/initiate", middleware.RequireRole(domain.RoleCitizen), m.h.Initiate)
	pay.GET("/history", middleware.RequireRole(domain.RoleCitizen), m.h.History)

	admin := rg.Group("/admin/payments", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.PUT("/:id/mark-paid", m.h.MarkPaid)
}
