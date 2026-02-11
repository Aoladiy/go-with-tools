package helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

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
