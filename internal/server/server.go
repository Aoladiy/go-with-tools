package server

import (
	"context"
	"fmt"
	"go-with-tools/internal/auth"
	"go-with-tools/internal/brand"
	"go-with-tools/internal/category"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/product"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-with-tools/internal/database"
)

type Server struct {
	port int

	db       database.Service
	q        *queries.Queries
	Server   *http.Server
	brand    *brand.Service
	category *category.Service
	product  *product.Service
	auth     *auth.Service
}

func NewServer() *Server {
	portEnv, exists := os.LookupEnv("PORT")
	if !exists {
		portEnv = "8080"
	}
	port, _ := strconv.Atoi(portEnv)
	db := database.New()
	pool := db.GetPool()
	q := queries.New(db.GetPool())
	newServer := &Server{
		port:     port,
		db:       db,
		q:        q,
		brand:    brand.New(q, pool),
		category: category.New(q, pool),
		product:  product.New(q, pool),
		auth:     auth.New(q, pool),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	newServer.Server = server

	return newServer
}

func (s *Server) ShutdownServer(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	s.db.Close()
	return nil
}
