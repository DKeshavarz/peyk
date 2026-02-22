package domain

type PlatformStatus string

const (
	PlatformActive   PlatformStatus = "active"
	PlatformInactive PlatformStatus = "inactive"
)

type Platform struct {
	ID     string
	Name   string
	Status PlatformStatus
}

type PlatformName string

const (
	Telegram PlatformName = "telegram"
	Bale     PlatformName = "bale"
)

func (p PlatformName) Valid() bool {
	switch p {
	case Telegram, Bale:
		return true
	}
	return false
}
