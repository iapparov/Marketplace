package web

import (
	"net/http"
)

type MarketHandler struct {

}

func NewMarketHandler() *MarketHandler {
	return &MarketHandler{}
}

func (h *MarketHandler) NewAd(w http.ResponseWriter, r *http.Request) {
}

func (h *MarketHandler) AdsList(w http.ResponseWriter, r *http.Request) {
}