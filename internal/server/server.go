package server

import (
	"context"
	"fmt"
	"go-with-tools/internal/auth"
	"go-with-tools/internal/brand"
	"go-with-tools/internal/category"
	"go-with-tools/internal/config"
	"go-with-tools/internal/database/queries"
	"go-with-tools/internal/product"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-with-tools/internal/database"
)

type Server struct {
	c        config.Config
	db       database.Service
	q        *queries.Queries
	Server   *http.Server
	brand    *brand.Service
	category *category.Service
	product  *product.Service
	auth     *auth.Service
}

func NewServer(c config.Config) *Server {
	db := database.New(c)
	pool := db.GetPool()
	q := queries.New(db.GetPool())
	newServer := &Server{
		c:        c,
		db:       db,
		q:        q,
		brand:    brand.New(q, pool),
		category: category.New(q, pool),
		product:  product.New(q, pool),
		auth:     auth.New(q, pool, c),
	}

	// Declare Server c
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", newServer.c.AppHost, newServer.c.AppPort),
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
