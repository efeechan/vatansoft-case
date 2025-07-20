package config

import (
	"fmt"
	"log"
	"time"

	"github.com/efecan/vatansoft-case/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		GetEnv("DB_HOST", "localhost"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "postgres"),
		GetEnv("DB_NAME", "vatansoft"),
	)

	var db *gorm.DB
	var err error

	// Retry loop
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Database connected")
			err = db.AutoMigrate(
				&models.Address{},
				&models.Hospital{},
				&models.User{},
				&models.Doctor{},
				&models.DepartmentType{},
				&models.Department{},
				&models.City{},
				&models.District{},
			)
			if err != nil {
				log.Fatalf("⚠️ AutoMigrate failed: %v", err)
			}
			DB = db
			return
		}
		log.Printf("⏳ Waiting for database... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("⚠️ Failed to connect to database: %v", err)
}
