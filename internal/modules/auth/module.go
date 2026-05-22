package auth

import (
	"context"

	"driving-authority-backend/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module struct {
	h   *Handlers
	svc *Service
	jwt *middleware.JWT
}

func NewModule(db *mongo.Database, jwt *middleware.JWT, bootstrapAdminSecret string) *Module {
	repo := NewUserRepo(db)
	tokens := NewTokenRepo(db, repo)
	svc := NewService(repo, tokens, jwt, bootstrapAdminSecret)
	return &Module{h: NewHandlers(svc), svc: svc, jwt: jwt}
}

func (m *Module) SeedDemoUsers(ctx context.Context) (SeedDemoOutput, error) {
	return m.svc.SeedDemoUsers(ctx)
}

func (m *Module) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", m.h.Register)
	auth.POST("/login", m.h.Login)
	auth.POST("/verify-email", m.h.VerifyEmail)
	auth.POST("/forgot-password", m.h.ForgotPassword)
	auth.POST("/reset-password", m.h.ResetPassword)
	auth.POST("/bootstrap-admin", m.h.BootstrapAdmin)
	auth.POST("/seed-demo", m.h.SeedDemo)
	rg.GET("/me", m.jwt.RequireAuth(), m.h.Me)
}
