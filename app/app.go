package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"workout-tracker/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

func NewLog() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize the WorkoutHandler
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
