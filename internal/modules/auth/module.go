package auth

import (
	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	h *Handlers
}

func NewModule(db *mongo.Database, jwt *middleware.JWT, bootstrapAdminSecret string) *Module {
	repo := NewUserRepo(db)
	svc := NewService(repo, jwt, bootstrapAdminSecret)
	return &Module{h: NewHandlers(svc)}
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", m.h.Register)
	auth.POST("/login", m.h.Login)
	auth.POST("/bootstrap-admin", m.h.BootstrapAdmin)
}
