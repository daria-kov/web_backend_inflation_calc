package handler

import (
	"web_backend/internal/app/repository"

	"github.com/redis/go-redis/v9"
)

type Handler struct {
	Repository *repository.Repository
	Redis      *redis.Client
}

func NewHandler(r *repository.Repository, redisClient *redis.Client) *Handler {
	return &Handler{
		Repository: r,
		Redis:      redisClient,
	}
}
