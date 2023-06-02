package db

import (
	"github.com/prodigy00/hotel-reservation-api/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUser(string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
	}
}
