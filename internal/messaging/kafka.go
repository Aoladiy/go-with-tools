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
	AuthAdminUserSignedIn        = "auth.admin_user.signed_in"
	AuthAdminUserSignedInGroupId = "main-microservice"
)

type Kafka struct {
	errch                  chan *errs.AppError
	rAuthAdminUserSignedIn *kafka.Reader
}

func New(c config.Config) *Kafka {
	rAuthAdminUserSignedIn := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{c.KafkaAddr},
		Topic:   AuthAdminUserSignedIn,
		GroupID: AuthAdminUserSignedInGroupId,
	})
	return &Kafka{errch: make(chan *errs.AppError, 100), rAuthAdminUserSignedIn: rAuthAdminUserSignedIn}
}

func (k *Kafka) ReadMessages(ctx context.Context) chan bool {
	done := make(chan bool)
	go k.ReadAuthAdminUserSignedInEvent(ctx, k.errch)

	go k.logErrors()
	go k.GracefulShutdown(ctx, done)
	return done
}

func (k *Kafka) GracefulShutdown(ctx context.Context, done chan bool) {
	<-ctx.Done()
	err := k.rAuthAdminUserSignedIn.Close()
	if err != nil {
		slog.Error("cannot close reader rAuthAdminUserSignedIn", "error", err)
	}
	slog.Info("kafka's rAuthAdminUserSignedIn closed successfully")
	done <- true
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
		msg, err := k.rAuthAdminUserSignedIn.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			errch <- errs.Internal(fmt.Errorf("ReadAuthAdminUserSignedInEvent: %w", err))
			continue
		}
		var tmp gen.JWTResponse
		err = proto.Unmarshal(msg.Value, &tmp)
		if err != nil {
			errch <- errs.Internal(fmt.Errorf("ReadAuthAdminUserSignedInEvent: %w", err))
			continue
		}
		slog.Info("successful ReadAuthAdminUserSignedInEvent",
			"accessToken", tmp.AccessToken,
			"refreshToken", tmp.RefreshToken)
		err = k.rAuthAdminUserSignedIn.CommitMessages(ctx, msg)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			errch <- errs.Internal(fmt.Errorf("ReadAuthAdminUserSignedInEvent: %w", err))
		}
	}
}
