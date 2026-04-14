package routes

import (
	"encoding/json"
	"fmt"
	"hng14-stage0-api-data-processing/models"
	"hng14-stage0-api-data-processing/utils"
	"net/http"
	"strconv"
	"time"
)

func GenderizeName(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET
	if r.Method != http.MethodGet {
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get the name parameter from the query string
	name := r.URL.Query().Get("name")

	// Validate name parameter is not empty
	if name == "" {
		utils.HandleError(w, http.StatusBadRequest, "Name parameter is required")
		return
	}
	// Validate name parameter is a string
	_, err := strconv.Atoi(name)
	if err == nil {
		utils.HandleError(w, http.StatusUnprocessableEntity, "Name parameter must be a string")
		return
	}

	// Call the Genderize API
	url := fmt.Sprintf("https://api.genderize.io?name=%s", name)
	resp, err := http.Get(url)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to fetch")
		return
	}

	defer resp.Body.Close()

	// Check if the API call was successful
	if resp.StatusCode != http.StatusOK {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to fetch data from Genderize API")
		return
	}

	// Decode the response from the Genderize API
	var data models.GenderizeData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to decode response")
		return
	}

	// Edge cases: Check if the API returns a count of 0 or an empty gender
	if data.Count == 0 || data.Gender == "" {
		utils.HandleError(w, http.StatusUnprocessableEntity, "No prediction available for the provided name")
		return
	}

	// Determine confidence based on probability and count
	isConfident := false
	if data.Probability >= 0.7 && data.Count >= 100 {
		isConfident = true
	}

	// Build the processed data
	processedData := models.ProcessedData{
		Name:        data.Name,
		Gender:      data.Gender,
		Probability: data.Probability,
		SampleSize:  data.Count,
		IsConfident: isConfident,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Create the success response
	successResponse := models.SuccessResponse{
		Status: "success",
		Data: processedData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}
