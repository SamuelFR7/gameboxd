package handlers

import (
	"gameboxd/data"

	"github.com/gofiber/fiber/v2"
)

func HandleCreateRating(c *fiber.Ctx) error {
	u := new(data.Rating)

	if err := c.BodyParser(u); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	err := u.CreateRating()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusCreated)
}
