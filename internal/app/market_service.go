package app

import (
	"errors"
	"fmt"
	"marketplace/internal/config"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
	"github.com/google/uuid"
)


func NewMarketService(marketrepo MarketRepository, userrepo UserRepository) *MarketService {
	return &MarketService{
		Marketrepo: marketrepo,
		Userrepo:   userrepo,
	}
}

func (s *MarketService) NewAd(ad Ad, config config.Config, userid uuid.UUID) (Ad, error) {

	if ad.Title == "" || ad.Description == "" || ad.ImageURL == "" || ad.Price < config.Ad.PriceMin {
		return Ad{}, errors.New("all fields must be filled and price must be greater than zero")
	}

	if utf8.RuneCountInString(ad.Title) < config.Ad.MinLengthTitle || utf8.RuneCountInString(ad.Title) > config.Ad.MaxLengthTitle {
		return Ad{}, fmt.Errorf("title must be between %d and %d characters", config.Ad.MinLengthTitle, config.Ad.MaxLengthTitle)
	}

	if utf8.RuneCountInString(ad.Description) < config.Ad.MinLengthDescription || utf8.RuneCountInString(ad.Description) > config.Ad.MaxLengthDescription {
		return Ad{}, fmt.Errorf("description must be between %d and %d characters", config.Ad.MinLengthDescription, config.Ad.MaxLengthDescription)
	}

	ext := filepath.Ext(strings.ToLower(ad.ImageURL))

	if !config.Ad.AllowedImgTypesMap[ext] {
		return Ad{}, fmt.Errorf("image type %s is not allowed", ext)
	}

	ad.CreatedAt = time.Now()
	ad.UserID = userid
	user, err := s.Userrepo.FindByUUID(userid.String())
	if err != nil {
		return Ad{}, fmt.Errorf("FindByUUID error: %w", err)
	}
	ad.Username = user.Login
	ad.UUID = uuid.New()
	return s.Marketrepo.SaveAd(ad)
}

func (s *MarketService) AdsList(params AdsListParams, id uuid.UUID) ([]Ad, error) {

	Adslist, err := s.Marketrepo.GetAdsList(params, id.String())
	if err != nil {
		return nil, fmt.Errorf("getadslist error: %w", err)
	}
	if len(Adslist) == 0 {
		return nil, errors.New("list is empty")
	}

	return Adslist, nil
}