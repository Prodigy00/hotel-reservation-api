package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DbName     = "hotel-reservation"
	DbUri      = "mongodb://localhost:27017"
	TestDbName = "hotel-reservation-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func ToObjectId(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}
