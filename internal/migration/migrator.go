package migration

import (
	"fmt"
	"os"
	"time"

	"github.com/natkazb/sql-migrator/internal/dbsql" //nolint:depguard
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

const (
	template = `--SQL
-- Up begin
-- SQL statements for applying the migration
-- Up end

-- Down begin
-- SQL statements for rolling back the migration
-- Down end
`
	templateGo = `--GO
-- Up begin
-- GO
-- Up end

-- Down begin
-- GO
-- Down end
`
	format    = "20060102150405"
	FormatSQL = "SQL"
	FormatGO  = "GO"
)

type Migrator struct {
	Path string
	DB   *dbsql.DB
	log  Logger
}

func New(dsn, driver, path string, l Logger) *Migrator {
	return &Migrator{
		Path: path,
		DB:   dbsql.New(dsn, driver, l),
		log:  l,
	}
}

func (m *Migrator) CreateMigration(name, format string) {
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	timestamp := time.Now().Format(format)
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	file, err := os.Create(fmt.Sprintf("%s/%s", m.Path, filename))
	if err != nil {
		m.log.Error(fmt.Sprintf("Error creating file: %s/%s : %s", m.Path, filename, err.Error()))
	}
	defer file.Close()

	template := template
	if format == FormatGO {
		template = templateGo
	}
	_, err = file.WriteString(template)
	if err != nil {
		m.log.Error(fmt.Sprintf("Error writing file: %s/%s : %s", m.Path, filename, err.Error()))
	}
	m.log.Info(fmt.Sprintf("Migration created in %s/%s", m.Path, filename))
}

func (m *Migrator) ApplyMigrations() {
	// Apply all pending migrations
	fmt.Println("Applying migrations...")
	// Implementation here
}

func (m *Migrator) RollbackMigration() {
	// Rollback the last applied migration
	fmt.Println("Rolling back the last migration...")
	// Implementation here
}

func (m *Migrator) RedoMigration() {
	// Redo the last migration
	fmt.Println("Redoing the last migration...")
	// Implementation here
}

func (m *Migrator) ShowStatus() {
	// Show the status of all migrations
	fmt.Println("Migration status:")
	// Implementation here
}

func (m *Migrator) ShowDBVersion() {
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	// номер последней примененной миграции
	info, err := m.DB.ShowLast()
	if err != nil {
		m.log.Error(err.Error())
	}
	m.log.Info(fmt.Sprintf("Current database version: %s", info))
}
