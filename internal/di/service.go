package di

import (
	"context"
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"marketplace/internal/config"
	"marketplace/internal/web"
	"go.uber.org/zap"
)


func StartHTTPServer(lc fx.Lifecycle, user_handler *web.UserHandler, market_handler *web.MarketHandler, config *config.Config, logger *zap.Logger) {
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
			logger.Info("Server started", zap.Int("port", config.Http_port))
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Shutting down server...")
			return server.Close()
		},
	})
}