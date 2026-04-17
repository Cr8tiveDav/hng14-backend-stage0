package services

import (
	"fmt"
	"hng14-stage0-api-data-processing/models"
	"hng14-stage0-api-data-processing/pkg/network"
	"hng14-stage0-api-data-processing/pkg/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func ProcessIdentity(name string) (*models.ApiResult, error) {
	// Manage multiple concurrency and catch error
	var taskManager errgroup.Group

	// Declare receiver struct
	var genderData models.GenderizeResp
	var ageData models.AgifyResp
	var nationData models.NationalizeResp

	// Start background task using Goroutine
	// Call External API simultaneously
	// Call Genderize API
	taskManager.Go(func() error {
		return network.FetchData(fmt.Sprintf("https://api.genderize.io?name=%s", name), "Genderize", &genderData)
	})
	// Call Agify API
	taskManager.Go(func() error {
		return network.FetchData(fmt.Sprintf("https://api.agify.io?name=%s", name), "Agify", &ageData)
	})
	// Call Nationalize API
	taskManager.Go(func() error {
		return network.FetchData(fmt.Sprintf("https://api.nationalize.io/?name=%s", name), "Nationalize", &nationData)
	})

	// Wait for all task to complete
	// If any task fails, taskManage.Go short circuit and return the error immediately
	err := taskManager.Wait()
	if err != nil {
		return nil, err // Return exact network error
	}

	// Edge cases:
	// Check if Genderize API returns a count of 0 or a null gender
	if genderData.Count == 0 || genderData.Gender == nil {
		return nil, utils.External("Genderize returned an invalid response", err)
	}

	// Check if Agify API returns an age of null
	if ageData.Age == nil {
		return nil, utils.External("Agify returned an invalid response", err)
	}

	// Check if Nationalize API returns no country data
	if len(nationData.Country) == 0 {
		return nil, utils.External("Nationalize returned an invalid response", err)
	}

	// Classify age data
	ageGroup := "senior" // Default if age is 60+
	switch {
	case *ageData.Age <= 12:
		ageGroup = "child"
	case *ageData.Age <= 19:
		ageGroup = "teenager"
	case *ageData.Age <= 59:
		ageGroup = "adult"
	}

	var countryID string
	var probability float64

	// Get the first item in the country list
	// since the data is pre-sorted with the highest probability
	if len(nationData.Country) > 0 {
		countryID = nationData.Country[0].ID
		probability = nationData.Country[0].Probability
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return nil, utils.Internal("Failed to generate uuid", err)
	}

	result := models.ApiResult{
		ID:                 newID,
		Name:               name,
		Gender:             *genderData.Gender,
		GenderProbability:  genderData.GenderProbability,
		SampleSize:         genderData.Count,
		Age:                *ageData.Age,
		AgeGroup:           ageGroup,
		CountryID:          countryID,
		CountryProbability: probability,
		CreatedAt:          time.Now().UTC().Format(time.RFC3339),
	}

	return &result, nil
}
