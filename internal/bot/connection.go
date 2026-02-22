package bot

import (
	"context"
	"fmt"

	"github.com/DKeshavarz/peyk/internal/domain"
	"github.com/DKeshavarz/peyk/internal/service"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	platform   domain.PlatformName
	connection service.ConnectionUsecase
}

func New(connection service.ConnectionUsecase, platform domain.PlatformName) *Handler {
	return &Handler{
		connection: connection,
		platform:   platform,
	}
}

func (h *Handler) generateCode(c tele.Context) error {
	ctx := context.Background()
	strID := fmt.Sprintf("%d", c.Chat().ID)
	code, err := h.connection.GenerateCode(ctx, strID, h.platform)

	msg := `کد موقت گروه
	%s`


	if err != nil {
		c.Send("somthing happend %s", err.Error())
	}
	return c.Send(fmt.Sprintf(msg, code))
}
