package handlers

import (
	"gameboxd/data"
	"gameboxd/services"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleListGames(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	name := c.Query("name")

	result, err := data.ListGames(page, name)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func HandleGetGameBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	result, err := data.GetGameBySlug(slug)

	if err != nil {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func HandleImportGames(c *fiber.Ctx) error {
	authResponse, err := services.GetAuthToken()
	if err != nil {
		return fiber.ErrInternalServerError
	}
	client := &http.Client{}

	dbCount, err := data.GetTotalGames()
	if err != nil {
		return fiber.ErrInternalServerError
	}

	apiCount, err := services.GetTotalGames(authResponse.AccessToken, client)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if dbCount == apiCount {
		return c.SendStatus(fiber.StatusOK)
	}

	missingGamesQt := apiCount - dbCount

	importGames := func() error {
		limit := 500
		callsUntilFinish := int(math.Ceil(float64(missingGamesQt) / float64(limit)))

		offset := dbCount
		lastTime := time.Now().Add(-250 * 1e6)

		for i := 0; i < callsUntilFinish; i++ {
			executionDuration := time.Since(lastTime).Milliseconds()
			if executionDuration < 250 {
				difference := 250 - executionDuration

				time.Sleep(time.Duration(difference * 1e6))
			}

			games, err := services.GetGames(offset, limit, authResponse.AccessToken, client)
			if err != nil {
				return err
			}
			lastTime = time.Now()

			err = data.CreateMultipleGamesFromIgdb(games)

			if err != nil {
				log.Printf("Stopped at offset %d\n", offset)
				return err
			}

			offset += 500
		}

		return nil
	}

	if missingGamesQt > 0 {
		go importGames()
	}

	return c.SendStatus(fiber.StatusOK)
}
