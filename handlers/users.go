package handlers

import (
	"gameboxd/data"
	"gameboxd/utils"

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

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Path:     "/",
		Name:     utils.SESSION_KEY,
		Value:    session.Id,
		MaxAge:   30 * 86400,
	})

	return c.SendStatus(fiber.StatusOK)
}

func SignOut(c *fiber.Ctx) error {
	sessionId := c.Cookies(utils.SESSION_KEY)

	if sessionId == "" {
		return fiber.ErrUnauthorized
	}

	err := data.DestroySession(sessionId)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	c.Cookie(&fiber.Cookie{
		HTTPOnly: true,
		Path:     "/",
		Name:     utils.SESSION_KEY,
		Value:    "",
		MaxAge:   0,
	})

	return c.SendStatus(fiber.StatusOK)
}
