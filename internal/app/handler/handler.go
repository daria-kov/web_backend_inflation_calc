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

func (h *Handler) GetPlanets(ctx *gin.Context) {
	var planets []repository.Planet
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		planets, err = h.Repository.GetPlanets()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		planets, err = h.Repository.GetPlanetByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	cartDevices, err := h.Repository.GetCart()
	CartCount := 0
	if err == nil {
		CartCount = len(cartDevices)
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"planets":   planets,
		"query":     searchQuery,
		"CartCount": CartCount,
	})
}

func (h *Handler) GetPlanet(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	planet, err := h.Repository.GetPlanet(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"planet": planet,
	})
}

func (h *Handler) GetCart(ctx *gin.Context) {
	var planets []repository.Planet
	var err error

	planets, err = h.Repository.GetCart()
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"planets": planets,
	})
}
