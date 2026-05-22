package http

import (
	"net/http"

	"driving-authority-backend/internal/config"
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"
	"driving-authority-backend/internal/modules/analytics"
	"driving-authority-backend/internal/modules/auth"
	"driving-authority-backend/internal/modules/exams"
	"driving-authority-backend/internal/modules/identity"
	"driving-authority-backend/internal/modules/inspections"
	"driving-authority-backend/internal/modules/licenses"
	"driving-authority-backend/internal/modules/monitoring"
	"driving-authority-backend/internal/modules/notifications"
	"driving-authority-backend/internal/modules/payments"
	"driving-authority-backend/internal/modules/platform"
	"driving-authority-backend/internal/modules/practical"
	"driving-authority-backend/internal/modules/vehicles"
	"driving-authority-backend/internal/modules/violations"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cfg config.Config, db *mongo.Database) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	jwt := middleware.NewJWT(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTLMinute)

	r.GET("/", Root)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.GET("/health", Health)

	authModule := auth.NewModule(db, jwt, cfg.BootstrapAdminSecret)
	authModule.RegisterRoutes(api)

	identityModule := identity.NewModule(db)
	identityModule.RegisterRoutes(api, jwt)

	licensesModule := licenses.NewModule(db)
	licensesModule.RegisterRoutes(api, jwt)

	examsModule := exams.NewModule(db)
	examsModule.RegisterRoutes(api, jwt)

	practicalModule := practical.NewModule(db)
	practicalModule.RegisterRoutes(api, jwt)

	vehiclesModule := vehicles.NewModule(db)
	vehiclesModule.RegisterRoutes(api, jwt)

	inspectionsModule := inspections.NewModule(db)
	inspectionsModule.RegisterRoutes(api, jwt)

	monitoringModule := monitoring.NewModule(db)
	monitoringModule.RegisterRoutes(api, jwt)

	paymentsModule := payments.NewModule(db)
	paymentsModule.RegisterRoutes(api, jwt)

	analyticsModule := analytics.NewModule(db)
	analyticsModule.RegisterRoutes(api, jwt)

	notificationsModule := notifications.NewModule(db)
	notificationsModule.RegisterRoutes(api, jwt)

	violationsModule := violations.NewModule(db)
	violationsModule.RegisterRoutes(api, jwt)

	platformModule := platform.NewModule(db)
	platformModule.RegisterRoutes(api, jwt)

	api.GET("/admin/ping", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin), AdminPing)

	return r
}
