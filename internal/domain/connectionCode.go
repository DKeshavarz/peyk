package domain

import "time"



type ConnectionCode struct {
	Code      string
	ChatID    string
	Platform  PlatformName
	ExpiresAt time.Time
}
