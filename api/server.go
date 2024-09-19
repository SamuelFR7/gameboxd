package api

import (
	"gameboxd/api/db"
	"gameboxd/api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer() *fiber.App {
    db.CreateDatabase()

    app := fiber.New()
    
    app.Use(logger.New())
    app.Use(recover.New())

    api := app.Group("/api")

    v1 := api.Group("/v1")

    v1.Get("/games", handlers.HandleListGames)
    v1.Get("/games/:slug", handlers.HandleGetGameBySlug)
    v1.Post("/games/import", handlers.HandleImportGames)
    v1.Post("/ratings", handlers.HandleCreateRating)

    return app
}
