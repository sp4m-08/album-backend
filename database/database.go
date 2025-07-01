package database

import (
	"fmt"
	"log"
	"rest/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "postgresql://postgres.tybnqguzgkaugyaosuuz:UbEIACnr0hxHSIRE@aws-0-ap-south-1.pooler.supabase.com:5432/postgres?sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&models.Album{})
	DB = db
	fmt.Println("Database connection established!")
}
