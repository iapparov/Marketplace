package app

import "github.com/google/uuid"

type MarketRepository interface {
	SaveAd(ad Ad) (Ad, error)
	GetAdsList(params AdsListParams, id uuid.UUID) ([]Ad, error)
}
