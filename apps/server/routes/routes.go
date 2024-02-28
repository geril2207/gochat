package routes

import (
	"github.com/geril2207/gochat/apps/server/auth"
	"github.com/geril2207/gochat/packages/config"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	server *echo.Echo,
	authController *auth.AuthController,
	config config.EnvConfig,
) {
	apiGroup := server.Group("/api")
	jwtMiddlewareConifg := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: []byte(config.JwtKey),
	}
	jwtMiddleware := echojwt.WithConfig(jwtMiddlewareConifg)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/register", authController.Register)
	authGroup.Use(jwtMiddleware)
	authGroup.POST("/refresh", authController.Refresh)
}
