package domain

import "time"

type PlatformName string

const (
	Telegram PlatformName = "telegram"
	Bale     PlatformName = "bale"
)

type ConnectionCode struct {
	Code      string
	ChatID    string
	Platform  PlatformName
	ExpiresAt time.Time
}
