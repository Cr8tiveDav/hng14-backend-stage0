package main

import (
	"hng14-stage0-api-data-processing/middleware"
	"hng14-stage0-api-data-processing/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	// Register handler
	mux.HandleFunc("/api/classify", routes.GenderizeName)

port := os.Getenv("PORT")
if port == "" {
	port = "8080"
}

	log.Println("Server is running on port " + port)
	err := http.ListenAndServe(":"+port, middleware.EnableCORS(mux))
	if err != nil {
		log.Fatal(err)
	}
}
