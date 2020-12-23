package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/bot"
	"github.com/caquillo07/rotom-bot/conf"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "bot",
		Short: "Run the Rotom-Bot server",
		Run:   runBotCommand,
	})
}

func runBotCommand(_ *cobra.Command, _ []string) {
	logger := zap.L()
	config, err := conf.LoadConfig(viper.GetViper())
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	b := bot.NewBot(config)
	if err := b.Run(); err != nil {
		logger.Fatal("failed to run the bot", zap.Error(err))
	}
}
