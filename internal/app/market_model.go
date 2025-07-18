package app

import (
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
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	Owner      	bool      `json:"owner,omitempty"` 
}

type MarketService struct {
	Marketrepo MarketRepository
	Userrepo   UserRepository
}

type AdsListParams struct {
	Page     int     `query:"page"`      
	Limit    int     `query:"limit"`     
	SortBy   string  `query:"sort_by"`   
	Order    string  `query:"order"` // "asc" or "desc"
	MinPrice int     `query:"min_price"` 
	MaxPrice int     `query:"max_price"` 
}