package main

import (
	"marketplace/internal/app"
	"marketplace/internal/config"
	"marketplace/internal/datasource"
	"marketplace/internal/di"
	"marketplace/internal/web"
	"go.uber.org/fx"
	"go.uber.org/zap"
)



func main() {

	app := fx.New(
		
		fx.Provide(
			config.MustLoad,
			provideLogger,
			app.NewJwtProvider,
			app.NewMarketService,
			app.NewUserService,
			datasource.NewStorage,
			datasource.NewMarketRepo,
			datasource.NewUserRepo,
			web.NewUserHandler,
			web.NewMarketHandler,
			func (repo *datasource.MarketRepo) app.MarketRepository{
				return repo
			},
			func (market *app.MarketService) app.MarketServicer{
				return market
			},
			func (repo *datasource.UserRepo) app.UserRepository{
				return repo
			},
			func (user *app.UserService) app.UserServicer{
				return user
			},

		),

		fx.Invoke(di.StartHTTPServer),
	)

	app.Run()
}

func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	switch cfg.Env {
	case "prod":
		zapCfg := zap.NewProductionConfig()
		return zapCfg.Build()
	default:
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.Encoding = "console"
		return zapCfg.Build()
	}
}
