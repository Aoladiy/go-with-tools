package auth

import (
	"log"
	"strconv"

	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(c config.Config) gen.AuthMicroserviceClient {
	conn, err := grpc.NewClient(c.AuthHost+":"+strconv.Itoa(c.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return gen.NewAuthMicroserviceClient(conn)
}
