package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetConsumerCategories(filterDate string) ([]ds.ConsumerCategory, error) {
	var categories []ds.ConsumerCategory

	query := r.db.
		Model(&ds.ConsumerCategory{}).
		Select("consumer_categories.*, COUNT(likes.like_id) AS likes_count").
		Joins("LEFT JOIN likes ON likes.category_id = consumer_categories.category_id").
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

func (r *Repository) GetConsumerCategoryByID(id int, next bool) (ds.ConsumerCategory, error) {
	var category ds.ConsumerCategory

	if !next {
		err := r.db.
			Model(&ds.ConsumerCategory{}).
			Select("consumer_categories.*, COUNT(likes.like_id) AS likes_count").
			Joins("LEFT JOIN likes ON likes.category_id = consumer_categories.category_id").
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
		Select("consumer_categories.*, COUNT(likes.like_id) AS likes_count").
		Joins("LEFT JOIN likes ON likes.category_id = consumer_categories.category_id").
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
			Select("consumer_categories.*, COUNT(likes.like_id) AS likes_count").
			Joins("LEFT JOIN likes ON likes.category_id = consumer_categories.category_id").
			Where("consumer_categories.status = ?", "published").
			Group("consumer_categories.category_id").
			Order("consumer_categories.category_id ASC").
			First(&category).Error
	}

	if err != nil {
		return ds.ConsumerCategory{}, err
	}

	return category, nil
}

func (r *Repository) GetConsumerCategoryDraft() (ds.ConsumerCategory, error) {
	var category ds.ConsumerCategory

	err := r.db.
		Where("status = ?", "draft").
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

func (r *Repository) CreateConsumerCategoryDraft(title string) error {
	category := ds.ConsumerCategory{
		Title:          title,
		Description:    "",
		Status:         "draft",
		ImageURL:       "",
		VideoURL:       "",
		BasketWeight:   0,
		HistoricalDate: time.Now(),
		CreatedAt:      time.Now(),
		CreatorID:      1,
	}

	return r.db.Create(&category).Error
}

func (r *Repository) PublishConsumerCategory(
	id int,
	description string,
	historicalDate time.Time,
	basketWeight float64,
) error {

	return r.db.
		Model(&ds.ConsumerCategory{}).
		Where("category_id = ? AND status = ?", id, "draft").
		Updates(map[string]interface{}{
			"description":     description,
			"historical_date": historicalDate,
			"basket_weight":   basketWeight,
			"status":          "published",
		}).Error
}

func (r *Repository) DeleteConsumerCategory(id int) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE consumer_categories SET status = 'deleted' WHERE category_id = $1",
		id,
	)

	return err
}
