package cmd

import (
	"log"
	"os"

	b "github.com/DKeshavarz/peyk/internal/bot"
	"github.com/DKeshavarz/peyk/internal/config"
	"github.com/DKeshavarz/peyk/internal/infra/bot"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "peyk",
	Short: "A platform that connects different messaging apps together",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.New()

		telebot, err := bot.New(&cfg.Telebot)
		if err != nil {
			log.Printf("can't create telegram bot: %s\n", err.Error())
		}

		b.Start(telebot)
		
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
