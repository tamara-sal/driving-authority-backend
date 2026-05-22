package vehicles

import (
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"
	"driving-authority-backend/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	h *Handlers
}

func NewModule(db *mongo.Database) *Module {
	repo := NewRepo(db)
	svc := NewService(repo)
	users := auth.NewUserRepo(db)
	return &Module{h: NewHandlers(svc, users)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup, jwt *middleware.JWT) {
	v := rg.Group("/vehicles", jwt.RequireAuth(), middleware.RequireRole(domain.RoleCitizen))
	v.POST("", m.h.Create)
	v.GET("/me", m.h.MyVehicles)
	v.POST("/:id/transfer", m.h.Transfer)

	adminVehicles := rg.Group("/admin/vehicles", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	adminVehicles.GET("", m.h.ListAll)

	admin := rg.Group("/admin/transfer", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin))
	admin.PUT("/:id/approve", m.h.ApproveTransfer)
}
