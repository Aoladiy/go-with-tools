package auth

import (
	"log"
	"strconv"
	"time"

	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	accessExp  = time.Minute * 15
	refreshExp = time.Hour * 24 * 7
	SignedOut  = "signed-out-token-"
)

func NewClient(c config.Config) gen.AuthMicroserviceClient {
	conn, err := grpc.NewClient(c.AuthHost+":"+strconv.Itoa(c.AuthPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return gen.NewAuthMicroserviceClient(conn)
}
