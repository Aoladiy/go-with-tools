package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/Aoladiy/go-with-tools/internal/auth"
	"github.com/Aoladiy/go-with-tools/internal/brand"
	"github.com/Aoladiy/go-with-tools/internal/category"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/database"
	"github.com/Aoladiy/go-with-tools/internal/database/queries"
	"github.com/Aoladiy/go-with-tools/internal/metrics"
	"github.com/Aoladiy/go-with-tools/internal/product"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	c             config.Config
	db            database.Service
	q             *queries.Queries
	Server        *http.Server
	metricsServer *http.Server
	brand         *brand.Service
	category      *category.Service
	product       *product.Service
	auth          gen.AuthMicroserviceClient
}

func New(c config.Config) *Server {
	metricsServer := metrics.NewServer(c)
	db := database.New(c)
	pool := db.GetPool()
	q := queries.New(db.GetPool())
	newServer := &Server{
		c:             c,
		db:            db,
		q:             q,
		metricsServer: metricsServer,
		brand:         brand.New(q, pool),
		category:      category.New(q, pool),
		product:       product.New(q, pool),
		auth:          auth.NewClient(c),
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
	if err := s.metricsServer.Shutdown(ctx); err != nil {
		return err
	}
	s.db.Close()
	return nil
}

func (s *Server) Serve() {
	go func() {
		err := s.metricsServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("http metrics server error: %s", err))
		}
	}()
	err := s.Server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http newServer error: %s", err))
	}
}
