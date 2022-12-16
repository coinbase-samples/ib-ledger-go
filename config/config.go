/**
 * Copyright 2022 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"fmt"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Port               string `mapstructure:"PORT"`
	DbCreds            string `mapstructure:"DB_CREDENTIALS"`
	DbPort             string `mapstructure:"DB_PORT"`
	DbHostname         string `mapstructure:"DB_HOSTNAME"`
	Env                string `mapstructure:"ENV_NAME"`
	LogLevel           string `mapstructure:"LOG_LEVEL"`
	NetworkName        string `mapstructure:"INTERNAL_API_HOSTNAME"`
	CoinbaseUsdAccount string `mapstructure:"COINBASE_USD_ACCOUNT"`
	NeoworksUsdAccount string `mapstructure:"NEOWORKS_USD_ACCOUNT"`
}

func (a AppConfig) IsLocalEnv() bool {
	return a.Env == "local"
}

func Setup(app *AppConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	// set defaults
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("PORT", "8443")
	viper.SetDefault("GRPC_PORT", "50002")
	viper.SetDefault("ENV_NAME", "local")
	viper.SetDefault("DB_CREDENTIALS", "{\"password\":\"postgres\",\"username\":\"postgres\"}")
	viper.SetDefault("DB_HOSTNAME", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("COINBASE_USD_ACCOUNT", "C4D0E14E-1B2B-4023-AFA6-8891AD1960C9")
	viper.SetDefault("NEOWORKS_USD_ACCOUNT", "B72D0E55-F53A-4DB0-897E-2CE4A73CB94B")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Missing env file %v\n", err)
	}

	err = viper.Unmarshal(&app)
	if err != nil {
		fmt.Printf("Cannot parse env file %v\n", err)
	}
}

func ValidateConfig(a AppConfig, l *log.Entry) {
	if a.Env == "" {
		l.Fatalln("no environment name set")
	}

	if a.CoinbaseUsdAccount == "" || a.NeoworksUsdAccount == "" {
		l.Fatalln("fee accounts not set")
	}
}
