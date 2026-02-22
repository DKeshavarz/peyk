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
		cache := cache.New()
		if err != nil {
			log.Printf("can't create telegram bot: %s\n", err.Error())
		}

		codeGen := codegen.NewCodeGenerator()
		codeRepo := cache_repo.NewConnectionCodeRepository(cache)
		
		connection := service.NewConnectionUsecase(codeGen, codeRepo, 5*time.Minute)
		b.Start(telebot, connection, domain.Telegram)

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
