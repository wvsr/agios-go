package database

import (
	"fmt"
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

		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		if host == "" || user == "" || password == "" || dbname == "" || port == "" || sslmode == "" {
			log.Println("One or more required DB environment variables are missing")
			return
		}

		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode,
		)

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
