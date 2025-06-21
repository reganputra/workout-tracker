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

func (wh *WorkoutHandler) GetWorkoutById(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "WorkoutHandler with id %d", workoutId)
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
