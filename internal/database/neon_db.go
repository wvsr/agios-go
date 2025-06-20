package database

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func ConnectToNeonDB() error {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found")
		}

		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Println("DATABASE_URL not set in .env file")
			return
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("failed to connect to NeonDB: %v", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("failed to get DB object from GORM: %v", err)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			log.Printf("failed to ping NeonDB: %v", err)
			return
		}

		DB = db
		log.Println("Successfully connected to NeonDB with GORM!")
	})

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func CloseNeonDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("NeonDB GORM connection closed.")
		}
	}
}
