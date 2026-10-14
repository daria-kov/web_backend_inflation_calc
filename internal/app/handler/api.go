package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// GetConsumerCategoriesAPI godoc
// @Summary Получить список категорий
// @Description Возвращает список опубликованных категорий с количеством лайков и информацией о лайке текущего пользователя.
// @Tags Consumer Categories
// @Produce json
// @Param date query string false "Дата исторического события" Format(YYYY-MM-DD)
// @Success 200 {array} serializer.ConsumerCategoryResponse
// @Failure 500 "Внутренняя ошибка сервера"
// @Security BearerAuth
// @Router /api/consumer_categories [get]
func (h *Handler) GetConsumerCategoriesAPI(ctx *gin.Context) {
	filterDate := ctx.Query("date")
	currentUser := getCurrentUser(ctx)

	categories, err := h.Repository.GetConsumerCategories(
		filterDate,
		currentUser.UserID,
	)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
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

// CreateConsumerCategoryAPI godoc
// @Summary Создать черновик категории
// @Description Создаёт черновик категории и загружает изображение и видео в хранилище.
// @Tags Consumer Categories
// @Accept mpfd
// @Produce json
// @Param title formData string true "Название категории"
// @Param image formData file true "Изображение категории"
// @Param video formData file true "Видео категории"
// @Success 201 "Черновик создан"
// @Failure 400 "Некорректные данные или у пользователя уже есть черновик"
// @Failure 500 "Ошибка обработки файла или сервера"
// @Security BearerAuth
// @Router /api/consumer_categories [post]
func (h *Handler) CreateConsumerCategoryAPI(ctx *gin.Context) {
	title := ctx.PostForm("title")

	if strings.TrimSpace(title) == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	video, err := ctx.FormFile("video")
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(ctx)

	imageExt := filepath.Ext(image.Filename)
	videoExt := filepath.Ext(video.Filename)

	imageFilename, err := generateFilename(imageExt)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	videoFilename, err := generateFilename(videoExt)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	imageFile, err := image.Open()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
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
		ctx.Status(http.StatusInternalServerError)
		return
	}

	videoFile, err := video.Open()
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
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
		ctx.Status(http.StatusInternalServerError)
		return
	}

	err = h.Repository.CreateConsumerCategoryDraft(
		title,
		imageFilename,
		videoFilename,
		currentUser.UserID,
	)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusCreated)
}

