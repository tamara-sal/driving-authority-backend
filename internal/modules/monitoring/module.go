package monitoring

import (
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
	rg.POST("/devices/data", m.h.IngestData)

	mon := rg.Group("/monitoring", jwt.RequireAuth())
	mon.GET("/trips/:vehicleId", m.h.TripsByVehicle)
	mon.GET("/score/:userId", m.h.ScoreByUser)
}
