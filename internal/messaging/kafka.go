package messaging

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/errs"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

const (
	AuthAdminUserSignedIn = "auth.admin_user.signed_in"
)

type Kafka struct {
	errch                  chan *errs.AppError
	rAuthAdminUserSignedIn *kafka.Reader
}

func New(c config.Config) *Kafka {
	rAuthAdminUserSignedIn := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.KafkaAddr},
		Topic:   AuthAdminUserSignedIn,
	})
	return &Kafka{errch: make(chan *errs.AppError), rAuthAdminUserSignedIn: rAuthAdminUserSignedIn}
}

func (k *Kafka) ReadMessages() {
	go k.ReadAuthAdminUserSignedInEvent(context.Background(), k.errch)

	go k.logErrors()
}

func (k *Kafka) logErrors() {
	for {
		err := <-k.errch
		slog.Error("failed kafka message",
			"error", err)
	}
}

func (k *Kafka) ReadAuthAdminUserSignedInEvent(ctx context.Context, errch chan *errs.AppError) {
	for {
		msg, err := k.rAuthAdminUserSignedIn.ReadMessage(ctx)
		if err != nil {
			errch <- errs.Internal(fmt.Errorf("ReadAuthAdminUserSignedInEvent: %w", err))
			continue
		}
		var tmp gen.JWTResponse
		err = proto.Unmarshal(msg.Value, &tmp)
		if err != nil {
			errch <- errs.Internal(fmt.Errorf("ReadAuthAdminUserSignedInEvent: %w", err))
			continue
		}
		if err == nil {
			slog.Info("successful ReadAuthAdminUserSignedInEvent",
				"accessToken", tmp.AccessToken,
				"refreshToken", tmp.RefreshToken)
		}
	}
}
