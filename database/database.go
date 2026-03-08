package database

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var SqlDB *sql.DB

const DSN = "postgres://postgres:postgres@localhost:5433/crud_api_dev?sslmode=disable"

func Connect() {
	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	DB = db

	// Get raw SQL DB for migrations
	var sqlErr error
	SqlDB, sqlErr = db.DB()
	if sqlErr != nil {
		log.Fatal("Failed to get SQL DB:", sqlErr)
	}

	log.Println("✓ Database connected successfully")
}
