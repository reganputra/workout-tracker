package routes

import (
	"github.com/go-chi/chi/v5"
	"workout-tracker/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	routes := chi.NewRouter()

	routes.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		routes.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
		routes.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		routes.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkout))
		routes.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
	})

	routes.Get("/health", app.HealthCheck)
	routes.Post("/users", app.UserHandler.HandleRegisterUser)
	routes.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	return routes
}
