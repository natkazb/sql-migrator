package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/natkazb/sql-migrator/internal/dbsql" //nolint:depguard
	"github.com/natkazb/sql-migrator/internal/mpath" //nolint:depguard
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
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	timestamp := time.Now().Format(format)
	filename := fmt.Sprintf("%s_%s.sql", timestamp, name)
	file, err := os.Create(fmt.Sprintf("%s/%s", m.Mpath.Path, filename))
	if err != nil {
		m.log.Error(fmt.Sprintf("Error creating file: %s/%s : %s", m.Mpath.Path, filename, err.Error()))
	}
	defer file.Close()

	template := template
	if format == FormatGO {
		template = templateGo
	}
	_, err = file.WriteString(template)
	if err != nil {
		m.log.Error(fmt.Sprintf("Error writing file: %s/%s : %s", m.Mpath.Path, filename, err.Error()))
	}
	m.log.Info(fmt.Sprintf("Migration created in %s/%s", m.Mpath.Path, filename))
}

func (m *Migrator) ApplyMigrations() {
	m.log.Info("START APPLY")
	err := m.DB.Init()
	defer m.DB.Close()
	if err != nil {
		m.log.Error(err.Error())
	}
	files, err := m.Mpath.GetList()
	if err != nil {
		m.log.Error(err.Error())
	}
	migrationsDone, err := m.DB.GetListDone()
	if err != nil {
		m.log.Error(err.Error())
	}
	migrationsError, err := m.DB.GetListError()
	if err != nil {
		m.log.Error(err.Error())
	}

	// Apply each migration file
	for _, file := range files {
		// Check if migration is already applied
		if stringInSlice(file, migrationsDone) {
			m.log.Info("Skipping already applied migration: " + file)
			continue
		}
		if stringInSlice(file, migrationsError) {
			// @todo: выполнить миграцию, обновить статус в бд
			continue
		}
		// @todo: что делать с записями в бд, которых нет среди списка файлов?

		// это новая миграция, выполняем, делаем новую запись в бд
		filePath := filepath.Join(m.Mpath.Path, file)
		content, err := os.ReadFile(filePath)
		if err != nil {
			m.log.Error("Failed to read migration file: " + err.Error())
		}

		// Execute migration
		err = m.DB.ProcessMigrate(file, string(content))
		if err == nil {
			m.log.Info("Successfully applied migration: " + file)
		} else {
			m.log.Error("Failed to apply migration: " + err.Error())
		}
	}
	m.log.Info("FINISH APPLY")
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
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
