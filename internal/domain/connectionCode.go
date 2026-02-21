package domain

import "time"

type ConnectionCode struct {
	Code       string
	SourceChat string
	ExpiresAt  time.Time
}
