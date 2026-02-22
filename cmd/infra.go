package cmd

import (
	"log"

	"github.com/DKeshavarz/peyk/internal/config"
	"github.com/DKeshavarz/peyk/internal/infra/bot"
	"github.com/DKeshavarz/peyk/internal/infra/cache"
	"gopkg.in/telebot.v4"
)

type infra struct {
	telebot *telebot.Bot
	balebot *telebot.Bot
	cache   *cache.Cache
}

func newInfra(cfg *config.Config) *infra {
	telebot, err := bot.New(&cfg.Telebot)
	if err != nil {
		log.Printf("can't create telegram bot: %s\n", err.Error())
	}

	balebot, err := bot.New(&cfg.Balebot)
	if err != nil {
		log.Printf("can't create bale bot: %s\n", err.Error())
	}

	cache := cache.New()

	return &infra{
		telebot: telebot,
		balebot: balebot,
		cache: cache,
	}
}
