package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"openedu/internal/models"
	"openedu/internal/parser"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fall back to default if environment variable is not set
		panic("Warning: DATABASE_URL not set, using default connection string")
	}
	slog.Info("got database url", "url", dsn)

	db, err := connectDB(dsn)
	if err != nil {
		panic(err)
	}
	var record models.DBPayload
	db.Model(&models.DBPayload{}).Where("id = ?", 12).First(&record)

	problems := parser.ParseContent(record.Payload.Data)
	for i, v := range problems {
		fmt.Println(i, v)
	}
	// for idx, v := range problems {
	// 	// p := models.FromProblem(v)
	// 	// db.Save(&p)
	// 	fmt.Println("saved", idx, v.Title)
	// }
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
