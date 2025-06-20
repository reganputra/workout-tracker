package store

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")
	if err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}
	fmt.Println("Connected to the database")

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, nil
}

func MigrateFs(db *sql.DB, migration fs.FS, dir string) error {
	goose.SetBaseFS(migration)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}
