package db

import (
	"context"
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	CreateBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBooking(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DbName).Collection("bookings"),
	}
}

func (s *MongoBookingStore) CreateBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cursor, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *MongoBookingStore) GetBooking(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := ToObjectId(id)
	if err != nil {
		return nil, err
	}
	var booking *types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := ToObjectId(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	return err
}
