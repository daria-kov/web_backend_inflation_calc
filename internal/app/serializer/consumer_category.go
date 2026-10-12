package serializer

import "time"

type ConsumerCategoryResponse struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	ImageURL       string    `json:"image_url"`
	VideoURL       string    `json:"video_url"`
	BasketWeight   float64   `json:"basket_weight"`
	HistoricalDate time.Time `json:"historical_date"`
	CreatedAt      time.Time `json:"created_at"`
	CreatorID      uint      `json:"creator_id"`
	LikesCount     int       `json:"likes_count"`
	IsCreator      int       `json:"is_creator"`
	IsLiked        int       `json:"is_liked"`
}