// GetConsumerCategoryDraftAPI godoc
// @Summary Получить свой черновик
// @Description Возвращает последний черновик текущего пользователя.
// @Tags Consumer Categories
// @Produce json
// @Success 200 {object} serializer.ConsumerCategoryResponse
// @Failure 404 "Черновик не найден"
// @Failure 500 "Ошибка сервера"
// @Security BearerAuth
// @Router /api/consumer_categories/draft [get]
func (h *Handler) GetConsumerCategoryDraftAPI(ctx *gin.Context) {
	currentUser := getCurrentUser(ctx)

	category, err := h.Repository.GetConsumerCategoryDraft(
		currentUser.UserID,
	)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if category.CategoryID == 0 {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, serializer.ConsumerCategoryResponse{
		ID:             category.CategoryID,
		Title:          category.Title,
		Description:    category.Description,
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

// PublishConsumerCategoryAPI godoc
// @Summary Опубликовать категорию
// @Description Публикует черновик текущего пользователя, сохраняя описание, историческую дату и вес корзины.
// @Tags Consumer Categories
// @Accept mpfd
// @Produce json
// @Param id path int true "ID категории"
// @Param description formData string false "Описание категории"
// @Param historical_date formData string true "Историческая дата в формате YYYY-MM-DD"
// @Param basket_weight formData number true "Вес корзины"
// @Success 200 "Категория опубликована"
// @Failure 400 "Некорректные данные или черновик не найден"
// @Security BearerAuth
// @Router /api/consumer_categories/{id}/publish [put]
func (h *Handler) PublishConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	description := ctx.PostForm("description")

	historicalDate, err := time.Parse(
		"2006-01-02",
		ctx.PostForm("historical_date"),
	)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	basketWeight, err := strconv.ParseFloat(
		ctx.PostForm("basket_weight"),
		64,
	)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(ctx)

	err = h.Repository.PublishConsumerCategory(
		id,
		currentUser.UserID,
		description,
		historicalDate,
		basketWeight,
	)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.Status(http.StatusOK)
}

// GetConsumerCategoryFeedAPI godoc
// @Summary Получить первую категорию ленты
// @Description Возвращает первую опубликованную категорию. С неё начинается просмотр ленты.
// @Tags Consumer Categories
// @Produce json
// @Success 200 {object} serializer.ConsumerCategoryResponse
// @Failure 404 "Опубликованных категорий нет"
// @Security BearerAuth
// @Router /api/consumer_categories/feed [get]
func (h *Handler) GetConsumerCategoryFeedAPI(ctx *gin.Context) {
	currentUser := getCurrentUser(ctx)

	category, err := h.Repository.GetFirstConsumerCategory(currentUser.UserID)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	isCreator := 0
	if category.CreatorID == currentUser.UserID {
		isCreator = 1
	}

	ctx.JSON(http.StatusOK, serializer.ConsumerCategoryResponse{
		ID:             category.CategoryID,
		Title:          category.Title,
		Description:    category.Description,
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

// GetConsumerCategoryFeedByIDAPI godoc
// @Summary Получить категорию по ID
// @Description Возвращает опубликованную категорию. При next=true возвращает следующую опубликованную категорию, а при достижении конца списка переходит к первой.
// @Tags Consumer Categories
// @Produce json
// @Param id path int true "ID категории"
// @Param next query bool false "Перейти к следующей категории"
// @Success 200 {object} serializer.ConsumerCategoryResponse
// @Failure 400 "Некорректный ID"
// @Failure 404 "Категория не найдена"
// @Security BearerAuth
// @Router /api/consumer_categories/feed/{id} [get]
func (h *Handler) GetConsumerCategoryFeedByIDAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	next := ctx.Query("next") == "true"
	currentUser := getCurrentUser(ctx)

	category, err := h.Repository.GetConsumerCategoryByID(
		id,
		next,
		currentUser.UserID,
	)
	if err != nil {
		ctx.Status(http.StatusNotFound)
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

// LikeConsumerCategoryAPI godoc
// @Summary Поставить или снять лайк
// @Description Устанавливает или снимает лайк текущего пользователя для опубликованной категории.
// @Tags Consumer Categories
// @Accept mpfd
// @Produce json
// @Param id path int true "ID категории"
// @Param like formData int true "Действие: 1 — поставить лайк, 0 — снять"
// @Success 200 "Лайк обновлён"
// @Failure 400 "Некорректный ID или значение like"
// @Failure 500 "Ошибка обновления лайка"
// @Security BearerAuth
// @Router /api/consumer_categories/{id}/like [post]
func (h *Handler) LikeConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	likeValue, err := strconv.Atoi(ctx.PostForm("like"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if likeValue != 0 && likeValue != 1 {
		ctx.Status(http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(ctx)

	err = h.Repository.SetLike(
		uint(id),
		currentUser.UserID,
		likeValue,
	)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteConsumerCategoryAPI godoc
// @Summary Удалить категорию
// @Description Меняет статус категории на deleted. Удалять можно только свою категорию.
// @Tags Consumer Categories
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 "Категория удалена"
// @Failure 400 "Некорректный ID"
// @Failure 404 "Категория не найдена или принадлежит другому пользователю"
// @Security BearerAuth
// @Router /api/consumer_categories/{id} [delete]
func (h *Handler) DeleteConsumerCategoryAPI(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	currentUser := getCurrentUser(ctx)

	err = h.Repository.DeleteConsumerCategory(
		id,
		currentUser.UserID,
	)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.Status(http.StatusOK)
}

// RegisterUserAPI godoc
// @Summary Зарегистрировать пользователя
// @Description Создаёт пользователя с логином, именем и паролем. Пароль сохраняется в виде bcrypt-хеша.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body serializer.RegisterUserRequest true "Данные нового пользователя"
// @Success 201 {object} serializer.UserResponse
// @Failure 400 "Некорректные данные или ошибка создания пользователя"
// @Router /api/users/register [post]
func (h *Handler) RegisterUserAPI(ctx *gin.Context) {
	var request serializer.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	if request.Login == "" ||
		request.FullName == "" ||
		request.Password == "" {
		ctx.Status(http.StatusBadRequest)
		return
	}

	newUser, err := h.Repository.CreateUser(
		request.Login,
		request.FullName,
		request.Password,
	)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusCreated, serializer.UserResponse{
		ID:       newUser.UserID,
		Login:    newUser.Login,
		FullName: newUser.FullName,
	})
}

// AuthUserAPI godoc
// @Summary Авторизовать пользователя
// @Description Проверяет логин и пароль и возвращает JWT-токен со сроком действия 24 часа.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body serializer.AuthUserRequest true "Данные для входа"
// @Success 200 {object} map[string]interface{} "Токен доступа и данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неверный логин или пароль"
// @Failure 500 "Ошибка сервера"
// @Router /api/users/auth [post]
func (h *Handler) AuthUserAPI(ctx *gin.Context) {
	var request serializer.AuthUserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	foundUser, err := h.Repository.GetUserByLogin(request.Login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		ctx.Status(http.StatusInternalServerError)
		return
	}

	// Сначала проверяем bcrypt-хеш.
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.Password),
		[]byte(request.Password),
	)

	if err != nil {
		// Совместимость со старым аккаунтом, где пароль был открытым текстом.
		if foundUser.Password != request.Password {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		// При успешном входе заменяем старый пароль на bcrypt-хеш.
		passwordHash, hashErr := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)
		if hashErr != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if err := h.Repository.UpdateUserPassword(
			foundUser.UserID,
			string(passwordHash),
		); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": foundUser.UserID,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   86400,
		"user": serializer.UserResponse{
			ID:       foundUser.UserID,
			Login:    foundUser.Login,
			FullName: foundUser.FullName,
		},
	})
}

// DeauthUserAPI godoc
// @Summary Выйти из системы
// @Description Отзывает текущий JWT, добавляя его хеш в чёрный список Redis до истечения срока действия токена.
// @Tags Users
// @Produce json
// @Success 200 "Выход выполнен успешно"
// @Failure 401 "Токен отсутствует, недействителен или истёк"
// @Failure 500 "JWT_SECRET не настроен"
// @Failure 503 "Не удалось отозвать токен"
// @Security BearerAuth
// @Router /api/users/deauth [post]
func (h *Handler) DeauthUserAPI(ctx *gin.Context) {
	// Получаем JWT из заголовка Authorization.
	authHeader := ctx.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		strings.TrimSpace(parts[1]) == "" {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	rawToken := strings.TrimSpace(parts[1])

	// Проверяем токен и получаем срок его действия.
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil || token == nil || !token.Valid {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	expiration, err := claims.GetExpirationTime()
	if err != nil || expiration == nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	ttl := time.Until(expiration.Time)
	if ttl <= 0 {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	// Вычисляем хеш токена — такой же, как в middleware.
	tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))
	blacklistKey := "blacklist:" + tokenHash

	// Добавляем токен в чёрный список до окончания его срока действия.
	err = h.Redis.Set(
		ctx.Request.Context(),
		blacklistKey,
		"1",
		ttl,
	).Err()

	if err != nil {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}

	ctx.Status(http.StatusOK)
}

func getCurrentUser(ctx *gin.Context) *ds.Users {
	value, exists := ctx.Get("userID")
	if !exists {
		// Для публичных запросов без авторизации
		return &ds.Users{}
	}

	userID, ok := value.(uint)
	if !ok {
		return &ds.Users{}
	}

	return &ds.Users{UserID: userID}
}
