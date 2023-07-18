package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/omelaymy/users/config"
	_ "github.com/omelaymy/users/docs"
	"github.com/omelaymy/users/internal/api/http/delivery"
	"github.com/omelaymy/users/pkg/di"

	"github.com/samber/do"
)

// @title Swagger Users API
// @description This is API for service Users.
// @version 1.0
// @securityDefinitions.basic BasicAuth
// @name Authorization
// @host localhost:8888
// @BasePath /api
func main() {
	i := do.New()

	do.Provide(i, di.NewFlags)
	do.Provide(i, di.NewConfig)
	do.Provide(i, di.NewLogger)
	do.Provide(i, di.NewInMemoryDatabase)
	do.Provide(i, di.NewAuth)
	do.Provide(i, di.NewAuthRepository)
	do.Provide(i, di.NewUsers)
	do.Provide(i, di.NewUsersRepository)
	do.Provide(i, di.NewRoutes)
	do.Provide(i, di.NewHandlers)
	do.Provide(i, di.NewMWManager)
	do.Provide(i, di.NewHttpErrorHandler)
	do.Provide(i, di.NewFiberApp)
	do.Provide(i, di.NewTranslator)
	do.Provide(i, di.NewValidate)

	app := do.MustInvoke[*fiber.App](i)
	routes := do.MustInvoke[*delivery.Routes](i)
	routes.RegisterRoutes()

	cfg := do.MustInvoke[*config.Config](i)
	log.Fatal(app.Listen(cfg.Server.Address))
}
