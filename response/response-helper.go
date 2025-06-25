package response

import (
	"encoding/json"
	"log"
	"net/http"
)

var logger *log.Logger

func InitLogger(l *log.Logger) {
	logger = l
}

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

type UserResponse struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

// JSON sends a JSON response with the given status code
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Log the response if logger is initialized
	if logger != nil {
		logData, _ := json.Marshal(data)
		logger.Printf("Sending response [%d]: %s", statusCode, string(logData))
	}

	json.NewEncoder(w).Encode(data)
}

// Success sends a success response with the given message and data
func Success(w http.ResponseWriter, message string, data interface{}) {
	if logger != nil {
		logger.Printf("Success: %s", message)
	}

	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusOK, resp)
}

// Created sends a 201 Created response with the given message and data
func Created(w http.ResponseWriter, message string, data interface{}) {
	if logger != nil {
		logger.Printf("Created: %s", message)
	}

	resp := StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSON(w, http.StatusCreated, resp)
}

// Error sends an error response with the given status code, message, and error
func Error(w http.ResponseWriter, statusCode int, message string, err error) {
	if logger != nil {
		logger.Printf("Error [%d]: %s - %v", statusCode, message, err)
	}

	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, statusCode, resp)
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string) {
	if logger != nil {
		logger.Printf("NotFound: %s", message)
	}

	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   "Resource not found",
	}
	JSON(w, http.StatusNotFound, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, err error) {
	if logger != nil {
		logger.Printf("BadRequest: %s - %v", message, err)
	}

	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, http.StatusBadRequest, resp)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string, err error) {
	if logger != nil {
		logger.Printf("InternalServerError: %s - %v", message, err)
	}

	resp := ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSON(w, http.StatusInternalServerError, resp)
}

// WorkoutUpdated sends a response for a successfully updated workout
func WorkoutUpdated(w http.ResponseWriter, workoutID int, workout interface{}, updatedFields interface{}) {
	if logger != nil {
		logger.Printf("Workout updated: ID=%d, Fields=%v", workoutID, updatedFields)
	}

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
	if logger != nil {
		logger.Printf("Workout deleted: ID=%d", workoutID)
	}

	Success(w, "Workout successfully deleted", map[string]interface{}{
		"workout_id":   workoutID,
		"workout_info": workoutInfo,
	})
}

// WorkoutCreated sends a response for a successfully created workout
func WorkoutCreated(w http.ResponseWriter, workout interface{}) {
	if logger != nil {
		logger.Printf("Workout created")
	}

	Created(w, "Workout successfully created", workout)
}

func UserCreated(w http.ResponseWriter, user interface{}) {
	if logger != nil {
		logger.Printf("User created")
	}
	Created(w, "User successfully created", user)
}

func UserUpdated(w http.ResponseWriter, user interface{}) {
	if logger != nil {
		logger.Printf("User updated")
	}
	Success(w, "User successfully updated", user)
}
