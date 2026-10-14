package api

import (
	"context"
	"log"
	"os"
	"time"
	"web_backend/internal/app/dsn"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/middleware"
	"web_backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "web_backend/docs"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository(dsn.FromEnv())
	if err != nil {
		logrus.Fatal("Ошибка инициализации репозитория:", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logrus.Fatal("Ошибка подключения к Redis:", err)
	}

	h := handler.NewHandler(repo, redisClient)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// Публичные маршруты (токен необязателен, нужен только для is_liked)
	optionalAuth := middleware.OptionalAuthMiddleware(redisClient)

	r.GET("/api/consumer_categories", optionalAuth, h.GetConsumerCategoriesAPI)
	r.GET("/api/consumer_categories/feed", optionalAuth, h.GetConsumerCategoryFeedAPI)
	r.GET("/api/consumer_categories/feed/:id", optionalAuth, h.GetConsumerCategoryFeedByIDAPI)

	r.POST("/api/users/register", h.RegisterUserAPI)
	r.POST("/api/users/auth", h.AuthUserAPI)

	// Маршруты, требующие JWT
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(redisClient))

	protected.GET("/consumer_categories/draft", h.GetConsumerCategoryDraftAPI)
	protected.POST("/consumer_categories", h.CreateConsumerCategoryAPI)
	protected.PUT("/consumer_categories/:id/publish", h.PublishConsumerCategoryAPI)
	protected.POST("/consumer_categories/:id/like", h.LikeConsumerCategoryAPI)
	protected.DELETE("/consumer_categories/:id", h.DeleteConsumerCategoryAPI)

	protected.POST("/users/deauth", h.DeauthUserAPI)

	err = r.Run(":8084")
	if err != nil {
		logrus.Fatal(err)
	}

	log.Println("Server down")
}
