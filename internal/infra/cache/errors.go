package cache

import "errors"

var (
	ErrInvalidTime = errors.New("time is invalid")
	ErrNotFound    = errors.New("not found")
	ErrDelete      = errors.New("delete faild")
)
