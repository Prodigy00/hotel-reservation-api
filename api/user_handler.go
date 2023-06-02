package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUser(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Hagrid",
		LastName:  "Hogwarts",
	}

	return c.JSON(u)
}
