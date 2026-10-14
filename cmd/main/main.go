package main

import (
	"log"
	"web_backend/internal/api"

	"github.com/joho/godotenv"
)

// @title Consumer Categories API
// @version 1.0
// @description API для работы с пользовательскими категориями.
// @host localhost:8084
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	_ = godotenv.Load()

	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
