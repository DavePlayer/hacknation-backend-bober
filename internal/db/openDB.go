package db

import (
	"errors"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("no DATABASE_URL env set")

	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil

}
