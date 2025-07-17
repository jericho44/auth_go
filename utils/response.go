package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// APIError represents error details in API responses
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta represents metadata for API responses
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Version   string `json:"version,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Total     int    `json:"total,omitempty"`
}

// ResponseFormatter handles API response formatting
type ResponseFormatter struct {
	version string
}

// NewResponseFormatter creates a new response formatter
func NewResponseFormatter(version string) *ResponseFormatter {
	return &ResponseFormatter{
		version: version,
	}
}

// Success sends a successful API response
func (rf *ResponseFormatter) Success(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
		Meta: &Meta{
			Version: rf.version,
		},
	}

	rf.writeResponse(w, http.StatusOK, response)
}

// Created sends a 201 Created response
func (rf *ResponseFormatter) Created(w http.ResponseWriter, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
		Meta: &Meta{
			Version: rf.version,
		},
	}

	rf.writeResponse(w, http.StatusCreated, response)
}

// Error sends an error API response
func (rf *ResponseFormatter) Error(w http.ResponseWriter, statusCode int, code, message, details string) {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now().UTC(),
		Meta: &Meta{
			Version: rf.version,
		},
	}

	rf.writeResponse(w, statusCode, response)
}

// BadRequest sends a 400 Bad Request response
func (rf *ResponseFormatter) BadRequest(w http.ResponseWriter, message string, details string) {
	rf.Error(w, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// Unauthorized sends a 401 Unauthorized response
func (rf *ResponseFormatter) Unauthorized(w http.ResponseWriter, message string) {
	rf.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", message, "")
}

// Forbidden sends a 403 Forbidden response
func (rf *ResponseFormatter) Forbidden(w http.ResponseWriter, message string) {
	rf.Error(w, http.StatusForbidden, "FORBIDDEN", message, "")
}

// NotFound sends a 404 Not Found response
func (rf *ResponseFormatter) NotFound(w http.ResponseWriter, message string) {
	rf.Error(w, http.StatusNotFound, "NOT_FOUND", message, "")
}

// Conflict sends a 409 Conflict response
func (rf *ResponseFormatter) Conflict(w http.ResponseWriter, message string) {
	rf.Error(w, http.StatusConflict, "CONFLICT", message, "")
}

// InternalServerError sends a 500 Internal Server Error response
func (rf *ResponseFormatter) InternalServerError(w http.ResponseWriter, message string) {
	rf.Error(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, "")
}

// ValidationError sends a validation error response
func (rf *ResponseFormatter) ValidationError(w http.ResponseWriter, message string, details string) {
	rf.Error(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message, details)
}

// Paginated sends a paginated response
func (rf *ResponseFormatter) Paginated(w http.ResponseWriter, message string, data interface{}, page, limit, total int) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
		Meta: &Meta{
			Version: rf.version,
			Page:    page,
			Limit:   limit,
			Total:   total,
		},
	}

	rf.writeResponse(w, http.StatusOK, response)
}

// writeResponse writes the JSON response
func (rf *ResponseFormatter) writeResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback error response
		http.Error(w, `{"success":false,"error":{"code":"ENCODING_ERROR","message":"Failed to encode response"}}`, http.StatusInternalServerError)
	}
}

// Global response formatter instance
var DefaultFormatter = NewResponseFormatter("1.0")

// Convenience functions using the default formatter
func Success(w http.ResponseWriter, message string, data interface{}) {
	DefaultFormatter.Success(w, message, data)
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	DefaultFormatter.Created(w, message, data)
}

func BadRequest(w http.ResponseWriter, message string, details string) {
	DefaultFormatter.BadRequest(w, message, details)
}

func Unauthorized(w http.ResponseWriter, message string) {
	DefaultFormatter.Unauthorized(w, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	DefaultFormatter.Forbidden(w, message)
}

func NotFound(w http.ResponseWriter, message string) {
	DefaultFormatter.NotFound(w, message)
}

func Conflict(w http.ResponseWriter, message string) {
	DefaultFormatter.Conflict(w, message)
}

func InternalServerError(w http.ResponseWriter, message string) {
	DefaultFormatter.InternalServerError(w, message)
}

func ValidationError(w http.ResponseWriter, message string, details string) {
	DefaultFormatter.ValidationError(w, message, details)
}
