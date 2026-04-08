package database

import (
	"log"

	"database/sql"

	// Для чего здесь пакет с миграциями? В таком простом проекте в нём нет нужды.
	// Тебе достаточно описать актуальную схему данных и создавать её при каждом запуске контейнера в докере
	// Почитай - https://hub.docker.com/_/postgres#initialization-scripts
	// К слову, в проекте нигде не вижу схемы чтобы его запустить.
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
