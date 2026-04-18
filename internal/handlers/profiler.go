package handlers

import (
	"encoding/json"
	"fmt"
	"hng14-stage0-api-data-processing/internal/repository"
	"hng14-stage0-api-data-processing/internal/services"
	"hng14-stage0-api-data-processing/models"
	"hng14-stage0-api-data-processing/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	Repo *repository.ProfileRepository
}

// Dispatcher function for base handler
func (h *Handler) ProfilesBaseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListProfiles(w, r)
	case http.MethodPost:
		h.CreateProfile(w, r)
	default:
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Dispatcher function for ID handler
func (h *Handler) ProfileIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetProfile(w, r)
	case http.MethodDelete:
		h.DeleteProfile(w, r)
	default:
		utils.HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Determine the gender of a name (Task 0)
func (h *Handler) DetermineGender(w http.ResponseWriter, r *http.Request) {

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
	if data.Count == 0 || data.Gender == nil {
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
		Gender:      *data.Gender,
		Probability: data.Probability,
		SampleSize:  data.Count,
		IsConfident: isConfident,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Create the success response
	successResponse := models.SuccessResponse{
		Status: "success",
		Data:   processedData,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

// Process and create profile (Task 1)
func (h *Handler) CreateProfile(w http.ResponseWriter, r *http.Request) {


// Get name from request body
	var req models.GenderizeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Missing or empty name")
		return
	}

	name := req.Name

	// Validate name parameter is not empty
	if name == "" {
		utils.HandleError(w, http.StatusBadRequest, "Name parameter is required")
		return
	}
	// Validate name parameter is a string
	_, err = strconv.Atoi(name)
	if err == nil {
		utils.HandleError(w, http.StatusUnprocessableEntity, "Name parameter must be a string")
		return
	}

	// Call external API
	result, err := services.ProcessIdentity(name)
	if err != nil {
		if netErr, ok := err.(*utils.NetworkError); ok {
			if netErr.Internal != nil {
				// Internal Error
				utils.HandleError(w, netErr.StatusCode, netErr.Message)
				return
			}
			// Other Error (Network and Logic Error)
			utils.HandleError(w, netErr.StatusCode, netErr.Message)
			return
		}
	}

	savedData, exists, err := h.Repo.Save(result)
	if err != nil {
		if netErr, ok := err.(*utils.NetworkError); ok {
			if netErr.Internal != nil {
				http.Error(w, netErr.Message, netErr.StatusCode)
				return
			}
			http.Error(w, netErr.Message, netErr.StatusCode)
			return
		}
	}

	// Check for existing name and return if with the data if true
	if exists {
		savedResponse := models.ProfileResponse{
			Status:  "success",
			Message: "Profile already exists",
			Data:    *savedData,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(savedResponse)
		return
	}

	profileResponse := models.ProfileResponse{
		Status: "success",
		Data:   *result,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profileResponse)

}

// Get profile by ID
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/profiles/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Fetch from Repo
	profile, err := h.Repo.GetByID(id)
	if err != nil {
		if netErr, ok := err.(*utils.NetworkError); ok {
			if netErr.Internal != nil {
				http.Error(w, netErr.Message, netErr.StatusCode)
				return
			}
			http.Error(w, netErr.Message, netErr.StatusCode)
			return
		}
	}

	profileResponse := models.ProfileResponse{
		Status: "success",
		Data:   *profile,
	}

	// Success Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profileResponse)
}

// Get profiles
func (h *Handler) ListProfiles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	gender := query.Get("gender")
	country := query.Get("country_id")
	ageGroup := query.Get("age_group")

	profiles, err := h.Repo.List(gender, country, ageGroup)
	if err != nil {
		fmt.Println("Error:",err)
		utils.HandleError(w, http.StatusInternalServerError, "Failed to fetch profiles")
		return
	}

	profilesResponse := models.ProfileResponse{
		Status: "success",
		Data:   profiles,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profilesResponse)
}

// Delete profile
func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/profiles/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	if err := h.Repo.Delete(id); err != nil {
		utils.HandleError(w, http.StatusInternalServerError, "Failed to delete profile")
		return
	}

	// No content on success
	w.WriteHeader(http.StatusNoContent)
}
