package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Hagrid",
		LastName:  "Hogwarts",
	}

	return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("Hagrid single user")
}
