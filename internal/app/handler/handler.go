package handler

import (
	"net/http"
	"strconv"
	"web_backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// 1. Страница "Плитка" (список всех категорий)
// Роут: GET /grid?filter=...
func (h *Handler) GetGrid(ctx *gin.Context) {
	// Получаем параметр фильтра из URL (например, ?filter=20)
	filterWeight := ctx.Query("filter")

	// Запрашиваем данные из репозитория с учетом фильтра
	categories, err := h.Repository.GetCategories(filterWeight)
	if err != nil {
		logrus.Error(err)
	}

	// Подготовка данных для шаблона: считаем количество лайков для каждой категории
	// В Go шаблонах нельзя вызывать len() внутри цикла удобно, поэтому посчитаем заранее
	type CategoryView struct {
		repository.Category
		LikesCount int // Новое поле для шаблона
	}

	var viewCategories []CategoryView
	for _, c := range categories {
		viewCategories = append(viewCategories, CategoryView{
			Category:   c,
			LikesCount: len(c.Likes), // Вычисляем длину массива лайков
		})
	}

	ctx.HTML(http.StatusOK, "grid.html", gin.H{
		"categories":  viewCategories,
		"filterValue": filterWeight, // Передаем значение обратно, чтобы поле input не пустело
	})
}

// 2. Страница "Лента" (одна категория с видео)
// Роут: GET /feed/:id?next=true
func (h *Handler) GetFeed(ctx *gin.Context) {
	// Получаем ID из URL. Если его нет (вкладка "Лента" без ID), будет пустая строка
	idStr := ctx.Param("id")
	id := 0
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			logrus.Error("Неверный ID:", err)
			// Если ID кривой, считаем как 0 (покажем первую)
			id = 0
		}
	}

	// Проверяем параметр ?next=true
	nextParam := ctx.Query("next")
	isNext := (nextParam == "true")

	// Получаем категорию из репозитория
	category, err := h.Repository.GetCategoryByID(id, isNext)
	if err != nil {
		logrus.Error(err)
		// Если ошибка (например, нет данных), можно редиректить на главную или показать ошибку
		// Но пока просто передадим пустую структуру, шаблон должен это пережить
	}

	// Считаем лайки для текущей страницы
	likesCount := len(category.Likes)

	ctx.HTML(http.StatusOK, "feed.html", gin.H{
		"category":   category,
		"likesCount": likesCount,
	})
}

// 3. Страница "Добавление" (черновик)
// Роут: GET /add
func (h *Handler) GetAdd(ctx *gin.Context) {
	draft, err := h.Repository.GetDraft()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "add.html", gin.H{
		"draft": draft,
	})
}
