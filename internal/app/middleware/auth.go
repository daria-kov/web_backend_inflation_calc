package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Получаем токен из заголовка Authorization.
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimSpace(parts[1])

		// 2. Получаем секрет для проверки подписи JWT.
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 3. Проверяем подпись и срок действия токена.
		token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 4. Проверяем, не отозван ли токен в Redis.
		tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))
		blacklistKey := "blacklist:" + tokenHash

		exists, err := redisClient.Exists(
			ctx.Request.Context(),
			blacklistKey,
		).Result()

		if err != nil {
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
			return
		}

		if exists > 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 5. Получаем user_id из проверенного токена.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userIDValue, ok := claims["user_id"].(float64)
		if !ok || userIDValue <= 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 6. Сохраняем ID пользователя в контексте Gin.
		ctx.Set("userID", uint(userIDValue))

		ctx.Next()
	}
}

// OptionalAuthMiddleware — для публичных маршрутов.
// Если передан валидный токен, сохраняет userID в контексте.
// Если токена нет или он невалидный, запрос продолжается как гостевой.
func OptionalAuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. Получаем токен из заголовка Authorization.
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			strings.TrimSpace(parts[1]) == "" {
			ctx.Next()
			return
		}

		rawToken := strings.TrimSpace(parts[1])

		// 2. Проверяем подпись и срок действия токена.
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			ctx.Next()
			return
		}

		token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || token == nil || !token.Valid {
			ctx.Next()
			return
		}

		// 3. Проверяем, не отозван ли токен в Redis.
		tokenHash := fmt.Sprintf("%x", sha256.Sum256([]byte(rawToken)))
		exists, err := redisClient.Exists(
			ctx.Request.Context(),
			"blacklist:"+tokenHash,
		).Result()

		if err != nil || exists > 0 {
			ctx.Next()
			return
		}

		// 4. Получаем user_id и сохраняем в контексте.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.Next()
			return
		}

		userIDValue, ok := claims["user_id"].(float64)
		if ok && userIDValue > 0 {
			ctx.Set("userID", uint(userIDValue))
		}

		ctx.Next()
	}
}
