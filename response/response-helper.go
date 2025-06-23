package response

import (
	"encoding/json"
	"net/http"
)

// StandardResponse is the base structure for all API responses
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// WorkoutResponse contains workout-specific response data
type WorkoutResponse struct {
	WorkoutID     int         `json:"workout_id"`
	UpdatedFields interface{} `json:"updated_fields,omitempty"`
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Success sends a success response with the given message and data
func Success(w http.ResponseWriter, message string, data interface{}) {
	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, resp)
}

// Created sends a 201 Created response with the given message and data
func Created(w http.ResponseWriter, message string, data interface{}) {
	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusCreated, resp)
}

// Error sends an error response with the given status code, message, and error
func Error(w http.ResponseWriter, statusCode int, message string, err error) {
	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, statusCode, resp)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   "Resource not found",
	}
	JSON(w, http.StatusNotFound, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, err error) {
	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, http.StatusBadRequest, resp)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string, err error) {
	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, http.StatusInternalServerError, resp)
}

// WorkoutUpdated sends a response for a successfully updated workout
func WorkoutUpdated(w http.ResponseWriter, workoutID int, workout interface{}, updatedFields interface{}) {
	workoutResp := WorkoutResponse{
		WorkoutID:     workoutID,
		UpdatedFields: updatedFields,
	}

	Success(w, "Workout successfully updated", map[string]interface{}{
		"workout_info": workoutResp,
		"workout":      workout,
	})
}

// WorkoutDeleted sends a response for a successfully deleted workout
func WorkoutDeleted(w http.ResponseWriter, workoutID int, workoutInfo interface{}) {
	Success(w, "Workout successfully deleted", map[string]interface{}{
		"workout_id":   workoutID,
		"workout_info": workoutInfo,
	})
}

// WorkoutCreated sends a response for a successfully created workout
func WorkoutCreated(w http.ResponseWriter, workout interface{}) {
	Created(w, "Workout successfully created", workout)
}
