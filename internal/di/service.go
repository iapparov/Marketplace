package di

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"marketplace/internal/config"
	"marketplace/internal/web"
)

func StartHTTPServer(lc fx.Lifecycle, user_handler *web.UserHandler, market_handler *web.MarketHandler, config *config.Config) {
	// Регистрируем маршруты
	router := chi.NewRouter()
	web.RegisterRoutes(router, user_handler, market_handler)

	addres := fmt.Sprintf(":%d", config.Http_port)
	server := &http.Server{
		Addr:    addres,
		Handler: router, 
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server started on http://localhost:%d\n", config.Http_port)
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return server.Close()
		},
	})
}