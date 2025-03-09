package main

import (
	"database/sql"
	"log/slog"
	"openedu/internal/handlers"
	"openedu/internal/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	slog.Info("starting application...")
	// Get database connection string from environment variable or use default
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fall back to default if environment variable is not set
		panic("Warning: DATABASE_URL not set, using default connection string")
	}
	slog.Info("got database url", "url", dsn)

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 1024,
	})
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(logger.New()) // Add this line to enable logging
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
	}))

	db, err := connectDB(dsn)
	if err != nil {
		panic("failed to connect database")
	}
	slog.Info("connected to db")

	models.Migrate(db)

	h := handlers.New(db)

	app.Post("/collect", h.CollectHandler)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

func connectDB(dsn string) (*gorm.DB, error) {
	pg, err := sql.Open(
		"pgx",
		dsn,
	) // openedu
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.New(postgres.Config{
		Conn: pg,
	}), &gorm.Config{})
}
