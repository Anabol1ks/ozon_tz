package storage

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Store Storage
)

func InitStorage(storageType string) error {
	switch storageType {
	case "memory":
		Store = NewMemoryStorage()
		log.Println("Подключено к memory")
		return nil
	case "postgres":
		if err := ConnectDatabase(); err != nil {
			return err
		}
		Store = NewPostgresStorage(DB)
		log.Println("Подключено к postgres")
		return nil
	default:
		return fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func ConnectDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL не задана, собираем DSN вручную")

		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	DB = db
	log.Println("Подключение к базе данных успешно!")
	return nil
}
