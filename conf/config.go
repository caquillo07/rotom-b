package conf

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/caquillo07/rotom-bot/repository"
)

type Config struct {
	Discord struct {
		Token            string
		PlayingStatus    string
		InviteURL        string
		SupportServerURL string
	}

	DBL struct {
		Enable bool
		Token string
	}

	Bot struct {
		Prefix          string
		EmbedColor      int
		ErrorEmbedColor int
	}

	// Config provides database configuration
	Database repository.Config
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
