package db

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq" // Postgres driver
)

//go:embed ../../migrations/*.sql
var embedMigrations embed.FS

func Connect() (*sql.DB, error) {
	// Get connection details
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" { dbHost = "localhost" }
	if dbPort == "" { dbPort = "5432" }
	if dbUser == "" { dbUser = "postgres" }
	if dbName == "" { dbName = "postgres" }

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Open Connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL database")

	// TRIGGER MIGRATIONS (The "Automated" Part)
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	log.Println("Running database migrations...")
	// Point to "../../migrations" because that's where the files live relative to this file
	if err := goose.Up(db, "../../migrations"); err != nil {
		return nil, fmt.Errorf("migration failed: %v", err)
	}

	return db, nil
}
