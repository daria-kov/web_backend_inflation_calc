package main

import (
	"log"
	"web_backend/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	log.Println("Application start!")
	api.StartServer()
	log.Println("Application terminated!")
}
