package routes

import (
	"github.com/go-chi/chi/v5"
	"workout-tracker/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	routes := chi.NewRouter()

	routes.Get("/health", app.HealthCheck)
	routes.Get("/workouts/{id}", app.WorkoutHandler.GetWorkoutById)
	routes.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	return routes
}
