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
