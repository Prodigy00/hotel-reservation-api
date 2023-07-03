package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type BookingsHandler struct {
	store *db.Store
}

func NewBookingsHandler(store *db.Store) *BookingsHandler {
	return &BookingsHandler{
		store: store,
	}
}

// needs to be admin authorized
func (h *BookingsHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookingsHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBooking(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericRes{
			Status: "error",
			Msg:    "not authorized",
		})
	}
	return c.JSON(booking)
}

func (h *BookingsHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBooking(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericRes{
			Status: "error",
			Msg:    "not authorized",
		})
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericRes{Status: "success", Msg: "booking canceled"})
}
