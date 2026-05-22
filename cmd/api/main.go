// @title           Driving Authority API
// @version         1.0
// @description     REST API for auth, JWT, RBAC, and identity verification.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter: Bearer {your JWT access token}

package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"driving-authority-backend/internal/config"
	"driving-authority-backend/internal/db"
	"driving-authority-backend/internal/http/middleware"
	apihttp "driving-authority-backend/internal/http"
	"driving-authority-backend/internal/modules/auth"

	_ "driving-authority-backend/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := db.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	database := client.Database(cfg.MongoDB)
	if cfg.SeedDemoUsers {
		jwt := middleware.NewJWT(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTTLMinute)
		mod := auth.NewModule(database, jwt, cfg.BootstrapAdminSecret)
		out, err := mod.SeedDemoUsers(ctx)
		if err != nil {
			log.Printf("demo user seed failed: %v", err)
		} else {
			log.Printf("demo users seeded: %v (password: %s)", out.Accounts, out.Password)
		}
	}

	router := apihttp.NewRouter(cfg, database)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("API listening on :%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
