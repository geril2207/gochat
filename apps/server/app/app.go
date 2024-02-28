package app

import (
	"context"
	"fmt"

	"github.com/geril2207/gochat/apps/server/auth"
	"github.com/geril2207/gochat/apps/server/routes"
	"github.com/geril2207/gochat/packages/config"
	"github.com/geril2207/gochat/packages/db"
	"github.com/geril2207/gochat/packages/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		config.ProvideEnvConfig,
		func(config config.EnvConfig, lc fx.Lifecycle) *pgxpool.Pool {
			pool, err := db.DbConnect(config.DatabaseUrl)
			if err != nil {
				panic(err)
			}
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return db.DbCloseConnection(pool)
				},
			})
			return pool
		},
		services.ProvideUsersService,
		auth.ProvideAuthController,
		NewServer,
	), fx.Invoke(InvokeServer, routes.SetupRoutes),
)

func NewServer(settings config.EnvConfig) *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Recover(),
	)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${uri} ${method} ${status} ${latency_human}\n",
	}))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Validator = &RequestValidator{validator: NewValidator()}

	return e
}

func InvokeServer(lifecycle fx.Lifecycle, server *echo.Echo, settings config.EnvConfig) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(":" + settings.ServerPort)
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Http Server.")
			return server.Shutdown(ctx)
		},
	})
}
