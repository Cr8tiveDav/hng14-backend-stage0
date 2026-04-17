# Stage 1: Profile Intelligence Service

A Go-based REST API that enriches user names with demographic data using Genderize, Agify, and Nationalize APIs. Built for the HNG Internship Stage 1 (Backend) task.

**Live URL:** [hng14-backend-stage1-production.up.railway.app](https://hng14-backend-stage1-production.up.railway.app)

## Features

- **Multi-API Integration:** Aggregates data from Genderize, Agify, and Nationalize.
- **Data Persistence:** Stores processed profile intelligence in a PostgreSQL database.
- **Idempotency:** Intelligent handling of duplicate submissions to prevent redundant records.
- **Advanced Filtering:** Supports case-insensitive query parameters for gender, country, and age group.
- **UUID v7 & ISO 8601:** Uses the latest UUID standards for identifiers and UTC timestamps for consistency.

## Tech Stack

- **Language:** Go (Vanilla net/http)
- **Database:** PostgreSQL
- **Integrations:** Genderize.io, Agify.io, Nationalize.io
- **Deployment:** Railway

## API Documentation

### 1. Create Profile

Enriches a name and stores the result. Returns existing data if the name has already been processed.

- **URL:** `/api/profiles`
- **Method:** `POST`
- **Body:** `{ "name": "ella" }`

### 2. Get Profile by ID

Retrieves a specific profile using its UUID.

- **URL:** `/api/profiles/{id}`
- **Method:** `GET`

### 3. List Profiles

Retrieves all stored profiles with optional filtering.

- **URL:** `/api/profiles`
- **Method:** `GET`
- **Query Params:** `gender`, `country_id`, `age_group`
- **Example:** `/api/profiles?gender=female&country_id=US`

### 4. Delete Profile

Removes a profile from the database.

- **URL:** `/api/profiles/{id}`
- **Method:** `DELETE`
- **Response:** `204 No Content`

### Local Setup

Clone the repository:

```zsh
git clone <your-repo-link>
cd <repo-folder>
Initialize modules:
```

````

Initialize modules

```zsh
go mod tidy
Run the server:
```

Run the server

```zsh
go run main.go
```
````
