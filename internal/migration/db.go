package migration

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	TABLE_NAME        = "migrations"
	STATUS_PROCESSING = "process"
	STATUS_DONE       = "done"
	STATUS_ERROR      = "error"
)

/*type Migrations struct {
	Id      int
	Name    string
	Status  string
	Applied string
} */

type DB struct {
	Dsn    string
	Driver string
	db     *sqlx.DB
}

func New(dsn, driver string) *DB {
	return &DB{
		Dsn:    dsn,
		Driver: driver,
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
	fmt.Println("Migrations table create") // @todo: log
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
		//@todo: записать в лог, вывести пользователю
	}
	return err
}
