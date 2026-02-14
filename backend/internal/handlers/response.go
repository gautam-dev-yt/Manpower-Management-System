package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// ErrorResponse represents a standardized API error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// JSONError writes a JSON error response.
func JSONError(w http.ResponseWriter, status int, message string) {
	errResp := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Status:  status,
	}
	JSON(w, status, errResp)
}

// PaginationMeta holds pagination metadata for list responses.
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResponse wraps a list with pagination metadata.
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}
