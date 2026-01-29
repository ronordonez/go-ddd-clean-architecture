package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"go-architecture/internal/product/application"
	"go-architecture/internal/product/infra/http"
	sharedhttp "go-architecture/internal/shared/http"
	"go-architecture/internal/product/infra/mssql"
	"go-architecture/internal/shared/config"
	"go-architecture/internal/shared/logger"
	"go-architecture/internal/shared/middleware"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting application...")

	// Load .env file (if present) so env vars from .env are available
	if err := godotenv.Load(); err != nil {
		// not fatal — proceed, maybe env vars are set externally
		log.Info(".env file not found or could not be loaded; relying on environment variables")
	} else {
		log.Info("Loaded .env file")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize database
	db, err := initDatabase(cfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		AppName:      "Go Architecture API",
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))
	app.Use(middleware.RequestLogger(log))
	app.Use(middleware.RateLimiter())

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Initialize dependencies - Product module
	productRepo := mssql.NewProductRepository(db)
	productService := application.NewProductService(productRepo)
	productHandler := http.NewProductHandler(productService, log)

	// API routes
	api := app.Group("/api/v1")

	// Auth
	api.Post("/login", sharedhttp.LoginHandler(cfg))

	// Product routes with JWT protection
	products := api.Group("/products")
	products.Get("/", productHandler.GetAll)
	products.Get("/:id", productHandler.GetByID)
	products.Post("/", middleware.JWTProtected(cfg.JWT.Secret), productHandler.Create)
	products.Put("/:id", middleware.JWTProtected(cfg.JWT.Secret), productHandler.Update)
	products.Delete("/:id", middleware.JWTProtected(cfg.JWT.Secret), productHandler.Delete)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatal("Server error", "error", err)
		}
	}()

	log.Info("Server started on port " + cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	if err := app.ShutdownWithContext(context.Background()); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server stopped")
}

func initDatabase(cfg *config.Config, log *logger.Logger) (*sqlx.DB, error) {
	// Log masked DSN for debugging (password masked)
	log.Info("DB DSN", "dsn", maskDSN(cfg.Database.DSN))

	var db *sqlx.DB
	var err error
	// retry a few times in case SQL Server is starting
	attempts := 3
	for i := 1; i <= attempts; i++ {
		db, err = sqlx.Connect("sqlserver", cfg.Database.DSN)
		if err == nil {
			break
		}
		log.Error("DB connection attempt failed", "attempt", i, "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		// Provide a clearer error when authentication fails
		if strings.Contains(strings.ToLower(err.Error()), "login failed") || strings.Contains(strings.ToLower(err.Error()), "login error") {
			host := extractHost(cfg.Database.DSN)
			return nil, fmt.Errorf("mssql: login failed for user. Host=%s. Check SQL Server authentication mode (enable SQL auth), ensure 'sa' is enabled and password is correct, and that TCP/IP is enabled. You can test with: sqlcmd -S %s -U sa -P '<password>'", host, host)
		}
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	return db, nil
}

func maskDSN(dsn string) string {
	// crude mask: replace password between : and @ for sqlserver://user:pass@host
	// keep scheme and host
	if len(dsn) == 0 {
		return dsn
	}
	// find '://' and '@'
	at := strings.Index(dsn, "@")
	if at == -1 {
		return dsn
	}
	// find last ':' before @ (password separator)
	slash := strings.Index(dsn, "://")
	if slash == -1 || slash+3 >= at {
		return dsn
	}
	sub := dsn[slash+3 : at]
	colon := strings.Index(sub, ":")
	if colon == -1 {
		return dsn
	}
	user := sub[:colon]
	masked := user + ":*****"
	return dsn[:slash+3] + masked + dsn[at:]
}

func extractHost(dsn string) string {
	// Attempt to extract host:port from DSN like sqlserver://user:pass@host:port?...
	at := strings.Index(dsn, "@")
	if at == -1 || at+1 >= len(dsn) {
		return "localhost:1433"
	}
	rest := dsn[at+1:]
	// stop at ?
	q := strings.Index(rest, "?")
	if q != -1 {
		rest = rest[:q]
	}
	// rest may be host or host:port or \instance; normalize commas
	// replace backslash with empty (instance names require SQL Browser)
	rest = strings.ReplaceAll(rest, "\\", "")
	return rest
}
