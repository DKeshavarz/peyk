package cmd

import (
	"os"
	"time"

	b "github.com/DKeshavarz/peyk/internal/bot"
	"github.com/DKeshavarz/peyk/internal/config"
	"github.com/DKeshavarz/peyk/internal/domain"

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

		infra := newInfra(cfg)

		codeGen := codegen.NewCodeGenerator()
		codeRepo := cache_repo.NewConnectionCodeRepository(infra.cache)

		connection := service.NewConnectionUsecase(codeGen, codeRepo, 5*time.Minute)
		go b.Start(infra.telebot, connection, domain.Telegram)
		b.Start(infra.balebot, connection, domain.Bale)
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
