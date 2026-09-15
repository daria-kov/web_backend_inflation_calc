package repository

import (
	"fmt"
	"strings"
	"time"
)

func ParseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)

	// Пробуем формат DD.MM.YYYY
	t, err := time.Parse("02.01.2006", dateStr)
	if err == nil {
		return t, nil
	}

	// Пробуем формат YYYY-MM-DD
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("неверный формат даты: %s", dateStr)
}

func FormatDateToInput(dateStr string) string {
	t, err := ParseDate(dateStr)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type ConsumerCategory struct {
	ID           int
	Title        string  // Название категории (например, "Продукты")
	Status       string  // "published", "draft", "deleted"
	ImageKey     string  // Ключ изображения в Minio (латиница)
	VideoKey     string  // Ключ видео в Minio (латиница)
	Description  string  // Текстовое описание (для кнопки "ещё")
	BasketWeight float64 // Вес в корзине (%) - числовое поле для фильтрации
	Date         string
	SumDate      int   // Сумма тогда (руб)
	SumNow       int   // Сумма сейчас (руб)
	Likes        []int // Массив ID пользователей, поставивших лайк
}

func (r *Repository) GetConsumerCategories(filterFrom, filterTo string) ([]ConsumerCategory, error) {
	allCategories := r.getAllData()
	var result []ConsumerCategory

	var fromDate, toDate time.Time
	var hasFromFilter, hasToFilter bool

	if filterFrom != "" {
		t, err := time.Parse("2006-01-02", filterFrom)
		if err == nil {
			fromDate = t
			hasFromFilter = true
		}
	}

	if filterTo != "" {
		t, err := time.Parse("2006-01-02", filterTo)
		if err == nil {
			toDate = t
			hasToFilter = true
		}
	}

	for _, c := range allCategories {
		if c.Status == "deleted" || c.Status == "draft" {
			continue
		}

		// Фильтрация по датам
		if hasFromFilter || hasToFilter {
			catDate, err := ParseDate(c.Date)
			if err != nil {
				continue
			}

			if hasFromFilter && catDate.Before(fromDate) {
				continue
			}
			if hasToFilter && catDate.After(toDate) {
				continue
			}
		}

		result = append(result, c)
	}

	return result, nil
}

func (r *Repository) GetConsumerCategoryByID(id int, next bool) (ConsumerCategory, error) {
	allCategories := r.getAllData()

	var visible []ConsumerCategory
	for _, c := range allCategories {
		if c.Status == "published" {
			visible = append(visible, c)
		}
	}

	if len(visible) == 0 {
		return ConsumerCategory{}, fmt.Errorf("нет доступных категорий")
	}

	if !next {
		if id == 0 {
			return visible[0], nil
		}
		for _, c := range visible {
			if c.ID == id {
				return c, nil
			}
		}
		return ConsumerCategory{}, fmt.Errorf("категория не найдена")
	}

	foundIndex := -1
	for i, c := range visible {
		if c.ID == id {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return visible[0], nil
	}
	nextIndex := (foundIndex + 1) % len(visible)
	return visible[nextIndex], nil
}

func (r *Repository) GetConsumerCategoryDraft() (ConsumerCategory, error) {
	allCategories := r.getAllData()
	for _, c := range allCategories {
		if c.Status == "draft" {
			return c, nil
		}
	}
	return ConsumerCategory{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) getAllData() []ConsumerCategory {
	return []ConsumerCategory{
		{
			ID:           1,
			Title:        "Продукты",
			Status:       "published",
			ImageKey:     "goods.jpg",
			VideoKey:     "goods.mp4",
			Description:  "Самая тяжелая статья расходов. Включает мясо, молоко, овощи, фрукты и хлебобулочные изделия. Цены на эту категорию растут быстрее всего из-за сезонности и логистики.",
			BasketWeight: 35.0,
			Date:         "03.02.2026",
			SumDate:      48590,
			SumNow:       51240,
			Likes:        []int{101, 102, 105, 108, 112},
		},
		{
			ID:           2,
			Title:        "Одежда и обувь",
			Status:       "published",
			ImageKey:     "clothes.jpg",
			VideoKey:     "clothes.mp4",
			Description:  "Сезонные расходы на гардероб. Сюда входят верхняя одежда, обувь и повседневные вещи. Вес категории сильно зависит от времени года.",
			BasketWeight: 12.5,
			Date:         "09.11.2026",
			SumDate:      34990,
			SumNow:       37800,
			Likes:        []int{101, 103},
		},
		{
			ID:           3,
			Title:        "Развлечения",
			Status:       "draft",
			ImageKey:     "amuzement.jpg",
			VideoKey:     "amuzement.mp4",
			Description:  "Отдыхать тоже важно, поэтому не забываем о развлечениях: музеи, парки, театры, кафе и кино.",
			BasketWeight: 18.2,
			Date:         "15.04.2026",
			SumDate:      5900,
			SumNow:       4500,
			Likes:        []int{101, 102, 103, 104, 105, 106, 107},
		},
	}
}
