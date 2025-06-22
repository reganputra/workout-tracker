package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"workout-tracker/store"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	if params == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get workout with id %d: %v", workoutId, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(workout)
	if err != nil {
		http.Error(w, "Failed to encode workout", http.StatusInternalServerError)
		return
	}

}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)
}

func (wh *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	if params == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(params, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get workout with id %d: %v", workoutId, err), http.StatusInternalServerError)
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}
	var updatedWorkout struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedWorkout)
	if err != nil {
		http.Error(w, "Failed to decode workout update", http.StatusBadRequest)
		return
	}
	if updatedWorkout.Title != nil {
		existingWorkout.Title = *updatedWorkout.Title
	}
	if updatedWorkout.Description != nil {
		existingWorkout.Description = *updatedWorkout.Description
	}
	if updatedWorkout.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updatedWorkout.DurationMinutes
	}
	if updatedWorkout.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updatedWorkout.CaloriesBurned
	}
	if updatedWorkout.Entries != nil {
		existingWorkout.Entries = updatedWorkout.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update workout with id %d: %v", workoutId, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingWorkout)
}
