package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Aoladiy/go-with-tools/cmd/api/docs"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/logs"
	"github.com/Aoladiy/go-with-tools/internal/messaging"
	"github.com/Aoladiy/go-with-tools/internal/server"
)

func gracefulShutdown(apiServer *server.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := apiServer.ShutdownServer(ctx)
	if err != nil {
		slog.Info(fmt.Sprintf("Server forced to shutdown with error: %v", err))
	}

	slog.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

//	@title			My API
//	@version		1.0
//	@description	A brief description of my API.

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the JWT token.
func main() {
	c := config.Config{}
	err := c.LoadEnv()
	if err != nil {
		log.Fatalln("cannot load env variables", err)
	}
	logs.Init(c)
	newServer := server.New(c)
	k:= messaging.New(c)
	k.ReadMessages()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(newServer, done)

	newServer.Serve()

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}
