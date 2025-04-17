package migration

import (
	"fmt"
	"os"
	"time"

	"github.com/natkazb/sql-migrator/internal/db_migr"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

const (
	TEMPLATE = `--SQL
-- Up begin
-- SQL statements for applying the migration
-- Up end

-- Down begin
-- SQL statements for rolling back the migration
-- Down end
`
	TEMPLATE_GO = `--GO
-- Up begin
-- GO
-- Up end

-- Down begin
-- GO
-- Down end
`
	FORMAT     = "20060102150405"
	FORMAT_SQL = "SQL"
	FORMAT_GO  = "GO"
)

type Migrator struct {
	Path string
	Db   *db_migr.DB
	log  Logger
}

func New(dsn, driver, path string, l Logger) *Migrator {
	return &Migrator{
		Path: path,
		Db:   db_migr.New(dsn, driver, l),
		log:  l,
	}
}

func (m *Migrator) CreateMigration(name, format string) {
	err := m.Db.Init()
	defer m.Db.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	timestamp := time.Now().Format(FORMAT)
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	file, err := os.Create(fmt.Sprintf("%s/%s", m.Path, filename))
	if err != nil {
		m.log.Error(fmt.Sprintf("Error creating file: %s/%s : %s", m.Path, filename, err.Error()))
	}
	defer file.Close()

	template := TEMPLATE
	if format == FORMAT_GO {
		template = TEMPLATE_GO
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
	err := m.Db.Init()
	defer m.Db.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	// номер последней примененной миграции
	info, err := m.Db.ShowLast()
	if err != nil {
		m.log.Error(err.Error())
	}
	m.log.Info(fmt.Sprintf("Current database version: %s", info))
}
