package platform

import (
	"driving-authority-backend/internal/audit"
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"
	"driving-authority-backend/internal/modules/auth"
	"driving-authority-backend/internal/modules/inspections"
	"driving-authority-backend/internal/modules/licenses"
	"driving-authority-backend/internal/modules/vehicles"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	h *Handlers
}

func NewModule(db *mongo.Database) *Module {
	svc := NewService(
		auth.NewUserRepo(db),
		licenses.NewRepo(db),
		vehicles.NewRepo(db),
		inspections.NewRepo(db),
		audit.NewLogger(db),
	)
	return &Module{h: NewHandlers(svc)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwt *middleware.JWT) {
	rg.GET("/activity", jwt.RequireAuth(), m.h.ListActivity)

	admin := rg.Group("/admin", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.GET("/users", m.h.ListUsers)
	admin.GET("/applications", m.h.ListApplications)
	admin.GET("/audit-logs", m.h.ListAuditLogs)
}
