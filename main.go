package main

import (
	"net/http"
	"time"
	"workout-tracker/app"
	"workout-tracker/routes"
)

func main() {
	// Create a new application instance
	application, err := app.NewLog()
	if err != nil {
		panic(err)
	}
	defer application.Db.Close()

	// Start the application
	application.Logger.Println("Application started successfully")

	// Set up the routes
	r := routes.SetupRoutes(application)

	server := &http.Server{
		Addr:         "localhost:1500",
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      r,
	}

	// Start the HTTP server
	err = server.ListenAndServe()
	if err != nil {
		application.Logger.Fatalf("Failed to start server: %v", err)
	}
}
