package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/prodigy00/hotel-reservation-api/db"
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (b *BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(b.FromDate) || now.After(b.TillDate) {
		return fmt.Errorf("cannot backdate room booking")
	}
	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.Validate(); err != nil {
		return err
	}
	roomID := c.Params("id")
	roomOID, err := db.ToObjectId(roomID)
	if err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericRes{
			Status: "error",
			Msg:    "internal server error",
		})
	}

	isAvailable, err := h.isRoomAvailableForBooking(c.Context(), roomOID, params)

	if err != nil {
		return err
	}

	if !isAvailable {
		return c.Status(http.StatusBadRequest).JSON(genericRes{
			Status: "error",
			Msg:    fmt.Sprintf("room %s already booked", roomID),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	created, err := h.store.Booking.CreateBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(created)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	condition := bson.M{
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"roomID": roomID,
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, condition)
	if err != nil {
		return false, err
	}
	//fmt.Println(bookings)
	isAvailable := len(bookings) == 0

	return isAvailable, nil
}
