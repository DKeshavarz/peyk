package bot

import (
	"log"

	"github.com/DKeshavarz/peyk/internal/domain"
	"github.com/DKeshavarz/peyk/internal/service"
	tele "gopkg.in/telebot.v4"
)

func Start(bot *tele.Bot, connection service.ConnectionUsecase, platform domain.PlatformName) {
	if bot == nil {
		return
	}
	bot.Handle("/start", func(c tele.Context) error {

		if c.Chat().Type == tele.ChatGroup || c.Chat().Type == tele.ChatSuperGroup {
			return c.Send("سلام من یک ربات ابر شاهکار هستم.\n .")
		}

		return c.Send("شالام")
	})

	h := New(connection, platform)
	bot.Handle("/gen", h.generateCode)

	log.Println("Bot started successfully!")
	bot.Start()
}
