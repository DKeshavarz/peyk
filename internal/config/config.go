package config

import "github.com/DKeshavarz/peyk/internal/infra/bot"

type Config struct {
	Telebot bot.Config `envPrefix:"TELEBOT_"`
}