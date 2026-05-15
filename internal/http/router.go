package http

import (
	"net/http"

	"driving-authority-backend/internal/config"
	"driving-authority-backend/internal/domain"
	"driving-authority-backend/internal/http/middleware"
	"driving-authority-backend/internal/modules/auth"
	"driving-authority-backend/internal/modules/identity"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(cfg config.Config, db *mongo.Database) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	jwt := middleware.NewJWT(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTLMinute)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.GET("/health", Health)

	authModule := auth.NewModule(db, jwt, cfg.BootstrapAdminSecret)
	authModule.RegisterRoutes(api)

	identityModule := identity.NewModule(db)
	identityModule.RegisterRoutes(api, jwt)

	api.GET("/me", jwt.RequireAuth(), Me)

	api.GET("/admin/ping", jwt.RequireAuth(), middleware.RequireRole(domain.RoleAdmin), AdminPing)

	return r
}
