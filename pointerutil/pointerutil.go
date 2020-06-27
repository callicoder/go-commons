package pointerutil

import (
	"time"
)

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
