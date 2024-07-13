package dbconvert

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func toBool(b bool) pgtype.Bool {
	if !b {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: b, Valid: true}
}

func toDate(t time.Time) pgtype.Date {
	if t.IsZero() {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: t, Valid: true}
}

func toInt8(i int64) pgtype.Int8 {
	if i == 0 {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: i, Valid: true}
}

func toText(str string) pgtype.Text {
	if str == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: str, Valid: true}
}

func toTimestamptz(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}
