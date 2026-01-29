package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type CORSConfig struct {
	AllowedOrigins string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			DSN:             getEnv("DATABASE_DSN", "sqlserver://sa:YourStrong!Passw0rd@localhost:1433?database=goarch"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 5),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnvAsInt("JWT_EXPIRATION", 24),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
		},
	}

	// If running in development and using SQL Server DSN, disable encryption by default
	if cfg.Server.Env == "development" {
		dsn := cfg.Database.DSN
		if len(dsn) >= 11 && dsn[:11] == "sqlserver://" {
			if !containsParam(dsn, "encrypt") {
				// append encrypt=disable preserving existing query string
				if strings.Contains(dsn, "?") {
					dsn = dsn + "&encrypt=disable"
				} else {
					dsn = dsn + "?encrypt=disable"
				}
				cfg.Database.DSN = dsn
			}
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func containsParam(dsn, param string) bool {
	// simple check for presence of param name in query string
	return strings.Contains(dsn, param+"=")
}
