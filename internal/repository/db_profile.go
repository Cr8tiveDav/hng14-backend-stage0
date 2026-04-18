package repository

import (
	"database/sql"
	"fmt"
	"hng14-stage0-api-data-processing/models"
	"hng14-stage0-api-data-processing/pkg/utils"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ProfileRepository struct {
	DB *sql.DB
}

// Save result to the database
func (r *ProfileRepository) Save(p *models.ApiResult) (*models.ApiResult, bool, error) {
	query := `
        INSERT INTO user_profiles (
            id, name, gender, gender_probability, sample_size,
            age, age_group, country_id, country_probability, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id`

	// Generate UUID if not already present
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	err := r.DB.QueryRow(query,
		p.ID, p.Name, p.Gender, p.GenderProbability, p.SampleSize,
		p.Age, p.AgeGroup, p.CountryID, p.CountryProbability, p.CreatedAt,
	).Scan(&p.ID)

	if err != nil {
		// Check if error is a Unique Constraint violation (Postgres code 23505)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// Fetch the existing record and return it to the user
			existing, fetchErr := r.GetByName(p.Name)
			if fetchErr != nil {
				return nil, false, fetchErr
			}
			return existing, true, nil // true = "already existed"
		}
		return nil, false, err
	}

	return p, false, nil // false = "newly created"
}

// Helper function to get already existing name
func (r *ProfileRepository) GetByName(name string) (*models.ApiResult, error) {
	query := `
        SELECT
            id, name, gender, gender_probability, sample_size,
            age, age_group, country_id, country_probability, created_at
        FROM user_profiles
        WHERE name = $1`

	var p models.ApiResult

	err := r.DB.QueryRow(query, name).Scan(
		&p.ID,
		&p.Name,
		&p.Gender,
		&p.GenderProbability,
		&p.SampleSize,
		&p.Age,
		&p.AgeGroup,
		&p.CountryID,
		&p.CountryProbability,
		&p.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No profile found, not an error for the caller
		}
		return nil, utils.External("Upstream", err) // Actual database connection or syntax error
	}

	return &p, nil
}

// Get all profiles in the database
func (r *ProfileRepository) List(gender, country, ageGroup string) ([]models.DataResponse, error) {
	query := `SELECT name, gender, age, age_group, country_id FROM user_profiles WHERE 1=1`
	var args []interface{}
	counter := 1

	if gender != "" {
		query += fmt.Sprintf(" AND LOWER(gender) = LOWER($%d)", counter)
		args = append(args, gender)
		counter++
	}
	if country != "" {
		query += fmt.Sprintf(" AND LOWER(country_id) = LOWER($%d)", counter)
		args = append(args, country)
		counter++
	}
	if ageGroup != "" {
		query += fmt.Sprintf(" AND LOWER(age_group) = LOWER($%d)", counter)
		args = append(args, ageGroup)
		counter++
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.DataResponse
	idCount := 1
	for rows.Next() {
		var p models.DataResponse
		err := rows.Scan(&p.Name, &p.Gender, &p.Age, &p.AgeGroup, &p.CountryID)
		if err != nil {
			return nil, err
		}

		// Set custom ID
		p.ID = fmt.Sprintf("id-%d", idCount)
		idCount++
		results = append(results, p)
	}
	return results, nil
}

// Get by id
func (r *ProfileRepository) GetByID(id uuid.UUID) (*models.ApiResult, error) {
	query := `SELECT id, name, gender, gender_probability, sample_size, age, age_group, country_id, country_probability, created_at
	          FROM user_profiles WHERE id = $1`
	var p models.ApiResult
	err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Gender, &p.GenderProbability, &p.SampleSize, &p.Age, &p.AgeGroup, &p.CountryID, &p.CountryProbability, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, utils.NotFound("Profile not found")
	}
	return &p, err
}

// Delete data from the database
func (r *ProfileRepository) Delete(id uuid.UUID) error {
	_, err := r.DB.Exec("DELETE FROM user_profiles WHERE id = $1", id)
	return err
}
