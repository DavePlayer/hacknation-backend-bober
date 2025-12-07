package main

import (
	"log"

	"bober.app/internal/db"
	router "bober.app/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.SyncDatabase()
	router := router.New()
	// router.Use(cors.Default())

	router.Run()
}
