package bot

import (
	"context"
	"fmt"

	"github.com/DKeshavarz/peyk/internal/service"
	tele "gopkg.in/telebot.v4"
)

type Handler struct {
	connection service.ConnectionUsecase
}

func (h *Handler) generateCode(c tele.Context) error {
	ctx := context.Background()
	strID := fmt.Sprintf("%d", c.Chat().ID)
	code, err := h.connection.GenerateCode(ctx, strID)
	if err != nil {
		c.Send("somthing happend %s", err.Error())
	}
	return c.Send(code)
}
