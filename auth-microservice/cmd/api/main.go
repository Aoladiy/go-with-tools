package main

import (
	"log"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/cache"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/config"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/database"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/database/queries"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/server"
)

func main() {
	c := config.Config{}
	err := c.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	rdb := cache.New(c)
	db := database.New(c)
	q := queries.New(db.GetPool())
	s := server.NewServer(c, db, q, rdb)
	s.Serve()
}
