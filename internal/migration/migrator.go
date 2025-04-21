package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/natkazb/sql-migrator/internal/dbsql" //nolint:depguard
	"github.com/natkazb/sql-migrator/internal/mpath" //nolint:depguard
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type Migrator struct {
	DB    *dbsql.DB
	Mpath *mpath.MigrationPath
	log   Logger
}

func New(dsn, driver, path string, l Logger) *Migrator {
	return &Migrator{
		DB:    dbsql.New(dsn, driver, l),
		Mpath: mpath.New(path, l),
		log:   l,
	}
}

func (m *Migrator) CreateMigration(name, format string) {
	m.Mpath.CreateNew(name, format)
}

func (m *Migrator) ApplyMigrations() {
	m.log.Info("START APPLY")
	defer m.log.Info("FINISH APPLY")
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
		return
	}
	files, err := m.Mpath.GetList()
	if err != nil {
		m.log.Error(err.Error())
		return
	}
	migrationsDB, err := m.DB.GetList()
	if err != nil {
		m.log.Error(err.Error())
		return
	}

	// @todo: что делать с записями в бд, которых нет среди списка файлов?

	parser := &BaseParser{}
	for _, file := range files {
		if slices.Contains(migrationsDB, file) {
			m.log.Info("Skipping already applied migration: " + file)
			continue
		}

		filePath := filepath.Join(m.Mpath.Path, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			m.log.Error(fmt.Sprintf("Failed to read migration file %s: %s", file, err.Error()))
			break
		}

		// Parse migration file
		currMigr, err := parser.Parse(string(content))
		if err != nil {
			m.log.Error(fmt.Sprintf("Failed to parse migration file %s: %s", file, err.Error()))
			break
		}
		// Execute migration
		err = m.DB.ProcessMigrate(file, currMigr.Up)
		if err == nil {
			m.log.Info("Successfully applied migration: " + file)
		} else {
			m.log.Error(fmt.Sprintf("Failed to apply migration file %s: %s", file, err.Error()))
			break
		}
	}
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

func (m *Migrator) ShowStatus(limit int) {
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	info, err := m.DB.ShowStatus(limit)
	if err != nil {
		m.log.Error(err.Error())
	}
	m.log.Info(fmt.Sprintf("Status: %s", info))
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
