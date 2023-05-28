package configs

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

type envConfigs struct {
	DBAddress  string `mapstructure:"DB_ADDR"`
	DBPassword string `mapstructure:"DB_PASS"`
	AppPort    string `mapstructure:"APP_PORT`
	Domain     string `mapstructure:"DOMAIN"`
	APIQuota   int    `mapstructure:"API_QUOTA"`
}

func loadEnvVariables() (config *envConfigs) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
