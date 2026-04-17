package utils

import (
	"encoding/json"
	"hng14-stage0-api-data-processing/models"
	"net/http"
)

func HandleError(w http.ResponseWriter, statusCode int, message string) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	// Write status code
	w.WriteHeader(statusCode)

	errRes := models.ErrorResponse{
		Status:  "error",
		Message: message,
	}
	// Encode the error response as JSON and write it to the response writer
	json.NewEncoder(w).Encode(errRes)
}
