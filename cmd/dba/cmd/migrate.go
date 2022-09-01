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

package cmd

import (
	"log"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

var migrationsDir string

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [direction]",
	Short: "Migrates the database",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrates the database to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrate.New(
			"file://"+migrationsDir,
			"postgres://postgres:postgres@localhost:5432/database?sslmode=disable")

		if err != nil {
			log.Fatalf("Error loading migrations: %v", err)
		}

		return m.Up()
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Migrates the database down one version",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrate.New(
			"file://"+migrationsDir,
			"postgres://postgres:postgres@localhost:5432/database?sslmode=disable")

		if err != nil {
			log.Fatalf("Error loading migrations: %v", err)
		}

		return m.Down()
	},
}

var migrateForceCmd = &cobra.Command{
	Use:   "force [version]",
	Short: "Forces the version for the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := strconv.Atoi(args[0])

		if err != nil {
			log.Fatalf("Error reading version input: %v", err)
		}

		m, err := migrate.New(
			"file://"+migrationsDir,
			"postgres://postgres:postgres@localhost:5432/database?sslmode=disable")

		if err != nil {
			log.Fatalf("Error loading migrations: %v", err)
		}

		return m.Force(version)
	},
}
