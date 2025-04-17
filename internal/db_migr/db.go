package db_migr

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

const (
	TABLE_NAME        = "migrations"
	STATUS_PROCESSING = "process"
	STATUS_DONE       = "done"
	STATUS_ERROR      = "error"
	NO_DATA           = "no any migration has been applied"
)

type Migration struct {
	Id      int
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
		name TEXT NOT NULL,
		status VARCHAR(15) NOT NULL,
		applied_at TIMESTAMP
	)`, TABLE_NAME)
	_, err := d.db.Exec(query)
	d.log.Debug("Migrations table init")
	return err
}

func (d *DB) ProcessMigrate(name, query string) error {
	// добавим новую запись в таблицу миграций
	var id int
	err := d.db.QueryRow(fmt.Sprintf(`
	INSERT INTO %s 
	(name, status, applied_at) 
	VALUES ($1, $2, $3)
	RETURNING id
	`, TABLE_NAME),
		name,
		STATUS_PROCESSING,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return err
	}

	// попробуем выполнить саму миграцию
	_, err = d.db.Exec(query)
	status := STATUS_DONE
	if err != nil {
		status = STATUS_ERROR
	}
	_, errUpdt := d.db.Exec(fmt.Sprintf(`
	UPDATE %s SET 
	status = $2
	applied_at = $3
	WHERE id = $1
	`, TABLE_NAME),
		id,
		status,
		time.Now(),
	)
	if errUpdt != nil {
		d.log.Error(errUpdt.Error())
	}
	return err
}

func (d *DB) ShowLast() (string, error) {
	query := fmt.Sprintf(`
	SELECT id, name, status, applied_at
FROM %s
WHERE status = $1 
ORDER BY applied_at DESC
LIMIT 1`, TABLE_NAME)
	results := make([]Migration, 0)
	err := d.db.Select(&results, query, STATUS_DONE)
	if err != nil {
		d.log.Error(err.Error())
	}
	resultInfo := NO_DATA
	if len(results) > 0 {
		resultInfo = fmt.Sprintf("ID=%d NAME=%s STATUS=%s APPLIED=%s", results[0].Id, results[0].Name, results[0].Status, results[0].Applied)
	}
	return resultInfo, err
}
