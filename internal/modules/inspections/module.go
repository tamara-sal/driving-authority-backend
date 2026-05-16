package inspections

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
	insp := rg.Group("/inspection", jwt.RequireAuth(), middleware.RequireRole(domain.RoleCitizen))
	insp.POST("/schedule", m.h.Schedule)
	insp.POST("/:id/upload-report", m.h.UploadReport)
}
