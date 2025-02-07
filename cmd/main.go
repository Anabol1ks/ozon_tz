package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Anabol1ks/ozon_tz/internal/models"
	"github.com/Anabol1ks/ozon_tz/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	key := os.Getenv("DB_PORT")
	if key == "" {
		log.Println("\nПеременной среды нет, используется .env")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Ошибка загрузки .env файла")
		}
	}

	storage.ConnectDatabase()

	if err := storage.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "GraphQL Posts API"})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
