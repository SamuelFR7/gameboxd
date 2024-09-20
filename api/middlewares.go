package api

import (
	"gameboxd/data"
	"gameboxd/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sessionId := c.Cookies(utils.SESSION_KEY)

	if sessionId == "" {
		return fiber.ErrUnauthorized
	}

	_, err := data.ValidateSession(sessionId)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}
