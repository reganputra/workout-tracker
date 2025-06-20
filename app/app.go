package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"workout-tracker/api"
	"workout-tracker/migrations"
	"workout-tracker/store"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	Db             *sql.DB
}

func NewLog() (*Application, error) {

	pgDb, err := store.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	err = store.MigrateFs(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize the WorkoutHandler
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		Db:             pgDb,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
