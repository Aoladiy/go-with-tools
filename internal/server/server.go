package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Aoladiy/go-with-tools/internal/auth"
	"github.com/Aoladiy/go-with-tools/internal/brand"
	"github.com/Aoladiy/go-with-tools/internal/cache"
	"github.com/Aoladiy/go-with-tools/internal/category"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
	"github.com/Aoladiy/go-with-tools/internal/product"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"

	"github.com/Aoladiy/go-with-tools/internal/database"
)

type Server struct {
	c        config.Config
	db       database.Service
	rdb      *redis.Client
	q        *queries.Queries
	Server   *http.Server
	brand    *brand.Service
	category *category.Service
	product  *product.Service
	auth     *auth.Service
}

func NewServer(c config.Config) *Server {
	db := database.New(c)
	rdb := cache.New(c)
	pool := db.GetPool()
	q := queries.New(db.GetPool())
	newServer := &Server{
		c:        c,
		rdb:      rdb,
		db:       db,
		q:        q,
		brand:    brand.New(q, pool),
		category: category.New(q, pool),
		product:  product.New(q, pool),
		auth:     auth.New(q, rdb, pool, c),
	}

	// Declare Server config
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
