package main

import (
	"log"

	router "bober.app/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := router.New()

	router.Run()
}
