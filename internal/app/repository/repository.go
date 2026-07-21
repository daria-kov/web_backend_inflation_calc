package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Planet struct {
	ID          int
	Title       string
	Mass        float64 // Масса в кг (например, 3.30e23)
	Radius      int     // Радиус в км
	Picture     string  // Ключ изображения в Minio
	Description string  // Краткое описание
}

func (r *Repository) GetPlanets() ([]Planet, error) {
	planets := []Planet{
		{
			ID:          1,
			Title:       "Меркурий",
			Mass:        3.30e23,
			Radius:      2440,
			Picture:     "mercury.jpg",
			Description: "Ближайшая к Солнцу планета, самая маленькая в Солнечной системе. Температура колеблется от -180°C до +430°C. Атмосфера почти отсутствует, поверхность покрыта кратерами.",
		},
		{
			ID:          2,
			Title:       "Венера",
			Mass:        4.87e24,
			Radius:      6052,
			Picture:     "venus.jpg",
			Description: "Самая горячая планета (температура ~460°C) из-за парникового эффекта. Вращается в обратную сторону (ось наклонена на 177°). Плотная атмосфера из углекислого газа с облаками серной кислоты.",
		},
		{
			ID:          3,
			Title:       "Земля",
			Mass:        5.97e24,
			Radius:      6371,
			Picture:     "earth.jpg",
			Description: "Третья планета от Солнца, единственная известная планета с жизнью. Имеет жидкую воду на поверхности, азотно-кислородную атмосферу и магнитное поле, защищающее от солнечного ветра.",
		},
		{
			ID:          4,
			Title:       "Марс",
			Mass:        6.42e23,
			Radius:      3390,
			Picture:     "mars.jpg",
			Description: "Красная планета из-за оксида железа (ржавчины) на поверхности. Имеет самую высокую гору в Солнечной системе — Олимп (21.9 км). Атмосфера разрежена, состоит в основном из углекислого газа.",
		},
		{
			ID:          5,
			Title:       "Юпитер",
			Mass:        1.90e27,
			Radius:      69911,
			Picture:     "jupyter.jpg",
			Description: "Самая большая планета, масса превышает массу всех остальных планет вместе взятых. Состоит из водорода и гелия. Имеет Большое Красное Пятно — гигантский шторм, который длится сотни лет.",
		},
		{
			ID:          6,
			Title:       "Сатурн",
			Mass:        5.68e26,
			Radius:      58232,
			Picture:     "saturn.jpg",
			Description: "Знаменита своей системой колец, состоящих из льда и камней. Состоит в основном из водорода и гелия. Плотность планеты меньше плотности воды. Имеет 146 известных спутников.",
		},
		{
			ID:          7,
			Title:       "Уран",
			Mass:        8.68e25,
			Radius:      25362,
			Picture:     "uranus.jpg",
			Description: "Ледяной гигант с наклоном оси почти 98° (вращается 'на боку'). Состоит из воды, метана и аммиака. Имеет слабую систему колец и 27 спутников.",
		},
		{
			ID:          8,
			Title:       "Нептун",
			Mass:        1.02e26,
			Radius:      24622,
			Picture:     "neptunejpg.jpg",
			Description: "Самая далекая планета, самый сильный ветер в Солнечной системе (до 2100 км/ч). Имеет темно-синий цвет из-за метана в атмосфере. Была открыта математически (по возмущениям орбиты Урана).",
		},
	}

	if len(planets) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return planets, nil
}

func (r *Repository) GetPlanet(id int) (Planet, error) {
	planets, err := r.GetPlanets()
	if err != nil {
		return Planet{}, err
	}

	for _, planet := range planets {
		if planet.ID == id {
			return planet, nil
		}
	}
	return Planet{}, fmt.Errorf("Планета не найдена")
}

func (r *Repository) GetPlanetByTitle(title string) ([]Planet, error) {
	planets, err := r.GetPlanets()
	if err != nil {
		return []Planet{}, err
	}

	var result []Planet
	for _, planet := range planets {
		if strings.Contains(strings.ToLower(planet.Title), strings.ToLower(title)) {
			result = append(result, planet)
		}
	}
	return result, nil
}

func (r *Repository) GetCart() ([]Planet, error) {
	planets, err := r.GetPlanets()
	if err != nil {
		return []Planet{}, err
	}

	var result []Planet
	for _, planet := range planets {
		if planet.ID == 3 || planet.ID == 4 {
			result = append(result, planet)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Корзина пуста")
	}

	return result, nil
}
