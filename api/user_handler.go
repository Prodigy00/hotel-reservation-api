package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		id = c.Params("id")
	)

	user, err := h.userStore.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errs := params.Validate(); len(errs) > 0 {
		return c.JSON(errs)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	createdUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(createdUser)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		params = types.UpdateUserParams{}
		userId = c.Params("id")
	)

	if err := c.BodyParser(&params); err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"status": fmt.Sprintf("User with id: %s updated successfully", userId),
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}

	return c.JSON(map[string]string{"status": "OK"})
}
