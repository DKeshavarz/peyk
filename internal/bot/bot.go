package bot

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

func Start(bot *tele.Bot) {

	bot.Handle("/start", func(c tele.Context) error {
		
		if c.Chat().Type == tele.ChatGroup || c.Chat().Type == tele.ChatSuperGroup {
			return c.Send("Hello everyone in this group! I'm a bot with a single command.")
		}
	
		return c.Send("Hello! Welcome to my bot. I'm a simple bot with one command.")
	})

	
	bot.Handle("/help", func(c tele.Context) error {
		return c.Send("This bot has only one command: /start\nIt works in both private chats and groups!")
	})


	
	log.Println("Bot started successfully!")
	bot.Start()
}
