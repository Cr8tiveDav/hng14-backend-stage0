package models

import (
	"github.com/google/uuid"
)

type GenderizeRequest struct {
	Name string `json:"name"`
}

type GenderizeData struct {
	Count       int64   `json:"count"`
	Name        string  `json:"name"`
	Gender      *string `json:"gender"`
	Probability float64 `json:"probability"`
}

type ProcessedData struct {
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
	SampleSize  int64   `json:"sample_size"`
	IsConfident bool    `json:"is_confident"`
	ProcessedAt string  `json:"processed_at"`
}

type SuccessResponse struct {
	Status string        `json:"status"`
	Data   ProcessedData `json:"data"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// My API response
type ApiResult struct {
	ID                 uuid.UUID `json:"id,omitempty"`
	Name               string    `json:"name"`
	Gender             string    `json:"gender"`
	GenderProbability  float64   `json:"gender_probability,omitempty"`
	SampleSize         int64     `json:"sample_size,omitempty"`
	Age                int64     `json:"age"`
	AgeGroup           string    `json:"age_group"`
	CountryID          string    `json:"country_id"`
	CountryProbability float64   `json:"country_probability,omitempty"`
	CreatedAt          string    `json:"created_at,omitempty"`
}
type DataResponse struct {
	ID                 string  `json:"id,omitempty"`
	Name               string  `json:"name"`
	Gender             string  `json:"gender"`
	GenderProbability  float64 `json:"gender_probability,omitempty"`
	SampleSize         int64   `json:"sample_size,omitempty"`
	Age                int64   `json:"age"`
	AgeGroup           string  `json:"age_group"`
	CountryID          string  `json:"country_id"`
	CountryProbability float64 `json:"country_probability,omitempty"`
	CreatedAt          string  `json:"created_at,omitempty"`
}

type ProfileResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message,omitempty"` // Optional
	Count   int    `json:"count,omitempty"`   // Optional
}

// Internal API response
type GenderizeResp struct {
	Gender            *string `json:"gender"`
	GenderProbability float64 `json:"probability"`
	Count             int64   `json:"count"`
}

type AgifyResp struct {
	Age *int64 `json:"age"`
}

type CountryDetails struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}
type NationalizeResp struct {
	Country []CountryDetails `json:"country"`
}
