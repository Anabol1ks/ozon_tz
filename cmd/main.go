package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Anabol1ks/ozon_tz/graph"
	"github.com/Anabol1ks/ozon_tz/graph/model"
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

	resolver := &graph.Resolver{
		DB:               storage.DB,
		CommentObservers: make(map[string][]chan *model.Comment),
	}

	r := gin.Default()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	r.POST("/query", gin.WrapH(srv))
	r.GET("/query", gin.WrapH(srv))

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
