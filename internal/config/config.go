package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                    string
	ShouldMigrate                  bool
	Port                           string
	HydraAdminURL                  string
	KratosAdminURL                 string
	KratosPublicURL                string
	KratosBrowserURL               string
	KratosUIURL                    string
	PublicURL                      string
	HydraIntrospectionClientID     string
	HydraIntrospectionClientSecret string
}

var current Config

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DatabaseURL:                    strings.TrimSpace(getEnv("DATABASE_URL", "postgresql://root@cockroach:26257/senior_citizen?sslmode=disable")),
		ShouldMigrate:                  getEnvBool("SHOULD_MIGRATE", false),
		Port:                           strings.TrimSpace(getEnv("PORT", "3000")),
		HydraAdminURL:                  strings.TrimSpace(getEnv("HYDRA_ADMIN_URL", "http://127.0.0.1:4445")),
		KratosAdminURL:                 strings.TrimSpace(getEnv("KRATOS_ADMIN_URL", "http://127.0.0.1:4434")),
		KratosPublicURL:                strings.TrimSpace(getEnv("KRATOS_PUBLIC_URL", "http://127.0.0.1:4433")),
		KratosBrowserURL:               strings.TrimSpace(getEnv("KRATOS_BROWSER_URL", "http://127.0.0.1:4433")),
		KratosUIURL:                    strings.TrimSpace(getEnv("KRATOS_UI_URL", "http://127.0.0.1:3000/ui")),
		PublicURL:                      strings.TrimSpace(getEnv("PUBLIC_URL", "http://127.0.0.1:3000")),
		HydraIntrospectionClientID:     strings.TrimSpace(getEnv("HYDRA_INTROSPECTION_CLIENT_ID", "")),
		HydraIntrospectionClientSecret: strings.TrimSpace(getEnv("HYDRA_INTROSPECTION_CLIENT_SECRET", "")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	current = cfg
	return cfg, nil
}

func Get() Config {
	return current
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
