package config

import "github.com/DKeshavarz/peyk/internal/infra/bot"

type Config struct {
	Telebot bot.Config `envPrefix:"TELEBOT_"`
	Balebot bot.Config `envPrefix:"BALEBOT_"`
}