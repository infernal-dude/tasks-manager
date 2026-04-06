package database

import (
	"log"

	"database/sql"

	"github.com/pressly/goose"
)

func RunMigrations(db *sql.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatal(err)
	}
}
