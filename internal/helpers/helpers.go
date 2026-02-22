package helpers

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/auth"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func DerefString(pointer *string, defaultValue string) (result string) {
	if pointer != nil {
		return *pointer
	}
	return defaultValue
}

func DerefBool(pointer *bool, defaultValue bool) (result bool) {
	if pointer != nil {
		return *pointer
	}
	return defaultValue
}

func ParsePgTimestamptz(timestamptz pgtype.Timestamptz) (time *time.Time) {
	if timestamptz.Valid {
		time = &timestamptz.Time
	}
	return time
}

func ParsePgInt8(int8 pgtype.Int8) (parsed *int) {
	if int8.Valid {
		tmp := int(int8.Int64)
		parsed = &tmp
	}
	return parsed
}

func ToPgInt8(in *int) (out pgtype.Int8) {
	out = pgtype.Int8{
		Int64: 0,
		Valid: false,
	}
	if in != nil {
		out = pgtype.Int8{
			Int64: int64(*in),
			Valid: true,
		}
	}
	return out
}

func SafeGetUserID(ctx context.Context) (int64, error) {
	val := ctx.Value(auth.UserId)

	if val == nil {
		return 0, errors.New("user_id not found in context")
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("user_id has wrong type: expected int64, got %T", val)
	}

	if userID == 0 {
		return 0, errors.New("user_id is invalid (0)")
	}

	return userID, nil
}
