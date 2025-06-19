package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type WorkoutHandler struct{}

func NewWorkoutHandler() *WorkoutHandler {
	return &WorkoutHandler{}
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

func (wh *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new workout")
}
