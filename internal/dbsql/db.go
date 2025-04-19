package dbsql

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Importing the PostgreSQL driver for its side effects
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

const (
	tableName        = "migrations"
	noData           = "no any migration has been applied"
	noDataStatus     = "migration table is empty"
	limitStatus      = 5
	StatusProcessing = "process"
	StatusDone       = "done"
)

type Migration struct {
	ID      int
	Name    string
	Status  string
	Applied time.Time `db:"applied_at"`
}

type DB struct {
	Dsn    string
	Driver string
	db     *sqlx.DB
	log    Logger
}

func New(dsn, driver string, l Logger) *DB {
	return &DB{
		Dsn:    dsn,
		Driver: driver,
		log:    l,
	}
}

func (d *DB) Connect() (err error) {
	d.db, err = sqlx.Connect(d.Driver, d.Dsn)
	return err
}

func (d *DB) Close() error {
	if d.db == nil {
		return nil
	}
	return d.db.Close()
}

func (d *DB) Init() error {
	err := d.Connect()
	if err != nil {
		return err
	}
	err = d.CreateMigrationsTable()
	return err
}

func (d *DB) CreateMigrationsTable() error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id SERIAL PRIMARY KEY,
		name VARCHAR(200) NOT NULL UNIQUE,
		status VARCHAR(15) NOT NULL,
		applied_at TIMESTAMP
	)`, tableName)
	_, err := d.db.Exec(query)
	d.log.Debug("Migrations table init")
	return err
}

// Выполнить миграцию.
func (d *DB) ProcessMigrate(name, query string) error {
	tx, err := d.db.Beginx()
	if err != nil {
		d.log.Debug("Failed to start transaction: " + err.Error())
		return err
	}
	// добавим новую запись в таблицу миграций
	var id int
	err = tx.QueryRow(fmt.Sprintf(`
	INSERT INTO %s 
	(name, status, applied_at) 
	VALUES ($1, $2, $3)
	RETURNING id
	`, tableName),
		name,
		StatusProcessing,
		time.Now(),
	).Scan(&id)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			d.log.Debug("Failed to rollback transaction: " + errTx.Error())
		}
		d.log.Debug("Failed insert into migrations table: " + err.Error())
		return err
	}

	// попробуем выполнить саму миграцию
	_, err = tx.Exec(query)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			d.log.Debug("Failed to rollback transaction: " + errTx.Error())
		}
		d.log.Error(fmt.Sprintf("Failed execute: '%s' : %s", query, err.Error()))
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
	UPDATE %s SET 
	status = $2,
	applied_at = $3
	WHERE id = $1
	`, tableName),
		id,
		StatusDone,
		time.Now(),
	)
	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			d.log.Debug("Failed to rollback transaction: " + errTx.Error())
		}
		d.log.Debug("Failed update migrations table: " + err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		d.log.Error("Failed to commit transaction: " + err.Error())
		return err
	}

	return nil
}

// последняя примененная (statusDone) запись в таблице миграций (tableName).
func (d *DB) ShowLast() (string, error) {
	query := fmt.Sprintf(`
	SELECT id, name, status, applied_at
FROM %s
WHERE status = $1 
ORDER BY applied_at DESC
LIMIT 1`, tableName)
	results := make([]Migration, 0)
	err := d.db.Select(&results, query, StatusDone)
	if err != nil {
		d.log.Error(err.Error())
	}
	resultInfo := noData
	if len(results) > 0 {
		resultInfo = fmt.Sprintf("ID=%d NAME=%s STATUS=%s APPLIED=%s",
			results[0].ID,
			results[0].Name,
			results[0].Status,
			results[0].Applied)
	}
	return resultInfo, err
}

// вывод последних limit записей из tableName.
func (d *DB) ShowStatus(limit int) (string, error) {
	if limit == 0 {
		limit = limitStatus
	}
	query := fmt.Sprintf(`
	SELECT id, name, status, applied_at
FROM %s
ORDER BY applied_at DESC
LIMIT %d`, tableName, limit)
	results := make([]Migration, limit)
	err := d.db.Select(&results, query)
	if err != nil {
		d.log.Error(err.Error())
	}
	resultInfo := noDataStatus
	if len(results) > 0 {
		resultInfo = fmt.Sprintf("%v", results)
	}
	return resultInfo, err
}

func (d *DB) GetList() ([]string, error) {
	// @todo: limit, offset
	query := fmt.Sprintf(`
	SELECT name
FROM %s
ORDER BY applied_at ASC`, tableName)
	results := make([]string, 0)
	err := d.db.Select(&results, query)
	if err != nil {
		d.log.Error(err.Error())
	}
	return results, err
}
