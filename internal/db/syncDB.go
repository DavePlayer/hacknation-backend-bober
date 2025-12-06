package db

import (
	"log"

	"bober.app/models"
)

func SyncDatabase() {
	db, err := OpenDB()
	if err != nil {
		log.Fatal("Faile to open db on migrate %s", err)
		return
	}

	// if err := db.AutoMigrate(&models.Organization{}); err != nil {
	// 	log.Fatal("Failed to migrate model Organization %s", err)
	// 	return
	// }

	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate model Users %s", err)
		return
	}

}
