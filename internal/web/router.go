package web

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, userHandler *UserHandler, marketHandler *MarketHandler) {
	r.Post("/login", userHandler.Login)
	r.Post("/register", userHandler.Register)
	
	r.With(OptionalAuthMiddleware(userHandler.jwt)).Get("/ads-list", marketHandler.AdsList)

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(userHandler.jwt))
		r.Post("/new-ad", marketHandler.NewAd)
		r.Post("/refresh-access-token", userHandler.RefreshAccessToken)
	})
}