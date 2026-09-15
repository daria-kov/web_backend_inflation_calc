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

func (h *Handler) GetConsumerCategories(ctx *gin.Context) {
	filterFrom := ctx.Query("from")
	filterTo := ctx.Query("to")

	categories, err := h.Repository.GetConsumerCategories(filterFrom, filterTo)
	if err != nil {
		logrus.Error(err)
	}

	type CategoryView struct {
		repository.ConsumerCategory
		LikesCount int
	}

	var viewCategories []CategoryView
	for _, c := range categories {
		viewCategories = append(viewCategories, CategoryView{
			ConsumerCategory: c,
			LikesCount:       len(c.Likes),
		})
	}

	ctx.HTML(http.StatusOK, "consumer_categories.html", gin.H{
		"categories": viewCategories,
		"filterFrom": filterFrom,
		"filterTo":   filterTo,
	})
}

func (h *Handler) GetConsumerCategoryFeed(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id := 0
	if idStr != "" {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			logrus.Error("Неверный ID:", err)
			id = 0
		}
	}

	nextParam := ctx.Query("next")
	isNext := (nextParam == "true")

	category, err := h.Repository.GetConsumerCategoryByID(id, isNext)
	if err != nil {
		logrus.Error(err)
	}

	likesCount := len(category.Likes)

	ctx.HTML(http.StatusOK, "consumer_categories_feed.html", gin.H{
		"category":   category,
		"likesCount": likesCount,
	})
}

func (h *Handler) GetAddConsumerCategory(ctx *gin.Context) {
	draft, err := h.Repository.GetConsumerCategoryDraft()
	if err != nil {
		logrus.Error(err)
	}

	draftForView := draft
	draftForView.Date = repository.FormatDateToInput(draft.Date)

	ctx.HTML(http.StatusOK, "add_consumer_category.html", gin.H{
		"draft": draftForView,
	})
}
