package models

import(
	"log"
	"gorm.io/gorm"
)

func MigrateModels(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &ParkingZone{}, &Reservation{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}