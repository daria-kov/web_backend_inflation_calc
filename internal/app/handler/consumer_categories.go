package handler

import (
	"net/http"
	"strconv"
	"time"
	"web_backend/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CategoryView struct {
	ds.ConsumerCategory
	LikesCount int
}

func (h *Handler) GetConsumerCategories(ctx *gin.Context) {
	filterDate := ctx.Query("date")

	categories, err := h.Repository.GetConsumerCategories(filterDate)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var viewCategories []CategoryView

	for _, c := range categories {
		viewCategories = append(viewCategories, CategoryView{
			ConsumerCategory: c,
			LikesCount:       c.LikesCount,
		})
	}

	ctx.HTML(http.StatusOK, "consumer_categories.html", gin.H{
		"categories": viewCategories,
		"filterDate": filterDate,
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
	isNext := nextParam == "true"

	category, err := h.Repository.GetConsumerCategoryByID(id, isNext)
	if err != nil {
		logrus.Error(err)
		ctx.Redirect(http.StatusFound, "/consumer_categories_feed/1")
		return
	}

	ctx.HTML(http.StatusOK, "consumer_categories_feed.html", gin.H{
		"category":   category,
		"likesCount": category.LikesCount,
	})
}

func (h *Handler) GetAddConsumerCategory(ctx *gin.Context) {
	draft, err := h.Repository.GetConsumerCategoryDraft()
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.HTML(http.StatusOK, "add_consumer_category.html", gin.H{
		"draft": draft,
	})
}

func (h *Handler) CreateConsumerCategoryDraft(ctx *gin.Context) {
	name := ctx.PostForm("name")

	if name == "" {
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	err := h.Repository.CreateConsumerCategoryDraft(name)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
}

func (h *Handler) PublishConsumerCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.PostForm("category_id"))
	if err != nil {
		logrus.Error("Неверный ID категории:", err)
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	description := ctx.PostForm("description")

	historicalDate, err := time.Parse(
		"2006-01-02",
		ctx.PostForm("historical_date"),
	)
	if err != nil {
		logrus.Error("Неверная дата:", err)
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	basketWeight, err := strconv.ParseFloat(
		ctx.PostForm("basket_weight"),
		64,
	)
	if err != nil {
		logrus.Error("Неверный вес:", err)
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	if err != nil {
		logrus.Error("Неверная сумма на дату:", err)
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	if err != nil {
		logrus.Error("Неверная сумма сейчас:", err)
		ctx.Redirect(http.StatusSeeOther, "/add_consumer_category")
		return
	}

	err = h.Repository.PublishConsumerCategory(
		id,
		description,
		historicalDate,
		basketWeight,
	)

	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/consumer_categories")
}

func (h *Handler) DeleteConsumerCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		logrus.Error("Неверный ID категории:", err)
		ctx.Redirect(http.StatusSeeOther, "/consumer_categories")
		return
	}

	err = h.Repository.DeleteConsumerCategory(id)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/consumer_categories")
}
