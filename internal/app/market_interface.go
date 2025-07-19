package app

import(
	"marketplace/internal/config"
	"github.com/google/uuid"
)

type MarketServicer interface {
	NewAd(ad Ad, config config.Config, userid uuid.UUID) (Ad, error)
	AdsList(params AdsListParams, id uuid.UUID) ([]Ad, error)
}

type MarketRepository interface {
	SaveAd(ad Ad) (Ad, error)
	GetAdsList(params AdsListParams, user_id string) ([]Ad, error)
}