package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/natkazb/sql-migrator/internal/migrator"
)

var rootCmd = &cobra.Command{
	Use:   "gomigrator",
	Short: "Database migration tool",
}

func Execute() {
	rootCmd.AddCommand(createCmd, upCmd, downCmd, redoCmd, statusCmd, dbVersionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Migration name is required")
			return
		}
		migrator.CreateMigration(args[0])
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ApplyMigrations()
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.RollbackMigration()
	},
}

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redo the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.RedoMigration()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ShowStatus()
	},
}

var dbVersionCmd = &cobra.Command{
	Use:   "dbversion",
	Short: "Show the current database version",
	Run: func(cmd *cobra.Command, args []string) {
		migrator.ShowDBVersion()
	},
}