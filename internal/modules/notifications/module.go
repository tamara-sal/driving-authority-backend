package notifications

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"driving-authority-backend/internal/http/middleware"
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
	n := rg.Group("/notifications", jwt.RequireAuth())
	n.GET("", m.h.List)
	n.PATCH("/:id/read", m.h.MarkRead)
}
