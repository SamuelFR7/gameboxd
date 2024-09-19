package handlers

import (
	"fmt"
	"gameboxd/data"
	"gameboxd/utils"
	"os"

	"github.com/gofiber/fiber/v2"
)

type signInUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleCreateUser(c *fiber.Ctx) error {
	user := new(data.CreateUserParams)

	if err := c.BodyParser(user); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	passwordHash, err := utils.GeneratePasswordHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = passwordHash

	err = data.CreateUser(*user)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusCreated)
}

func SignIn(c *fiber.Ctx) error {
	user := new(signInUserRequest)

	if err := c.BodyParser(user); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	userExists, err := data.GetUserByEmail(user.Email)
	if err != nil {
		return fiber.ErrNotFound
	}

	passwordMatches := utils.ComparePasswordHash(userExists.PasswordHash, user.Password)

	if passwordMatches == false {
		return fiber.ErrUnauthorized
	}

	session, err := data.CreateSession(userExists.Id)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	var secret = os.Getenv("COOKIE_SECRET")

	if secret == "" {
		secret = "my_secret"
	}

	sessionHashed := utils.HashCookie(session.Id, secret)

	c.Response().Header.Add("Set-Cookie", fmt.Sprintf("_session=%s; HttpOnly; Path=/", sessionHashed))

	return c.SendStatus(fiber.StatusOK)
}
