package utils

import "net/http"

type NetworkError struct {
	StatusCode int
	Message    string
	Internal   error // Store the original/technical error
}

// Error method
func (e *NetworkError) Error() string {
	return e.Message
}

// Status error construction func
func BadRequest(message string) *NetworkError {
	return &NetworkError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

func UnProcessableEntity(message string) *NetworkError {
	return &NetworkError{
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
	}
}

func NotFound(message string) *NetworkError {
	return &NetworkError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func External(message string, err error) *NetworkError {
	return &NetworkError{
		StatusCode: http.StatusBadGateway,
		Message:    message,
		Internal:   err,
	}
}

func Internal(message string, err error) *NetworkError {
	return &NetworkError{
		StatusCode: http.StatusInternalServerError,
		Message:    message,
		Internal:   err,
	}
}
