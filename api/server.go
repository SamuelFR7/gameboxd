package api

import (
	"gameboxd/db"
	"gameboxd/handlers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer() *fiber.App {
	db.Init()

	cookieSecret := os.Getenv("COOKIE_SECRET")

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cookieSecret,
	}))

	api := app.Group("/api")

	v1 := api.Group("/v1")

	games := v1.Group("/games")
	ratings := v1.Group("/ratings")

	games.Use(AuthMiddleware)
	ratings.Use(AuthMiddleware)

	games.Get("/", handlers.HandleListGames)
	games.Get("/:slug", handlers.HandleGetGameBySlug)
	games.Post("/import", handlers.HandleImportGames)

	ratings.Post("/", handlers.HandleCreateRating)

	v1.Post("/users/", handlers.HandleCreateUser)
	v1.Post("/users/sign-in", handlers.SignIn)
	v1.Post("/users/sign-out", handlers.SignOut)

	return app
}
