package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}

	if !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}
	return c.Next()
}
