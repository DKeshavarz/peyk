package bot

import (
	"log"

	tele "gopkg.in/telebot.v4"
)

type Config struct {
	API   string `env:"API"`
	Token string `env:"TOKEN"`
}

func New(cfg *Config) (bot *tele.Bot, err error) {
	pref := tele.Settings{
		Token: cfg.Token,
	}
	if cfg.API != "" {
		pref.URL = cfg.API
	}

	bot, err = tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized bot %s", bot.Me.Username)
	return bot, err
}