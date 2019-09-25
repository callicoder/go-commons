package nullutil

import (
	"database/sql"
	"time"
)

func NewNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
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

func NewPointerString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func NewPointerInt64(i int64) *int64 {
	if i == 0 {
		return nil
	}
	return &i
}

func NewPointerFloat64(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

func NewPointerBool(b bool) *bool {
	if b == false {
		return nil
	}
	return &b
}

func NewPointerTime(t time.Time) *time.Time {
	if (t == time.Time{}) {
		return nil
	}
	return &t
}
