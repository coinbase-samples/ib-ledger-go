/**
 * Copyright 2022-present- Present Coinbase Global, Inc.
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
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Port                 string `mapstructure:"PORT"`
	DbCreds              string `mapstructure:"DB_CREDENTIALS"`
	DbPort               string `mapstructure:"DB_PORT"`
	DbHostname           string `mapstructure:"DB_HOSTNAME"`
	Env                  string `mapstructure:"ENV_NAME"`
	LogLevel             string `mapstructure:"LOG_LEVEL"`
	NetworkName          string `mapstructure:"INTERNAL_API_HOSTNAME"`
	CoinbaseUserId       string `mapstructure:"COINBASE_USER_ID"`
	NeoworksUserId       string `mapstructure:"NEOWORKS_USER_ID"`
	CoinbaseUsdAccountId string `mapstructure:"COINBASE_USD_ACCOUNT_ID"`
	NeoworksUsdAccountId string `mapstructure:"NEOWORKS_USD_ACCOUNT_ID"`
	DevRegion            string `mapstructure:"DEV_REGION"`
	QldbName             string `mapstructure:"QLDB_NAME"`
	QueueUrl             string `mapstructure:"QUEUE_URL"`
	LocalStackUrl        string `mapstructure:"LOCAL_STACK_URL"`
	DevProfile           string `mapstructure:"DEV_PROFILE"`
}

func (a AppConfig) IsLocalEnv() bool {
	return a.Env == "local" || a.Env == "awslocal"
}

func Setup(app *AppConfig) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("PORT", "8443")
	viper.SetDefault("GRPC_PORT", "50002")
	viper.SetDefault("ENV_NAME", "local")

	// These credentials are only used locally
	viper.SetDefault("DB_CREDENTIALS", "{\"password\":\"postgres\",\"username\":\"postgres\"}")

	viper.SetDefault("DB_HOSTNAME", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("COINBASE_USER_ID", "c4d0e14e-1b2b-4023-afa6-8891ad1960c9")
	viper.SetDefault("NEOWORKS_USER_ID", "b72d0e55-f53a-4db0-897e-2ce4a73cb94b")
	viper.SetDefault("COINBASE_USD_ACCOUNT_ID", "e8ee70ca-cf90-4e87-b43b-bf8496854333")
	viper.SetDefault("NEOWORKS_USD_ACCOUNT_ID", "ed7aa9d6-8fec-4472-9aa1-d5ff6d2115eb")
	viper.SetDefault("DB_NAME", "LedgerTest")
	viper.SetDefault("DEV_REGION", "us-east-1")
	viper.SetDefault("LOCAL_STACK_URL", "http://localhost:4565")
	viper.SetDefault("DEV_PROFILE", "sa-infra")

	err := viper.ReadInConfig()
	if err != nil {
		log.Debugf("missing env file %v\n", err)
	}

	err = viper.Unmarshal(&app)
	if err != nil {
		log.Debugf("cannot parse env file %v\n", err)
	}
}

func ValidateConfig(a AppConfig, l *log.Entry) {
	if a.Env == "" {
		l.Fatalln("no environment name set")
	}
}
