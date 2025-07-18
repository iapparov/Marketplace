package web

import (
	"net/http"
	"marketplace/internal/app"
	"encoding/json"
	"strconv"
	"github.com/google/uuid"
)

type AdsListResponse struct {
	Total int       `json:"total"`
	Items []app.Ad  `json:"items"`
}

type MarketHandler struct {
	app *app.MarketService
}

func NewMarketHandler(app *app.MarketService) *MarketHandler {
	return &MarketHandler{
		app: app,
	}
}

func (h *MarketHandler) NewAd(w http.ResponseWriter, r *http.Request) {
	
}

func (h *MarketHandler) AdsList(w http.ResponseWriter, r *http.Request) {

	rq := r.URL.Query()
	var params app.AdsListParams
	var err error
	params.Page, err = strconv.Atoi(rq.Get("page"))
	if params.Page == 0 || params.Page < 1 || err != nil {
		params.Page = 1
	}
	params.Limit, err = strconv.Atoi(rq.Get("limit"))
	if params.Limit == 0 || params.Limit < 1 || err != nil {
		params.Limit = 10
	}
	params.SortBy = rq.Get("sort_by")
	if params.SortBy != "date" && params.SortBy != "price" {
		params.SortBy = "date"
	}
	params.Order = rq.Get("order")
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "asc"
	}
	params.MinPrice, err = strconv.Atoi(rq.Get("min_price"))

	if err != nil {
		params.MinPrice = 0
	}

	params.MaxPrice, err = strconv.Atoi(rq.Get("max_price"))
	if err != nil || params.MaxPrice < 0 {
		params.MaxPrice = 1000000
	}


	// Примерный проброс для контекста с user_id, доделать
	userIDVal := r.Context().Value("user_id")
	useruuid, ok := userIDVal.(string)
	if !ok {
		http.Error(w, "invalid user_id in context", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(useruuid)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	AddList, err := h.app.AddList(params, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AddList)
}