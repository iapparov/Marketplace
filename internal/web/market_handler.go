package web

import (
	"encoding/json"
	"marketplace/internal/app"
	"marketplace/internal/config"
	"net/http"
	"strconv"
	"go.uber.org/zap"
	"github.com/google/uuid"
)

type AdsListResponse struct {
	Total int       `json:"total"`
	Items []app.Ad  `json:"items"`
}

type MarketHandler struct {
	app app.MarketServicer
	config *config.Config
	logger *zap.Logger
}

func NewMarketHandler(app app.MarketServicer, config *config.Config, logger *zap.Logger) *MarketHandler {
	return &MarketHandler{
		app: app,
		config: config,
		logger: logger,
	}
}

func (h *MarketHandler) NewAd(w http.ResponseWriter, r *http.Request) {
	var ad app.Ad
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		h.logger.Warn("invalid ad request body", zap.Error(err))
		http.Error(w, "bad request in body", http.StatusBadRequest)
		return
	}
	userIDVal := r.Context().Value(UserIDKey)
	useruuid, val_err := userIDVal.(string)
	if !val_err {
		h.logger.Warn("invalid user_id in context")
		http.Error(w, "invalid user_id in context", http.StatusBadRequest)
		return
	}
	uuid, err := uuid.Parse(useruuid)
	if err != nil {
		h.logger.Warn("invalid user_id format", zap.Error(err))
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}
	
	Adresp, err := h.app.NewAd(ad, *h.config, uuid)

	if err != nil {
		h.logger.Warn("failed to create new ad", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("new ad created successfully", zap.String("ad_id", Adresp.UUID.String()))
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Adresp)
	
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
		params.Order = "desc"
	}
	params.MinPrice, err = strconv.Atoi(rq.Get("min_price"))

	if err != nil {
		params.MinPrice = 0
	}

	params.MaxPrice, err = strconv.Atoi(rq.Get("max_price"))
	if err != nil || params.MaxPrice < 0 {
		params.MaxPrice = 1000000
	}

	userIDVal := r.Context().Value(UserIDKey)
	useruuid, val_err := userIDVal.(string)
	var id uuid.UUID
	if val_err {
		id, err = uuid.Parse(useruuid)
		if err != nil {
			id = uuid.Nil
		}
	} else {
		id = uuid.Nil
	}

	AdsList, err := h.app.AdsList(params, id)
	if err != nil {
		h.logger.Warn("failed to get ads list", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Info("ads list retrieved successfully", zap.Int("total", len(AdsList)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AdsList)
}