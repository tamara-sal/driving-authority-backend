package licenses

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
	lic := rg.Group("/licenses", jwt.RequireAuth())
	lic.POST("", middleware.RequireRole(domain.RoleCitizen), m.h.Create)
	lic.GET("/me", middleware.RequireRole(domain.RoleCitizen), m.h.MyLicenses)
	lic.PUT("/:id/renew", middleware.RequireRole(domain.RoleCitizen), m.h.Renew)

	admin := rg.Group("/admin/licenses", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.PUT("/:id/approve", m.h.Approve)
}
