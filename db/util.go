package db

import (
	"database/sql"
	"time"
)

func NewNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) != 0,
	}
}

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: i != 0}
}

func NewNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: f != 0}
}

func NewNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: b != false}
}

func NewNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: t != time.Time{}}
}
