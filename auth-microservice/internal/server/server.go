package server

import (
	"log"
	"net"
	"strconv"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/auth"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/config"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/database"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/database/queries"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/middleware"
	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Server struct {
	c   config.Config
	db  database.Service
	q   *queries.Queries
	rdb *redis.Client
}

func NewServer(
	c config.Config,
	db database.Service,
	q *queries.Queries,
	rdb *redis.Client,
) *Server {
	return &Server{c: c, db: db, q: q, rdb: rdb}
}

func (s *Server) Serve() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.c.AppPort))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.Logger()))
	defer grpcServer.GracefulStop()
	gen.RegisterAuthMicroserviceServer(grpcServer, auth.New(s.q, s.rdb, s.db.GetPool(), s.c))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
