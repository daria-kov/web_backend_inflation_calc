package repository

import (
	"fmt"
)

// Структура репозитория (пока пустая, данные хардкодим в методах)
type Repository struct{}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

// Модель услуги (Категория потребительской корзины)
type Category struct {
	ID           int
	Title        string  // Название категории (например, "Продукты")
	Status       string  // "published", "draft", "deleted"
	ImageKey     string  // Ключ изображения в Minio (латиница)
	VideoKey     string  // Ключ видео в Minio (латиница)
	Description  string  // Текстовое описание (для кнопки "ещё")
	BasketWeight float64 // Вес в корзине (%) - числовое поле для фильтрации
	PreviousSum  int     // Сумма год назад (руб)
	TodaySum     int     // Сумма сейчас (руб)
	Likes        []int   // Массив ID пользователей, поставивших лайк
}

func (r *Repository) GetCategories(filterWeight string) ([]Category, error) {
	allCategories := r.getAllData()
	var result []Category

	for _, c := range allCategories {
		if c.Status == "deleted" {
			continue
		}

		if filterWeight != "" {
			var filterVal float64
			_, err := fmt.Sscanf(filterWeight, "%f", &filterVal)
			if err == nil {
				if c.BasketWeight < filterVal {
					continue
				}
			}
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *Repository) GetCategoryByID(id int, next bool) (Category, error) {
	allCategories := r.getAllData()

	// Собираем все, КРОМЕ удаленных (и published, и draft)
	var visible []Category
	for _, c := range allCategories {
		if c.Status != "deleted" {
			visible = append(visible, c)
		}
	}

	if len(visible) == 0 {
		return Category{}, fmt.Errorf("нет доступных категорий")
	}

	// Логика поиска (та же, что была для published, но теперь для visible)
	if !next {
		if id == 0 {
			return visible[0], nil
		}
		for _, c := range visible {
			if c.ID == id {
				return c, nil
			}
		}
		return Category{}, fmt.Errorf("категория не найдена")
	}

	// Логика next
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

func (r *Repository) GetDraft() (Category, error) {
	allCategories := r.getAllData()
	for _, c := range allCategories {
		if c.Status == "draft" {
			return c, nil
		}
	}
	return Category{}, fmt.Errorf("черновик не найден")
}

func (r *Repository) getAllData() []Category {
	return []Category{
		{
			ID:           1,
			Title:        "Продукты",
			Status:       "published",
			ImageKey:     "goods.jpg",
			VideoKey:     "goods.mp4",
			Description:  "Самая тяжелая статья расходов. Включает мясо, молоко, овощи, фрукты и хлебобулочные изделия. Цены на эту категорию растут быстрее всего из-за сезонности и логистики.",
			BasketWeight: 35.0,
			PreviousSum:  48590,
			TodaySum:     51240,
			Likes:        []int{101, 102, 105, 108, 112},
		},
		{
			ID:           2,
			Title:        "Одежда и обувь",
			Status:       "draft",
			ImageKey:     "clothes.jpg",
			VideoKey:     "clothes.mp4",
			Description:  "Сезонные расходы на гардероб. Сюда входят верхняя одежда, обувь и повседневные вещи. Вес категории сильно зависит от времени года.",
			BasketWeight: 12.5,
			PreviousSum:  34990,
			TodaySum:     37800,
			Likes:        []int{101, 103},
		},
		{
			ID:           3,
			Title:        "Развлечения",
			Status:       "deleted",
			ImageKey:     "amuzement.jpg",
			VideoKey:     "amuzement.mp4",
			Description:  "Отдыхать тоже важно, поэтому не забываем о развлечениях: музеи, парки, театры, кафе и кино.",
			BasketWeight: 18.2,
			PreviousSum:  5900,
			TodaySum:     4500,
			Likes:        []int{101, 102, 103, 104, 105, 106, 107},
		},
	}
}
