package main

import (
	"fmt"
	"hng14-stage0-api-data-processing/internal/handlers"
	"hng14-stage0-api-data-processing/internal/repository"
	"hng14-stage0-api-data-processing/middleware"
	"hng14-stage0-api-data-processing/pkg/db"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	// 1. Load the .env file first!
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Call db connect
	database, err := db.Connect()
	if err != nil {
		fmt.Printf("Database skipped: %v\n", err)
	} else {
		defer database.Close()
	}

	// Pass it to the repo
	repo := &repository.ProfileRepository{DB: database}
	// ... rest of setup

	h := &handlers.Handler{Repo: repo}

	mux := http.NewServeMux()
	// Register handler
	mux.HandleFunc("/api/classify", h.DetermineGender)
	mux.HandleFunc("/api/profiles", h.ProfilesBaseHandler)
	mux.HandleFunc("/api/profiles/", h.ProfileIDHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port " + port)
	err = http.ListenAndServe(":"+port, middleware.EnableCORS(mux))
	if err != nil {
		log.Fatal(err)
	}
}
