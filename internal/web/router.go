package web

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, userHandler *UserHandler, marketHandler *MarketHandler) {
	r.Post("/login", userHandler.Login)
	r.Post("/register", userHandler.Register)


	r.Post("/new-ad", marketHandler.NewAd)
	r.Get("/ads-list", marketHandler.AdsList)
}