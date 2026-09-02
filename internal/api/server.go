package api

import (
	"log"
	"web_backend/internal/app/handler"
	"web_backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("Ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// --- РОУТИНГ (3 GET запроса по заданию) ---

	// 1. Плитка (список всех карточек) + фильтрация
	// Пример: http://localhost:8081/grid?filter=20  НЕПРАВИЛЬНО!!!
	r.GET("/grid", h.GetGrid)

	// 2. Лента (одна карточка по ID) + кнопка "следующий" (?next=true)
	// Пример: http://localhost:8081/feed/1
	// Пример: http://localhost:8081/feed/1?next=true
	// Для вкладки "Лента" без ID: http://localhost:8081/feed/ (Gin отдаст это как id="", мы обработаем как 0)
	r.GET("/feed/:id", h.GetFeed) // Для случая с ID

	// 3. Добавление (черновик)
	// Пример: http://localhost:8081/add
	r.GET("/add", h.GetAdd)

	// Опционально: редирект с корня на плитку, чтобы при запуске сразу открывался список
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/grid")
	})

	err = r.Run(":8081")
	if err != nil {
		logrus.Fatal(err)
	}
	log.Println("Server down")
}
