package domain

type Chat struct {
	PlatformID string
	ChatID     string
}

func (c Chat) Equal(other Chat) bool {
	return c.PlatformID == other.PlatformID &&
		c.ChatID == other.ChatID
}
