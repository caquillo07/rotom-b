package conf

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Discord struct {
		Token         string
		PlayingStatus string
	}

	Bot struct {
		Prefix          string
		EmbedColor      int
		ErrorEmbedColor int
	}
	TradeCodes struct {
		TradeDesc		string
		Code			string
		TradeDesc2		string
		Code2			string
		TradeDesc3		string
		Code3			string	
		TradeDesc4		string
		Code4			string
		TradeDesc5		string
		Code5			string
		TradeDesc6		string
		Code6			string
		TradeDesc7		string
		Code7			string
		TradeDesc8		string
		Code8			string
		TradeDesc9		string
		Code9			string
		TradeDesc10		string
		Code10			string
	}
}

// LoadConfig loads configuration from the viper instance.
func LoadConfig(v *viper.Viper) (Config, error) {
	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, err
	}
	return config, nil
}

// InitViper initializes the global viper instance with the configuration file
func InitViper(configFile string) {
	if configFile == "" {

		// default to one in present directory named 'config'
		configFile = "config.yaml"
	}
	viper.SetConfigFile(configFile)

	// Default settings
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	zap.L().Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))
}
