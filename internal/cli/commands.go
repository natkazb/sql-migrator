package cli

import (
	"fmt"
	"os"

	"github.com/natkazb/sql-migrator/internal/config"
	"github.com/natkazb/sql-migrator/internal/logger"
	"github.com/natkazb/sql-migrator/internal/migration"
	"github.com/spf13/cobra"
)

var m *migration.Migrator
var l *logger.Logger
var configFile string

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "gomigrator",
		Short: "Database migration tool",
	}
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to the configuration file")
	rootCmd.AddCommand(createCmd, createCmdGo, upCmd, downCmd, redoCmd, statusCmd, dbVersionCmd)

	cobra.OnInitialize(loadConfig)

	if err := rootCmd.Execute(); err != nil {
		l.Error(err.Error())
	}
}

func loadConfig() {
	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't parse config file, %v", err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Sql.Host,
		conf.Sql.Port,
		conf.Sql.Username,
		conf.Sql.Password,
		conf.Sql.DBName)

	l = logger.New(conf.Logger.Level)
	m = migration.New(dsn, conf.Sql.Driver, conf.Path, l)
}

var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new migration by SQL format",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			l.Error("Migration name is required")
			return
		}
		m.CreateMigration(args[0], migration.FORMAT_SQL)
	},
}

var createCmdGo = &cobra.Command{
	Use:   "create-go <name>",
	Short: "Create a new migration by GO format",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			l.Error("Migration name is required")
			return
		}
		m.CreateMigration(args[0], migration.FORMAT_GO)
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m.ApplyMigrations()
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		m.RollbackMigration()
	},
}

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redo the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		m.RedoMigration()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		m.ShowStatus()
	},
}

var dbVersionCmd = &cobra.Command{
	Use:   "dbversion",
	Short: "Show the current database version",
	Run: func(cmd *cobra.Command, args []string) {
		m.ShowDBVersion()
	},
}
