package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetConsumerCategories(
	filterDate string,
) ([]ds.ConsumerCategory, error) {
	return r.GetConsumerCategoriesForUser(filterDate, 0)
}

func (r *Repository) GetConsumerCategoriesForUser(
	filterDate string,
	userID uint,
) ([]ds.ConsumerCategory, error) {
	var categories []ds.ConsumerCategory

	query := r.db.
		Model(&ds.ConsumerCategory{}).
		Select(`
			consumer_categories.*,
			COUNT(likes.like_id) AS likes_count,
			CASE
				WHEN EXISTS (
					SELECT 1
					FROM likes AS user_likes
					WHERE user_likes.category_id =
						consumer_categories.category_id
					AND user_likes.user_id = ?
				)
				THEN 1
				ELSE 0
			END AS is_liked
		`, userID).
		Joins(
			"LEFT JOIN likes ON likes.category_id = consumer_categories.category_id",
		).
		Where("consumer_categories.status = ?", "published").
		Group("consumer_categories.category_id")

	if filterDate != "" {
		date, err := time.Parse("2006-01-02", filterDate)
		if err != nil {
			return nil, err
		}

		query = query.Where(
			"consumer_categories.historical_date = ?",
			date,
		)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) GetConsumerCategoryByID(
	id int,
	next bool,
	userID uint,
) (ds.ConsumerCategory, error) {
	var category ds.ConsumerCategory

	if !next {
		err := r.db.
			Model(&ds.ConsumerCategory{}).
			Select(`
				consumer_categories.*,
				COUNT(likes.like_id) AS likes_count,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM likes AS user_likes
						WHERE user_likes.category_id =
							consumer_categories.category_id
						AND user_likes.user_id = ?
					)
					THEN 1
					ELSE 0
				END AS is_liked
			`, userID).
			Joins(
				"LEFT JOIN likes ON likes.category_id = consumer_categories.category_id",
			).
			Where(
				"consumer_categories.category_id = ? AND consumer_categories.status = ?",
				id,
				"published",
			).
			Group("consumer_categories.category_id").
			First(&category).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ds.ConsumerCategory{}, fmt.Errorf("категория не найдена")
			}

			return ds.ConsumerCategory{}, err
		}

		return category, nil
	}

	var current ds.ConsumerCategory

	err := r.db.
		Where(
			"category_id = ? AND status = ?",
			id,
			"published",
		).
		First(&current).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ConsumerCategory{}, fmt.Errorf("категория не найдена")
		}

		return ds.ConsumerCategory{}, err
	}

	err = r.db.
		Model(&ds.ConsumerCategory{}).
		Select(`
			consumer_categories.*,
			COUNT(likes.like_id) AS likes_count,
			CASE
				WHEN EXISTS (
					SELECT 1
					FROM likes AS user_likes
					WHERE user_likes.category_id =
						consumer_categories.category_id
					AND user_likes.user_id = ?
				)
				THEN 1
				ELSE 0
			END AS is_liked
		`, userID).
		Joins(
			"LEFT JOIN likes ON likes.category_id = consumer_categories.category_id",
		).
		Where(
			"consumer_categories.category_id > ? AND consumer_categories.status = ?",
			current.CategoryID,
			"published",
		).
		Group("consumer_categories.category_id").
		Order("consumer_categories.category_id ASC").
		First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.db.
			Model(&ds.ConsumerCategory{}).
			Select(`
				consumer_categories.*,
				COUNT(likes.like_id) AS likes_count,
				CASE
					WHEN EXISTS (
						SELECT 1
						FROM likes AS user_likes
						WHERE user_likes.category_id =
							consumer_categories.category_id
						AND user_likes.user_id = ?
					)
					THEN 1
					ELSE 0
				END AS is_liked
			`, userID).
			Joins(
				"LEFT JOIN likes ON likes.category_id = consumer_categories.category_id",
			).
			Where("consumer_categories.status = ?", "published").
			Group("consumer_categories.category_id").
			Order("consumer_categories.category_id ASC").
			First(&category).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ConsumerCategory{}, fmt.Errorf("категория не найдена")
		}

		return ds.ConsumerCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetConsumerCategoryDraft(creatorID uint) (ds.ConsumerCategory, error) {
	var category ds.ConsumerCategory

	err := r.db.
		Where(
			"creator_id = ? AND status = ?",
			creatorID,
			"draft",
		).
		Order("category_id DESC").
		First(&category).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ConsumerCategory{}, nil
		}

		return ds.ConsumerCategory{}, err
	}

	return category, nil
}

func (r *Repository) CreateConsumerCategoryDraft(title string, imageFilename string, videoFilename string, creatorID uint) error {
	var count int64

	err := r.db.
		Model(&ds.ConsumerCategory{}).
		Where("creator_id = ? AND status = ?", creatorID, "draft").
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("у пользователя уже есть draft")
	}

	category := ds.ConsumerCategory{
		Title:          title,
		Description:    "",
		Status:         "draft",
		ImageURL:       imageFilename,
		VideoURL:       videoFilename,
		BasketWeight:   0,
		HistoricalDate: time.Now(),
		CreatedAt:      time.Now(),
		CreatorID:      creatorID,
	}

	return r.db.Create(&category).Error
}

func (r *Repository) PublishConsumerCategory(
	id int,
	creatorID uint,
	description string,
	historicalDate time.Time,
	basketWeight float64,
) error {
	result := r.db.
		Model(&ds.ConsumerCategory{}).
		Where(
			"category_id = ? AND creator_id = ? AND status = ?",
			id,
			creatorID,
			"draft",
		).
		Updates(map[string]interface{}{
			"description":     description,
			"historical_date": historicalDate,
			"basket_weight":   basketWeight,
			"status":          "published",
			"published_at":    time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("draft не найден или принадлежит другому пользователю")
	}

	return nil
}

func (r *Repository) DeleteConsumerCategory(id int, creatorID uint) error {
	result := r.db.
		Model(&ds.ConsumerCategory{}).
		Where(
			"category_id = ? AND creator_id = ? AND status != ?",
			id,
			creatorID,
			"deleted",
		).
		Update("status", "deleted")

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("категория не найдена или принадлежит другому пользователю")
	}

	return nil
}

func (r *Repository) SetLike(categoryID uint, userID uint, likeValue int) error {
	var like ds.Like

	err := r.db.
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&like).Error

	if likeValue == 1 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			like = ds.Like{
				CategoryID: categoryID,
				UserID:     userID,
			}

			return r.db.Create(&like).Error
		}

		return err
	}

	if likeValue == 0 {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		return r.db.Delete(&like).Error
	}

	return fmt.Errorf("like должен быть 0 или 1")
}
