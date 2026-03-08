package database

import (
	"fmt"
	"log"
	"tasks-manager/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.Config) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Problem in connection to database")
	}

	if err = db.Ping(); err != nil {
		log.Println("Problem in getting answer from database")
	} else {
		fmt.Println("Database coneection is succeeded")
	}

	return db
}
