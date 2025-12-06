package db

import (
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_RUL")
	if dsn == "" {
		return nil, errors.New("brak zmiennej środowiskowej DATABASE_URL")

	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil

}
