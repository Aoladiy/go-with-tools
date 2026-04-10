package messaging

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/config"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/errs"
	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	AuthAdminUserSignedIn = "auth.admin_user.signed_in"
)

type Kafka struct {
	wAuthAdminUserSignedIn kafka.Writer
}

func New(c config.Config) *Kafka {
	return &Kafka{wAuthAdminUserSignedIn: kafka.Writer{Addr: kafka.TCP(c.KafkaAddr), Topic: AuthAdminUserSignedIn}}
}

func Init(c config.Config) {
	var conn *kafka.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = kafka.Dial("tcp", c.KafkaAddr)
		if err == nil {
			slog.Info("connection to kafka succeed")
			break
		}
		slog.Info(fmt.Sprintf("trying to connect to kafka, attempt №%v", i+1))
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		log.Fatal(err)
	}
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             AuthAdminUserSignedIn,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (k *Kafka) WriteAuthAdminUserSignedInEvent(ctx context.Context, msg *gen.JWTResponse) *errs.AppError {
	kMsg, err := proto.Marshal(msg)
	if err != nil {
		return errs.Internal(err)
	}
	err = k.wAuthAdminUserSignedIn.WriteMessages(ctx, kafka.Message{
		Value: kMsg,
	})
	if err != nil {
		return errs.Internal(err)
	}
	return nil
}
