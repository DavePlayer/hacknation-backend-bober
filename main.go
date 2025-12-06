package main

import (
	router "bober.app/routes"
)

func main() {
	router := router.New()

	router.Run()
}
