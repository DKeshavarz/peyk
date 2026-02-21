package domain

type BridgeStatus string

const (
	BridgeActive   BridgeStatus = "active"
	BridgeDisabled BridgeStatus = "disabled"
)

type Bridge struct {
	ID       string
	SourceID string
	TargetID string
	Status   BridgeStatus
}

func (b *Bridge) Disable() {
	b.Status = BridgeDisabled
}

func (b *Bridge) Enable() {
	b.Status = BridgeActive
}
