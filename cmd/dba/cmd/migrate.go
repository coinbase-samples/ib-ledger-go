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
