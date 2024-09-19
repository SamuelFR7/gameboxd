package api

import (
	"gameboxd/db"
	"gameboxd/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer() *fiber.App {
	db.Init()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Get("/games", handlers.HandleListGames)
	v1.Get("/games/:slug", handlers.HandleGetGameBySlug)
	v1.Post("/games/import", handlers.HandleImportGames)

	v1.Post("/ratings", handlers.HandleCreateRating)

	v1.Post("/users/", handlers.HandleCreateUser)
	v1.Post("/users/sign-in", handlers.SignIn)

	return app
}
