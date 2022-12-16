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
	"os"

	"github.com/sirupsen/logrus"
)

func LogInit(app AppConfig) *logrus.Entry {
	logger := logrus.New()
	logLevel, _ := logrus.ParseLevel(app.LogLevel)
	logger.SetLevel(logLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetOutput(os.Stdout)
	return logrus.NewEntry(logger)
}
