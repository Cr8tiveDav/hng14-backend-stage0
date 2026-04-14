# Stage 0: Backend API Integration & Data Processing Assessment

A Go-based REST API that classifies a name's gender using the Genderize.io API and process the data. Built as part of the HNG Internship Stage 0 (Backend) task.
**Live URL:** [https://hng14-backend-stage0.onrender.com](https://hng14-backend-stage0.onrender.com)

## Features

- **Data Transformation:** Maps external API data to a custom internal structure.
- **Input Validation:** Strict handling of missing names (400) and numeric inputs (422).
- **CORS Enabled:** Fully accessible for cross-origin frontend requests. -**ISO 8601 Timestamps:** Every response includes a dynamic UTC processed_at field.

## Tech Stack

- **Language:** Go (Vanilla net/http)
- **Integration:** Genderize.io API
- **Deployment:** [Insert your platform here, e.g., Render/Railway]

## API Documentation

Get Gender Classification
Returns gender data and confidence scores for a provided name.

- **URL:** /api/classify
- **Method:** GET
- **Query Params:** name=[string]

### Successful Response (200 OK)

```JSON
{
  "status": "success",
  "data": {
    "name": "<name>",
    "gender": "male",
    "probability": 0.98,
    "sample_size": 1234,
    "is_confident": true,
    "processed_at": "2026-04-14T13:09:31Z"
  }
}
```

### Local Setup

Clone the repository:

```zsh
git clone <your-repo-link>
cd <repo-folder>
Initialize modules:
```

Initialize modules

```zsh
go mod tidy
Run the server:
```

Run the server

```zsh
go run main.go
```
