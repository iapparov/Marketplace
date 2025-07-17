package main

import (

	"go.uber.org/fx"
	"marketplace/internal/config"
	"marketplace/internal/di"
	"marketplace/internal/web"
	"marketplace/internal/datasource"
)



func main() {

	app := fx.New(
		
		fx.Provide(
			config.MustLoad,
			datasource.New,
			web.NewUserHandler,
			web.NewMarketHandler,
		),

		fx.Invoke(di.StartHTTPServer),
	)

	app.Run()
}