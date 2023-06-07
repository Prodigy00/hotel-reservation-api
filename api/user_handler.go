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
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.store.User.GetUser(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
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

	createdUser, err := h.store.User.CreateUser(c.Context(), user)
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
	if err := h.store.User.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}

	return c.JSON(map[string]string{
		"status": fmt.Sprintf("User with id: %s updated successfully", userId),
	})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	if err := h.store.User.DeleteUser(c.Context(), userId); err != nil {
		return err
	}

	return c.JSON(map[string]string{"status": "OK"})
}
