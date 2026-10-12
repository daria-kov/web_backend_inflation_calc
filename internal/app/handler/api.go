package handler

import (
	"crypto/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/user"

	"github.com/gin-gonic/gin"
)

func generateFilename(ext string) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 16)

	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	for i := range result {
		result[i] = letters[int(randomBytes[i])%len(letters)]
	}

	return string(result) + ext, nil
}

func (h *Handler) GetConsumerCategoriesAPI(ctx *gin.Context) {
	filterDate := ctx.Query("date")
	currentUser := user.GetCurrentUser()

	categories, err := h.Repository.GetConsumerCategoriesForUser(
		filterDate,
		currentUser.UserID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]serializer.ConsumerCategoryResponse, 0, len(categories))

	for _, category := range categories {
		isCreator := 0

		if category.CreatorID == currentUser.UserID {
			isCreator = 1
		}

		response = append(response, serializer.ConsumerCategoryResponse{
			ID:             category.CategoryID,
			Title:          category.Title,
			Description:    category.Description,
			Status:         category.Status,
			ImageURL:       category.ImageURL,
			VideoURL:       category.VideoURL,
			BasketWeight:   category.BasketWeight,
			HistoricalDate: category.HistoricalDate,
			CreatedAt:      category.CreatedAt,
			PublishedAt:    category.PublishedAt,
			CreatorID:      category.CreatorID,
			LikesCount:     category.LikesCount,
			IsCreator:      isCreator,
			IsLiked:        category.IsLiked,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) CreateConsumerCategoryAPI(ctx *gin.Context) {
	title := ctx.PostForm("title")

	if strings.TrimSpace(title) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "title is required",
		})
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "image is required",
		})
		return
	}

	video, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "video is required",
		})
		return
	}

	currentUser := user.GetCurrentUser()

	imageExt := filepath.Ext(image.Filename)
	videoExt := filepath.Ext(video.Filename)

	imageFilename, err := generateFilename(imageExt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	videoFilename, err := generateFilename(videoExt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	imageFile, err := image.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer imageFile.Close()

	err = h.Repository.UploadFile(
		imageFile,
		image.Size,
		imageFilename,
		image.Header.Get("Content-Type"),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	videoFile, err := video.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer videoFile.Close()

	err = h.Repository.UploadFile(
		videoFile,
		video.Size,
		videoFilename,
		video.Header.Get("Content-Type"),
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.Repository.CreateConsumerCategoryDraft(
		title,
		imageFilename,
		videoFilename,
		currentUser.UserID,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "draft created",
	})
}

func (h *Handler) GetConsumerCategoryDraftAPI(ctx *gin.Context) {
	currentUser := user.GetCurrentUser()

	category, err := h.Repository.GetConsumerCategoryDraft(
		currentUser.UserID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if category.CategoryID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "draft not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, serializer.ConsumerCategoryResponse{
		ID:             category.CategoryID,
		Title:          category.Title,
		Description:    category.Description,
		Status:         category.Status,
		ImageURL:       category.ImageURL,
		VideoURL:       category.VideoURL,
		BasketWeight:   category.BasketWeight,
		HistoricalDate: category.HistoricalDate,
		CreatedAt:      category.CreatedAt,
		PublishedAt:    category.PublishedAt,
		CreatorID:      category.CreatorID,
		LikesCount:     category.LikesCount,
		IsCreator:      1,
	})
}

func (h *Handler) PublishConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	description := ctx.PostForm("description")

	historicalDate, err := time.Parse(
		"2006-01-02",
		ctx.PostForm("historical_date"),
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid historical_date",
		})
		return
	}

	basketWeight, err := strconv.ParseFloat(
		ctx.PostForm("basket_weight"),
		64,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid basket_weight",
		})
		return
	}

	currentUser := user.GetCurrentUser()

	err = h.Repository.PublishConsumerCategory(
		id,
		currentUser.UserID,
		description,
		historicalDate,
		basketWeight,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "consumer category published",
	})
}

func (h *Handler) GetConsumerCategoryFeedByIDAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	next := ctx.Query("next") == "true"
	currentUser := user.GetCurrentUser()

	category, err := h.Repository.GetConsumerCategoryByID(
		id,
		next,
		currentUser.UserID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	isCreator := 0
	if category.CreatorID == currentUser.UserID {
		isCreator = 1
	}

	response := serializer.ConsumerCategoryResponse{
		ID:             category.CategoryID,
		Title:          category.Title,
		Description:    category.Description,
		Status:         category.Status,
		ImageURL:       category.ImageURL,
		VideoURL:       category.VideoURL,
		BasketWeight:   category.BasketWeight,
		HistoricalDate: category.HistoricalDate,
		CreatedAt:      category.CreatedAt,
		PublishedAt:    category.PublishedAt,
		CreatorID:      category.CreatorID,
		LikesCount:     category.LikesCount,
		IsCreator:      isCreator,
		IsLiked:        category.IsLiked,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) LikeConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	likeValue, err := strconv.Atoi(ctx.PostForm("like"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid like"})
		return
	}

	if likeValue != 0 && likeValue != 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "like must be 0 or 1"})
		return
	}

	currentUser := user.GetCurrentUser()

	err = h.Repository.SetLike(
		uint(id),
		currentUser.UserID,
		likeValue,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "like updated",
	})
}

func (h *Handler) DeleteConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	currentUser := user.GetCurrentUser()

	err = h.Repository.DeleteConsumerCategory(
		id,
		currentUser.UserID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "consumer category deleted",
	})
}

func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var request serializer.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if request.FullName == "" || request.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "full_name and email are required",
		})
		return
	}

	newUser, err := h.Repository.CreateUser(
		request.FullName,
		request.Email,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, serializer.UserResponse{
		ID:       newUser.UserID,
		FullName: newUser.FullName,
		Email:    newUser.Email,
	})
}

func (h *Handler) AuthUserAPI(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "auth is not implemented yet",
	})
}

func (h *Handler) DeauthUserAPI(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "deauth is not implemented yet",
	})
}
