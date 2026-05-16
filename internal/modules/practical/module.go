package practical

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
	centers := rg.Group("/centers", jwt.RequireAuth())
	centers.GET("", m.h.ListCenters)
	centers.GET("/:id/slots", m.h.ListSlots)

	practical := rg.Group("/practical", jwt.RequireAuth(), middleware.RequireRole(domain.RoleCitizen))
	practical.POST("/book", m.h.Book)

	examiner := rg.Group("/examiner/practical", jwt.RequireAuth(), middleware.RequireRole(domain.RoleExaminer))
	examiner.PUT("/:id/result", m.h.RecordResult)
}
