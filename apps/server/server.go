package main

import (
	"fmt"

	"github.com/geril2207/gochat/apps/server/app"
	"go.uber.org/fx"
)

// @title GoChat Api Documentation
// @description This is a simpe chat server implementation in golang
// @version 1.0
// @BasePath /api
func main() {
	fmt.Print(1)
	fx.New(
		app.App,
	).Run()
}
