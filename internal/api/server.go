package api

import (
	"log"
	"web_backend/internal/app/dsn"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository(dsn.FromEnv())
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/consumer_categories", h.GetConsumerCategories)
	r.GET("/consumer_categories_feed/:id", h.GetConsumerCategoryFeed)
	r.GET("/add_consumer_category", h.GetAddConsumerCategory)

	r.POST("/add_consumer_category", h.CreateConsumerCategoryDraft)
	r.POST("/publish_consumer_category", h.PublishConsumerCategory)
	r.POST("/delete_consumer_category/:id", h.DeleteConsumerCategory)

	err = r.Run(":8082")
	if err != nil {
		logrus.Fatal(err)
	}
	log.Println("Server down")
}
