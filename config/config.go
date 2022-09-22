package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port        string `mapstructure:"PORT"`
	DbCreds     string `mapstructure:"DB_CREDENTIALS"`
	DbPort      string `mapstructure:"DB_PORT"`
	DbHostname  string `mapstructure:"DB_HOSTNAME"`
	Env         string `mapstructure:"ENV_NAME"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
	NetworkName string `mapstructure:"INTERNAL_API_HOSTNAME"`
}

func Setup(app *AppConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	// set defaults
	viper.SetDefault("LOG_LEVEL", "warning")
	viper.SetDefault("PORT", "8445")
	viper.SetDefault("GRPC_PORT", "50002")
	viper.SetDefault("ENV_NAME", "local")
	viper.SetDefault("DB_CREDENTIALS", "{\"password\":\"postgres\",\"username\":\"postgres\"}")
	viper.SetDefault("DB_HOSTNAME", "localhost")
	viper.SetDefault("DB_PORT", "5432")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Missing env file %v\n", err)
	}

	err = viper.Unmarshal(&app)
	if err != nil {
		fmt.Printf("Cannot parse env file %v\n", err)
	}
}
