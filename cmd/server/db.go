package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func initDB(dsn string) (*sqlx.DB, error) {

	db, err := postgresConnect(dsn)
	if err != nil {
		return nil, err
	}

	if err = migrateFor(db.DB, "postgres"); err != nil {
		return nil, err
	}

	return db, nil
}

func postgresConnect(dsn string) (*sqlx.DB, error) {

	db, errOpen := sqlx.Open("postgres", dsn)
	if errOpen != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", errOpen)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connection to DB created, but Ping returned error: %v", err)
	}

	return db, nil
}

func migrateFor(db *sql.DB, driverDB string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		driverDB, driver)

	if err != nil {
		return err
	}

	if errUp := m.Up(); errUp != nil && errUp != migrate.ErrNoChange {
		return errUp
	}

	return nil
}
