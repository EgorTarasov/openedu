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
	"github.com/gofiber/template/html/v2"
	_ "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	slog.Info("starting application...")
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("Warning: DATABASE_URL not set, using default connection string")
	}
	slog.Info("got database url", "url", dsn)

	renderEngine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:     renderEngine,
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
	app.Get("/", h.Index)
	app.Post("/collect", h.CollectHandler)
	app.Get("/q", h.Search)

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

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pg,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	models.Migrate(db)
	return db, err
}
