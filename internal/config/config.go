package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port   string

	MongoURI string
	MongoDB  string

	JWTSecret          string
	JWTIssuer          string
	JWTAccessTTLMinute int

	BootstrapAdminSecret string
	SeedDemoUsers        bool
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:               getEnv("APP_ENV", "dev"),
		Port:                 getEnv("PORT", "8080"),
		MongoURI:             getEnv("MONGO_URI", ""),
		MongoDB:              getEnv("MONGO_DB", "driving_authority"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTIssuer:            getEnv("JWT_ISSUER", "driving-authority"),
		BootstrapAdminSecret: getEnv("BOOTSTRAP_ADMIN_SECRET", ""),
		SeedDemoUsers:        getEnvBool("SEED_DEMO_USERS", true),
	}

	ttlStr := getEnv("JWT_ACCESS_TTL_MINUTES", "60")
	ttl, err := strconv.Atoi(ttlStr)
	if err != nil || ttl <= 0 {
		return Config{}, errors.New("invalid JWT_ACCESS_TTL_MINUTES")
	}
	cfg.JWTAccessTTLMinute = ttl

	if cfg.MongoURI == "" {
		return Config{}, errors.New("MONGO_URI is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
