package app

import (
	"marketplace/internal/config"
	"time"
	"github.com/google/uuid"
)

type Ad struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}