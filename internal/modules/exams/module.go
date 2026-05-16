package exams

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
	exam := rg.Group("/exam", jwt.RequireAuth())
	exam.GET("/questions", middleware.RequireRole(domain.RoleAdmin), m.h.ListQuestions)
	exam.POST("/start", middleware.RequireRole(domain.RoleCitizen), m.h.Start)
	exam.POST("/:attemptId/submit", middleware.RequireRole(domain.RoleCitizen), m.h.Submit)
	exam.GET("/history", middleware.RequireRole(domain.RoleCitizen), m.h.History)
}
