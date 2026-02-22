package cmd

import (
	"log"
	"os"
	"time"

	b "github.com/DKeshavarz/peyk/internal/bot"
	"github.com/DKeshavarz/peyk/internal/config"
	"github.com/DKeshavarz/peyk/internal/domain"
	"github.com/DKeshavarz/peyk/internal/infra/bot"
	"github.com/DKeshavarz/peyk/internal/infra/cache"

	cache_repo "github.com/DKeshavarz/peyk/internal/repositories/cache"
	"github.com/DKeshavarz/peyk/internal/service"
	"github.com/DKeshavarz/peyk/pkg/codegen"
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
		balebot, err := bot.New(&cfg.Balebot)
		if err != nil {
			log.Printf("can't create bale bot: %s\n", err.Error())
		}
		cache := cache.New()

		codeGen := codegen.NewCodeGenerator()
		codeRepo := cache_repo.NewConnectionCodeRepository(cache)
		
		connection := service.NewConnectionUsecase(codeGen, codeRepo, 5*time.Minute)
		go b.Start(telebot, connection, domain.Telegram)
		b.Start(balebot, connection, domain.Bale)
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
