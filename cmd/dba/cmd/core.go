package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&migrationsDir,
		"migrations-dir",
		"m",
		"dba/migrations",
		"location of the migrations directory",
	)
}

var rootCmd = &cobra.Command{
	Use:   "dba",
	Short: "dba is a utility for managing the PostgreSQL database",
}

// Execute is the entrypoint for the CLI app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
