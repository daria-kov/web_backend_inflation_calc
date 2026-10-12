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
		logrus.Fatal("Ошибка инициализации репозитория:", err)
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/api/consumer_categories", h.GetConsumerCategoriesAPI)
	r.GET("/api/consumer_categories/feed/:id", h.GetConsumerCategoryFeedByIDAPI)
	r.GET("/api/consumer_categories/draft", h.GetConsumerCategoryDraftAPI)

	r.POST("/api/consumer_categories", h.CreateConsumerCategoryAPI)
	r.PUT("/api/consumer_categories/:id/publish", h.PublishConsumerCategoryAPI)
	r.POST("/api/consumer_categories/:id/like", h.LikeConsumerCategoryAPI)
	r.DELETE("/api/consumer_categories/:id", h.DeleteConsumerCategoryAPI)

	r.POST("/api/users/register", h.RegisterUserAPI)
	r.POST("/api/users/auth", h.AuthUserAPI)
	r.POST("/api/users/deauth", h.DeauthUserAPI)

	err = r.Run(":8084")
	if err != nil {
		logrus.Fatal(err)
	}
	log.Println("Server down")
}
