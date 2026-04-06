package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/errs"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		uuid := uuid2.New().String()
		ctx = context.WithValue(ctx, "request_id", uuid)

		start := time.Now()
		resp, err = handler(ctx, req)
		finish := time.Since(start)
		if err == nil {
			return resp, nil
		}

		var appErr *errs.AppError
		if !errors.As(err, &appErr) {
			e := errs.Internal(fmt.Errorf("error is not of type *errs.AppError: %w", err))
			logError(uuid, info.FullMethod, e)
			logDebug(uuid, info.FullMethod, finish)
			return nil, respondError(e)
		}

		if appErr.Code == errs.InternalErrCode {
			logError(uuid, info.FullMethod, appErr.Err)
			logDebug(uuid, info.FullMethod, finish)
			return nil, respondError(appErr)
		}

		logInfo(uuid, info.FullMethod, appErr.Err)
		logDebug(uuid, info.FullMethod, finish)
		return nil, respondError(appErr)
	}
}

func logError(requestId, method string, err error) {
	slog.Error("failed request",
		"request_id", requestId,
		"method", method,
		"err", err,
	)
}

func logInfo(requestId, method string, err error) {
	slog.Info("failed request",
		"request_id", requestId,
		"method", method,
		"err", err,
	)
}

func logDebug(requestId, method string, duration time.Duration) {
	slog.Debug("request finished",
		"request_id", requestId,
		"method", method,
		"duration", duration,
	)
}

func respondError(appErr *errs.AppError) error {
	return status.Error(appErr.GrpcCode(), appErr.Error())
}
