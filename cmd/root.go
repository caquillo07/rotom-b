package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/conf"
)

var cfgFile string

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	Use:   "den-bot",
	Short: "Den bot server",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// OnInitialize tells cobra what functions to run when it starts up.
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(func() { conf.InitViper(cfgFile) })

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().Bool("dev-log", false, "Development logging")
}

// initConfig reads in config file and ENV variables if set.
func initLogging() {
	var logger *zap.Logger
	if val, _ := rootCmd.PersistentFlags().GetBool("dev-log"); val {
		logger, _ = zap.NewDevelopment()
		logger.Info("Development logging enabled")
	} else {
		logger, _ = zap.NewProduction()
	}

	logger.Info("Starting Den-Bot")
	zap.ReplaceGlobals(logger)
}
