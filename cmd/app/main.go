package main

import (
	"car-rental-application/config"
	"car-rental-application/internal/routes"
	"car-rental-application/pkg"
	"context"
	"errors"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Setup logger
	pkg.SetupLogger()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	// xendit client
	config.InitXendit()

	// initialize database connection
	db, err := config.InitDB()
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}

	// Close db when main function ends
	pgSql, err := db.DB()
	if err != nil {
		logrus.Fatalf("Error getting database object: %v", err)
	}
	defer func() {
		if err := pgSql.Close(); err != nil {
			logrus.Errorf("Error closing database: %v", err)
		}
	}()

	// initialize controller
	c := routes.InitController(db, os.Getenv("JWT_SECRET"))
	// initialize echo
	e := echo.New()

	// middleware
	e.Use(pkg.LogrusLogger)
	e.Use(middleware.Recover())

	// register routes
	routes.RegisterRoutes(e, c, []byte(os.Getenv("JWT_SECRET")))
	startServer(e)
}

func startServer(e *echo.Echo) {
	// Get the server port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define the server configuration
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      e,
	}

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Start the server in a goroutine
	go func() {
		logrus.Infof("Starting server on port %s...", port)
		if err := e.StartServer(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Could not listen on port %s: %v\n", port, err)
		}
	}()

	// Block until we receive a signal
	<-quit

	// Graceful shutdown
	logrus.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	} else {
		logrus.Info("Server exiting")
	}
}
