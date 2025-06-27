package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"workout-tracker/middleware"
	"workout-tracker/response"
	"workout-tracker/store"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	if params == "" {
		response.NotFound(w, "Workout ID is required")
		return
	}

	workoutId, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		response.NotFound(w, "Invalid workout ID format")
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to get workout with ID %d", workoutId), err)
		return
	}

	if workout == nil {
		response.NotFound(w, fmt.Sprintf("Workout with ID %d not found", workoutId))
		return
	}

	response.Success(w, "Workout retrieved successfully", workout)

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		response.BadRequest(w, "Failed to decode workout data", err)
		return
	}

	currenUser := middleware.GetUser(r)
	wh.logger.Printf("User from context: %+v", currenUser)
	if currenUser == nil || currenUser == store.AnonymousUser {
		wh.logger.Println("No authenticated user found")
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}
	workout.UserId = currenUser.Id

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		response.InternalServerError(w, "Failed to create workout", err)
		return
	}

	response.WorkoutCreated(w, createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	if params == "" {
		response.NotFound(w, "Workout ID is required")
		return
	}

	workoutId, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		response.NotFound(w, "Invalid workout ID format")
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to get workout with ID %d", workoutId), err)
		return
	}

	if existingWorkout == nil {
		response.NotFound(w, fmt.Sprintf("Workout with ID %d not found", workoutId))
		return
	}

	// Store original values for comparison
	//originalWorkout := *existingWorkout

	var updatedWorkout struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updatedWorkout)
	if err != nil {
		response.BadRequest(w, "Failed to decode workout update data", err)
		return
	}

	// Track what fields were updated
	updatedFields := make(map[string]interface{})

	if updatedWorkout.Title != nil {
		existingWorkout.Title = *updatedWorkout.Title
		updatedFields["title"] = *updatedWorkout.Title
	}

	if updatedWorkout.Description != nil {
		existingWorkout.Description = *updatedWorkout.Description
		updatedFields["description"] = *updatedWorkout.Description
	}

	if updatedWorkout.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkout.DurationMinutes
		updatedFields["duration"] = *updatedWorkout.DurationMinutes
	}

	if updatedWorkout.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkout.CaloriesBurned
		updatedFields["calories_burned"] = *updatedWorkout.CaloriesBurned
	}

	if updatedWorkout.Entries != nil {
		existingWorkout.Entries = updatedWorkout.Entries
		updatedFields["entries"] = "Updated workout entries"
	}

	currenUser := middleware.GetUser(r)
	if currenUser == nil || currenUser == store.AnonymousUser {
		response.BadRequest(w, "User must be login to update a workout", nil)
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to get workout owner for ID %d", workoutId), err)
		return
	}

	if workoutOwner != currenUser.Id {
		response.Forbidden(w, fmt.Sprintf("User %d is not authorized to update workout %d", currenUser.Id, workoutId))
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to update workout with ID %d", workoutId), err)
		return
	}

	response.WorkoutUpdated(w, existingWorkout.Id, existingWorkout, updatedFields)
}

func (wh *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	if params == "" {
		response.NotFound(w, "Workout ID is required")
		return
	}

	workoutId, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		response.NotFound(w, "Invalid workout ID format")
		return
	}

	// Check if the workout exists
	currenUser := middleware.GetUser(r)
	if currenUser == nil || currenUser == store.AnonymousUser {
		response.BadRequest(w, "User must be login to delete a workout", nil)
		return
	}
	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to get workout owner for ID %d", workoutId), err)
		return
	}
	if workoutOwner != currenUser.Id {
		response.Forbidden(w, fmt.Sprintf("User %d is not authorized to delete workout %d", currenUser.Id, workoutId))
		return
	}

	// Get workout details before deletion for the response
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to get workout with ID %d", workoutId), err)
		return
	}

	if workout == nil {
		response.NotFound(w, fmt.Sprintf("Workout with ID %d not found", workoutId))
		return
	}

	// Store basic info for response
	workoutInfo := struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}{
		ID:    workout.Id,
		Title: workout.Title,
	}

	// Delete the workout
	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		response.InternalServerError(w, fmt.Sprintf("Failed to delete workout with ID %d", workoutId), err)
		return
	}

	response.WorkoutDeleted(w, workout.Id, workoutInfo)
}
