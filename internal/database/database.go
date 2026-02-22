package database

import (
	"context"
	"fmt"
	"go-with-tools/internal/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	Close()
	GetPool() *pgxpool.Pool
}

type service struct {
	db *pgxpool.Pool
	c  config.Config
}

var dbInstance *service

func New(c config.Config) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbDatabaseName, c.DbSchema)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: pool,
		c:  c,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	return stats
}

func (s *service) Close() {
	s.db.Close()
	log.Printf("Disconnected from database: %s", s.c.DbDatabaseName)
}

func (s *service) GetPool() *pgxpool.Pool {
	return s.db
}
